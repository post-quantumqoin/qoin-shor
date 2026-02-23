package miner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/arc/v2"
	"github.com/ipfs/go-cid"
	"go.opencensus.io/trace"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/abi"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/crypto"
	"github.com/post-quantumqoin/core-types/proof"

	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/api/v1api"
	"github.com/post-quantumqoin/qoin-shor/build"
	cliutil "github.com/post-quantumqoin/qoin-shor/cli/util"
	"github.com/post-quantumqoin/qoin-shor/core/contracts/builtin"
	"github.com/post-quantumqoin/qoin-shor/core/contracts/policy"
	"github.com/post-quantumqoin/qoin-shor/core/gen"
	"github.com/post-quantumqoin/qoin-shor/core/gen/slashfilter"
	lrand "github.com/post-quantumqoin/qoin-shor/core/rand"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/journal"
	"github.com/post-quantumqoin/qoin-shor/pqccrypto/shake3"
	"github.com/post-quantumqoin/qoin-shor/pqcpow"
)

// var log = logging.Logger("miner")

// Journal event types.
// const (
// 	evtTypeBlockMined = iota
// )

// waitFunc is expected to pace block mining at the configured network rate.
//
// baseTime is the timestamp of the mining base, i.e. the timestamp
// of the tipset we're planning to construct upon.
//
// Upon each mining loop iteration, the returned callback is called reporting
// whether we mined a block in this round or not.
// type waitFunc func(ctx context.Context, baseTime uint64) (func(bool, abi.ChainEpoch, error), abi.ChainEpoch, error)

// func randTimeOffset(width time.Duration) time.Duration {
// 	buf := make([]byte, 8)
// 	rand.Reader.Read(buf) //nolint:errcheck
// 	val := time.Duration(binary.BigEndian.Uint64(buf) % uint64(width))

// 	return val - (width / 2)
// }

// NewMiner instantiates a miner with a concrete WinningPoStProver and a miner
// address (which can be different from the worker's address).
func NewPqcMiner(api v1api.FullNode, addr address.Address, sf *slashfilter.SlashFilter, j journal.Journal) *PqcMiner {
	arc, err := arc.NewARC[abi.ChainEpoch, bool](10000)
	if err != nil {
		panic(err)
	}
	log.Infow("NewPqcMiner ", "addr:", addr.String())
	return &PqcMiner{
		api:     api,
		epp:     nil,
		address: addr,
		propagationWaitFunc: func(ctx context.Context, baseTime uint64) (func(bool, abi.ChainEpoch, error), abi.ChainEpoch, error) {
			// wait around for half the block time in case other parents come in
			//
			// if we're mining a block in the past via catch-up/rush mining,
			// such as when recovering from a network halt, this sleep will be
			// for a negative duration, and therefore **will return
			// immediately**.
			//
			// the result is that we WILL NOT wait, therefore fast-forwarding
			// and thus healing the chain by backfilling it with null rounds
			// rapidly.
			log.Infow("propagationWaitFunc 1",
				"time", build.Clock.Now())
			// bTime := uint64(build.Clock.Now().Unix())
			deadline := baseTime + build.PropagationDelaySecs
			// deadline := bTime + build.PropagationDelaySecs
			baseT := time.Unix(int64(deadline), 0)

			baseT = baseT.Add(randTimeOffset(time.Second))

			build.Clock.Sleep(build.Clock.Until(baseT))
			log.Infow("propagationWaitFunc 2",
				"time", build.Clock.Now())
			return func(bool, abi.ChainEpoch, error) {}, 0, nil
		},

		sf:                sf,
		minedBlockHeights: arc,
		evtTypes: [...]journal.EventType{
			evtTypeBlockMined: j.RegisterEventType("miner", "block_mined"),
		},
		journal: j,
	}
}

// Miner encapsulates the mining processes of the system.
//
// Refer to the godocs on mineOne and mine methods for more detail.
type PqcMiner struct {
	api v1api.FullNode

	epp gen.WinningPoStProver

	lk       sync.Mutex
	address  address.Address
	stop     chan struct{}
	stopping chan struct{}

	propagationWaitFunc waitFunc

	// lastWork holds the last MiningBase we built upon.
	lastWork *MiningBase

	sf *slashfilter.SlashFilter
	// minedBlockHeights is a safeguard that caches the last heights we mined.
	// It is consulted before publishing a newly mined block, for a sanity check
	// intended to avoid slashings in case of a bug.
	minedBlockHeights *arc.ARCCache[abi.ChainEpoch, bool]

	evtTypes [1]journal.EventType
	journal  journal.Journal
}

// Address returns the address of the miner.
func (m *PqcMiner) Address() address.Address {
	m.lk.Lock()
	defer m.lk.Unlock()

	return m.address
}

// Start starts the mining operation. It spawns a goroutine and returns
// immediately. Start is not idempotent.
func (m *PqcMiner) Start(_ context.Context) error {
	m.lk.Lock()
	defer m.lk.Unlock()
	if m.stop != nil {
		return fmt.Errorf("PqcMiner already started")
	}
	m.stop = make(chan struct{})
	go m.mine(context.TODO())
	return nil
}

// Stop stops the mining operation. It is not idempotent, and multiple adjacent
// calls to Stop will fail.
func (m *PqcMiner) Stop(ctx context.Context) error {
	m.lk.Lock()

	m.stopping = make(chan struct{})
	stopping := m.stopping
	close(m.stop)

	m.lk.Unlock()

	select {
	case <-stopping:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *PqcMiner) niceSleep(d time.Duration) bool {
	select {
	case <-build.Clock.After(d):
		return true
	case <-m.stop:
		log.Infow("received interrupt while trying to sleep in mining cycle")
		return false
	}
}

// mine runs the mining loop. It performs the following:
//
//  1. Queries our current best currently-known mining candidate (tipset to
//     build upon).
//  2. Waits until the propagation delay of the network has elapsed (currently
//     15 seconds). The waiting is done relative to the timestamp of the best
//     candidate, which means that if it's way in the past, we won't wait at
//     all (e.g. in catch-up or rush mining).
//  3. After the wait, we query our best mining candidate. This will be the one
//     we'll work with.
//  4. Sanity check that we _actually_ have a new mining base to mine on. If
//     not, wait one epoch + propagation delay, and go back to the top.
//  5. We attempt to mine a block, by calling mineOne (refer to godocs). This
//     method will either return a block if we were eligible to mine, or nil
//     if we weren't.
//     6a. If we mined a block, we update our state and push it out to the network
//     via gossipsub.
//     6b. If we didn't mine a block, we consider this to be a nil round on top of
//     the mining base we selected. If other PqcMiner or miners on the network
//     were eligible to mine, we will receive their blocks via gossipsub and
//     we will select that tipset on the next iteration of the loop, thus
//     discarding our null round.
func (m *PqcMiner) mine(ctx context.Context) {
	ctx, span := trace.StartSpan(ctx, "/mine")
	defer span.End()

	// Perform the Winning PoSt warmup in a separate goroutine.
	// go m.doWinPoStWarmup(ctx)

	var lastBase MiningBase

	// Start the main mining loop.
minerLoop:
	for {
		// Prepare a context for a single node operation.
		ctx := cliutil.OnSingleNode(ctx)

		// Handle stop signals.
		select {
		case <-m.stop:
			// If a stop signal is received, clean up and exit the mining loop.
			stopping := m.stopping
			m.stop = nil
			m.stopping = nil
			close(stopping)
			return

		default:
		}

		var base *MiningBase // NOTE: This points to m.lastWork; Incrementing nulls here will increment it there.
		var onDone func(bool, abi.ChainEpoch, error)
		var injectNulls abi.ChainEpoch

		// Look for the best mining candidate.
		for {
			if !m.niceSleep(2 * time.Second) {
				continue minerLoop
			}
			prebase, err := m.GetBestMiningCandidate(ctx)
			if err != nil {
				log.Errorf("failed to get best mining candidate: %s", err)
				if !m.niceSleep(time.Second * 5) {
					continue minerLoop
				}
				continue
			}

			btime := time.Unix(int64(prebase.TipSet.MinTimestamp()), 0)
			log.Infow("mined block in the past",
				"block-time", btime, "time", build.Clock.Now(), "prebase NullRounds:", prebase.NullRounds)
			// Check if we have a new base or if the current base is still valid.
			if base != nil && base.TipSet.Height() == prebase.TipSet.Height() && base.NullRounds == prebase.NullRounds {
				base = prebase
				break
			}
			if base != nil {
				onDone(false, 0, nil)
			}

			// TODO: need to change the orchestration here. the problem is that
			// we are waiting *after* we enter this loop and selecta mining
			// candidate, which is almost certain to change in multiminer
			// tests. Instead, we should block before entering the loop, so
			// that when the test 'MineOne' function is triggered, we pull our
			// best mining candidate at that time.

			// Wait until propagation delay period after block we plan to mine on
			onDone, injectNulls, err = m.propagationWaitFunc(ctx, prebase.TipSet.MinTimestamp())
			if err != nil {
				log.Error(err)
				continue
			}

			// if prebase != nil && prebase.TipSet.Height() == build.UpgradeYellowStoneHeight {
			// 	prebase.NullRounds += abi.ChainEpoch((uint64(build.Clock.Now().Unix()) - prebase.TipSet.MinTimestamp()) / build.BlockDelaySecs)
			// 	m.niceSleep(time.Duration(build.BlockDelaySecs-(uint64(build.Clock.Now().Unix())-prebase.TipSet.MinTimestamp())%build.BlockDelaySecs) * time.Second)
			// }

			// Ensure the beacon entry is available before finalizing the mining base.
			_, err = m.api.StateGetBeaconEntry(ctx, prebase.TipSet.Height()+prebase.NullRounds+1)
			if err != nil {
				log.Errorf("failed getting beacon entry: %s", err)
				if !m.niceSleep(time.Second) {
					continue minerLoop
				}
				continue
			}

			base = prebase
		}

		base.NullRounds += injectNulls // Adjust for testing purposes.

		// Check for repeated mining candidates and handle sleep for the next round.
		if base.TipSet.Equals(lastBase.TipSet) && lastBase.NullRounds == base.NullRounds {
			log.Warnf("BestMiningCandidate from the previous round: %s (nulls:%d)", lastBase.TipSet.Cids(), lastBase.NullRounds)
			if !m.niceSleep(time.Duration(build.BlockDelaySecs) * time.Second) {
				continue minerLoop
			}
			continue
		}

		// Attempt to mine a block.
		b, err := m.mineOne(ctx, base)
		if err != nil {
			log.Warnf("mining block failed: %+v", err)
			if !m.niceSleep(time.Second) {
				continue minerLoop
			}
			onDone(false, 0, err)
			continue
		}
		lastBase = *base

		var h abi.ChainEpoch
		if b != nil {
			h = b.Header.Height
		}
		onDone(b != nil, h, nil)

		// Process the mined block.
		if b != nil {
			// Record the event of mining a block.
			m.journal.RecordEvent(m.evtTypes[evtTypeBlockMined], func() interface{} {
				return map[string]interface{}{
					// Data about the mined block.
					"parents":   base.TipSet.Cids(),
					"nulls":     base.NullRounds,
					"epoch":     b.Header.Height,
					"timestamp": b.Header.Timestamp,
					"cid":       b.Header.Cid(),
				}
			})

			// Check for slash filter conditions.
			if os.Getenv("LOTUS_MINER_NO_SLASHFILTER") != "_yes_i_know_i_can_and_probably_will_lose_all_my_fil_and_power_" && !build.IsNearUpgrade(base.TipSet.Height(), build.UpgradeWatermelonFixHeight) {
				witness, fault, err := m.sf.MinedBlock(ctx, b.Header, base.TipSet.Height()+base.NullRounds)
				if err != nil {
					log.Errorf("<!!> SLASH FILTER ERRORED: %s", err)
					// Continue here, because it's _probably_ wiser to not submit this block
					continue
				}

				if fault {
					log.Errorf("<!!> SLASH FILTER DETECTED FAULT due to blocks %s and %s", b.Header.Cid(), witness)
					continue
				}
			}

			// Check for blocks created at the same height.
			if _, ok := m.minedBlockHeights.Get(b.Header.Height); ok {
				log.Warnw("Created a block at the same height as another block we've created", "height", b.Header.Height, "PqcMiner", b.Header.Miner, "parents", b.Header.Parents)
				continue
			}

			// Add the block height to the mined block heights.
			m.minedBlockHeights.Add(b.Header.Height, true)

			// Submit the newly mined block.
			if err := m.api.SyncSubmitBlock(ctx, b); err != nil {
				log.Errorf("failed to submit newly mined block: %+v", err)
			}

			btime := time.Unix(int64(b.Header.Timestamp), 0)
			now := build.Clock.Now()
			// Handle timing for broadcasting the block.
			switch {
			case btime == now:
				// block timestamp is perfectly aligned with time.
			case btime.After(now):
				// Wait until it's time to broadcast the block.
				if !m.niceSleep(build.Clock.Until(btime)) {
					log.Warnf("received interrupt while waiting to broadcast block, will shutdown after block is sent out")
					build.Clock.Sleep(build.Clock.Until(btime))
				}
			default:
				// Log if the block was mined in the past.
				log.Warnw("mined block in the past",
					"block-time", btime, "time", build.Clock.Now(), "difference", build.Clock.Since(btime))
			}
		} else {
			// If no block was mined, increase the null rounds and wait for the next epoch.
			base.NullRounds++
			log.Info("no block was mined increase the null rounds:", base.NullRounds)
			// Calculate the time for the next round.
			nextRound := time.Unix(int64(base.TipSet.MinTimestamp()+build.BlockDelaySecs*uint64(base.NullRounds))+int64(build.PropagationDelaySecs), 0)

			// Wait for the next round or stop signal.
			select {
			case <-build.Clock.After(build.Clock.Until(nextRound)):
			case <-m.stop:
				stopping := m.stopping
				m.stop = nil
				m.stopping = nil
				close(stopping)
				return
			}
		}
	}
}

// MiningBase is the tipset on top of which we plan to construct our next block.
// Refer to godocs on GetBestMiningCandidate.
// type MiningBase struct {
// 	TipSet      *types.TipSet
// 	ComputeTime time.Time
// 	NullRounds  abi.ChainEpoch
// }

// GetBestMiningCandidate implements the fork choice rule from a miner's
// perspective.
//
// It obtains the current chain head (HEAD), and compares it to the last tipset
// we selected as our mining base (LAST). If HEAD's weight is larger than
// LAST's weight, it selects HEAD to build on. Else, it selects LAST.
func (m *PqcMiner) GetBestMiningCandidate(ctx context.Context) (*MiningBase, error) {
	m.lk.Lock()
	defer m.lk.Unlock()

	bts, err := m.api.ChainHead(ctx)
	if err != nil {
		return nil, err
	}

	if m.lastWork != nil {
		if m.lastWork.TipSet.Equals(bts) {
			return m.lastWork, nil
		}

		btsw, err := m.api.ChainTipSetWeight(ctx, bts.Key())
		if err != nil {
			return nil, err
		}
		ltsw, err := m.api.ChainTipSetWeight(ctx, m.lastWork.TipSet.Key())
		if err != nil {
			m.lastWork = nil
			return nil, err
		}

		if types.BigCmp(btsw, ltsw) <= 0 {
			return m.lastWork, nil
		}
	}

	m.lastWork = &MiningBase{TipSet: bts, ComputeTime: time.Now()}
	return m.lastWork, nil
}

// mineOne attempts to mine a single block, and does so synchronously, if and
// only if we are eligible to mine.
//
// {hint/landmark}: This method coordinates all the steps involved in mining a
// block, including the condition of whether mine or not at all depending on
// whether we win the round or not.
//
// This method does the following:
//
//	1.
func (m *PqcMiner) mineOne(ctx context.Context, base *MiningBase) (minedBlock *types.BlockMsg, err error) {
	log.Debugw("attempting to mine a block", "tipset", types.LogCids(base.TipSet.Cids()))
	tStart := build.Clock.Now()

	round := base.TipSet.Height() + base.NullRounds + 1

	// always write out a log
	var winner *types.ElectionProof
	var mbi *api.MiningBaseInfo
	var rbase types.BeaconEntry
	defer func() {

		var hasMinPower bool

		// mbi can be nil if we are deep in penalty and there are 0 eligible sectors
		// in the current deadline. If this case - put together a dummy one for reporting
		// https://github.com/post-quantumqoin/qoin-shor/blob/v1.9.0/chain/stmgr/utils.go#L500-L502
		if mbi == nil {
			mbi = &api.MiningBaseInfo{
				NetworkPower:      big.NewInt(-1), // we do not know how big the network is at this point
				EligibleForMining: false,
				MinerPower:        big.NewInt(0), // but we do know we do not have anything eligible
			}

			// try to opportunistically pull actual power and plug it into the fake mbi
			if pow, err := m.api.StateMinerPower(ctx, m.address, base.TipSet.Key()); err == nil && pow != nil {
				hasMinPower = pow.HasMinPower
				mbi.MinerPower = pow.MinerPower.QualityAdjPower
				mbi.NetworkPower = pow.TotalPower.QualityAdjPower
			}
		}

		isLate := uint64(tStart.Unix()) > (base.TipSet.MinTimestamp() + uint64(base.NullRounds*builtin.EpochDurationSeconds) + build.PropagationDelaySecs)

		logStruct := []interface{}{
			"tookMilliseconds", (build.Clock.Now().UnixNano() - tStart.UnixNano()) / 1_000_000,
			"forRound", int64(round),
			"baseEpoch", int64(base.TipSet.Height()),
			"baseDeltaSeconds", uint64(tStart.Unix()) - base.TipSet.MinTimestamp(),
			"nullRounds", int64(base.NullRounds),
			"lateStart", isLate,
			"beaconEpoch", rbase.Round,
			"lookbackEpochs", int64(policy.ChainFinality), // hardcoded as it is unlikely to change again: https://github.com/post-quantumqoin/qoin-shor/blob/v1.8.0/chain/actors/policy/policy.go#L180-L186
			"networkPowerAtLookback", mbi.NetworkPower.String(),
			"minerPowerAtLookback", mbi.MinerPower.String(),
			"isEligible", mbi.EligibleForMining,
			"isWinner", (winner != nil),
			"error", err,
		}

		if err != nil {
			log.Errorw("completed mineOne", logStruct...)
		} else if isLate || (hasMinPower && !mbi.EligibleForMining) {
			// log.Warnw("completed mineOne", logStruct...)
		} else {
			// log.Infow("completed mineOne", logStruct...)
		}
	}()

	log.Infow("MinerGetBaseInfo ", " m.address:", m.address.String())
	mbi, err = m.api.MinerGetBaseInfo(ctx, m.address, round, base.TipSet.Key())
	if err != nil {
		err = xerrors.Errorf("failed to get mining base info: %w", err)
		return nil, err
	}
	if mbi == nil {
		return nil, nil
	}

	// if !mbi.EligibleForMining {
	// 	// slashed or just have no power yet
	// 	return nil, nil
	// }

	tPowercheck := build.Clock.Now()

	bvals := mbi.BeaconEntries
	rbase = mbi.PrevBeaconEntry
	if len(bvals) > 0 {
		rbase = bvals[len(bvals)-1]
	}
	log.Infow("computeTicket ", "WorkerKey:", mbi.WorkerKey.String(), "tPowercheck:", tPowercheck)
	ticket, err := m.computeTicket(ctx, &rbase, round, base.TipSet.MinTicket(), mbi)
	if err != nil {
		err = xerrors.Errorf("scratching ticket failed: %w", err)
		return nil, err
	}

	// winner, err = gen.IsRoundWinner(ctx, round, m.address, rbase, mbi, m.api)
	// if err != nil {
	// 	err = xerrors.Errorf("failed to check if we win next round: %w", err)
	// 	return nil, err
	// }

	// if winner == nil {
	// 	return nil, nil
	// }

	tTicket := build.Clock.Now()

	buf := new(bytes.Buffer)
	if err := m.address.MarshalCBOR(buf); err != nil {
		err = xerrors.Errorf("failed to marshal PqcMiner address: %w", err)
		return nil, err
	}

	// rand, err := lrand.DrawRandomnessFromBase(rbase.Data, crypto.DomainSeparationTag_WinningPoStChallengeSeed, round, buf.Bytes())
	// if err != nil {
	// 	err = xerrors.Errorf("failed to get randomness for winning post: %w", err)
	// 	return nil, err
	// }

	// prand := abi.PoStRandomness(rand)

	tSeed := build.Clock.Now()
	// nv, err := m.api.StateNetworkVersion(ctx, base.TipSet.Key())
	// if err != nil {
	// 	return nil, err
	// }

	// postProof, err := m.epp.ComputeProof(ctx, mbi.Sectors, prand, round, nv)
	// if err != nil {
	// 	err = xerrors.Errorf("failed to compute winning post proof: %w", err)
	// 	return nil, err
	// }

	tProof := build.Clock.Now()
	log.Infow("MpoolSelect")
	// get pending messages early,
	msgs, err := m.api.MpoolSelect(ctx, base.TipSet.Key(), ticket.Quality())
	if err != nil {
		err = xerrors.Errorf("failed to select messages for block: %w", err)
		return nil, err
	}

	tPending := build.Clock.Now()

	// This next block exists to "catch" equivocating miners,
	// who submit 2 blocks at the same height at different times in order to split the network.
	// To safeguard against this, we make sure it's been EquivocationDelaySecs since our base was calculated,
	// then re-calculate it.
	// If the daemon detected equivocated blocks, those blocks will no longer be in the new base.
	m.niceSleep(time.Until(base.ComputeTime.Add(time.Duration(build.EquivocationDelaySecs) * time.Second)))
	newBase, err := m.GetBestMiningCandidate(ctx)
	if err != nil {
		err = xerrors.Errorf("failed to refresh best mining candidate: %w", err)
		return nil, err
	}

	tEquivocateWait := build.Clock.Now()

	// If the base has changed, we take the _intersection_ of our old base and new base,
	// thus ejecting blocks from any equivocating miners, without taking any new blocks.
	if newBase.TipSet.Height() == base.TipSet.Height() && !newBase.TipSet.Equals(base.TipSet) {
		log.Warnf("base changed from %s to %s, taking intersection", base.TipSet.Key(), newBase.TipSet.Key())
		newBaseMap := map[cid.Cid]struct{}{}
		for _, newBaseBlk := range newBase.TipSet.Cids() {
			newBaseMap[newBaseBlk] = struct{}{}
		}

		refreshedBaseBlocks := make([]*types.BlockHeader, 0, len(base.TipSet.Cids()))
		for _, baseBlk := range base.TipSet.Blocks() {
			if _, ok := newBaseMap[baseBlk.Cid()]; ok {
				refreshedBaseBlocks = append(refreshedBaseBlocks, baseBlk)
			}
		}

		if len(refreshedBaseBlocks) != 0 && len(refreshedBaseBlocks) != len(base.TipSet.Blocks()) {
			refreshedBase, err := types.NewTipSet(refreshedBaseBlocks)
			if err != nil {
				err = xerrors.Errorf("failed to create new tipset when refreshing: %w", err)
				return nil, err
			}

			if !base.TipSet.MinTicket().Equals(refreshedBase.MinTicket()) {
				log.Warn("recomputing ticket due to base refresh")

				ticket, err = m.computeTicket(ctx, &rbase, round, refreshedBase.MinTicket(), mbi)
				if err != nil {
					err = xerrors.Errorf("failed to refresh ticket: %w", err)
					return nil, err
				}
			}

			log.Warn("re-selecting messages due to base refresh")
			// refresh messages, as the selected messages may no longer be valid
			msgs, err = m.api.MpoolSelect(ctx, refreshedBase.Key(), ticket.Quality())
			if err != nil {
				err = xerrors.Errorf("failed to re-select messages for block: %w", err)
				return nil, err
			}

			base.TipSet = refreshedBase
		}
	}

	tIntersectAndRefresh := build.Clock.Now()

	log.Infow("createBlock bvals:", "bvals:", bvals)
	minedBlock, err = m.createBlock(base, m.address, ticket, nil, bvals, nil, msgs, nil, 0)
	if err != nil {
		err = xerrors.Errorf("failed to create block: %w", err)
		return nil, err
	}
	///add nil Pqc-Pow-Proof to Header
	minedBlock.Header.PqcPowProof = &types.PqcPowProof{}

	log.Infow("minedBlock:", "Timestamp", time.Unix(int64(minedBlock.Header.Timestamp), 0))
	tCreateBlock := build.Clock.Now()
	dur := tCreateBlock.Sub(tStart)
	parentMiners := make([]address.Address, len(base.TipSet.Blocks()))
	for i, header := range base.TipSet.Blocks() {
		parentMiners[i] = header.Miner
	}

	log.Infow("pqcpow.GetNbit")
	nbit, err := pqcpow.GetNbit(ctx, m.api, minedBlock)
	if err != nil {
		err = xerrors.Errorf("failed GetNbit: %w", err)
		return nil, err
	}
	log.Infow("computePqcPowProof ", "nbit:", nbit)
	minedBlock.Header.PqcPowProof.Nbit = nbit

	// qpcseed, err := minedBlock.Header.PqcSeed()  //TODO:upgrade v2.0
	// qpcseed, err := minedBlock.Header.PqcProofSeed() //TODO:upgrade v2.0
	// if err != nil {
	// 	err = xerrors.Errorf("failed get qpcseed: %w", err)
	// 	return nil, err
	// }

	proofStart := uint64(build.Clock.Now().Unix())
	PqcPowProof, err := computePqcPowProof(ctx, nbit, minedBlock, m.api)
	if err != nil {
		if err == pqcpow.ErrXFoundOutTime {
			log.Warnf("compute pqc pow proof: %+v", err)
			return nil, nil
		}

		if err == pqcpow.NewBlockheads {
			log.Warnf("compute pqc pow proof: %+v", err)
			return nil, nil
		}
		err = xerrors.Errorf("failed computePqcPowProof: %w", err)
		return nil, err
	}
	// winning  post-quantum Proof
	// tPPProof := build.Clock.Now()
	proofEnd := uint64(build.Clock.Now().Unix())

	minedBlock.Header.PqcPowProof = PqcPowProof
	// copy(minedBlock.Header.PqcPowProof.Nbit, PqcPowProof.Nbit)
	// copy(minedBlock.Header.PqcPowProof.PowProof, PqcPowProof.PowProof)
	log.Infow("computePqcPowProof ", "PqcPowProof.Nbit:", PqcPowProof.Nbit, "PqcPowProof.PowProof:", PqcPowProof.PowProof, "len(PqcPowProof.PowProof):", len(PqcPowProof.PowProof))
	log.Infow("computePqcPowProof ", "Header.PqcPowProof.Nbit:", minedBlock.Header.PqcPowProof.Nbit, "Header.PqcPowProof.PowProof:", minedBlock.Header.PqcPowProof.PowProof)
	minedBlock.Header.MiningTimestamp = proofEnd - proofStart

	log.Infow("minedBlock:", "MiningTimestamp", minedBlock.Header.MiningTimestamp)

	log.Infow("mined new block", "cid", minedBlock.Cid(), "height", int64(minedBlock.Header.Height), "miner", minedBlock.Header.Miner, "parents", parentMiners, "parentTipset", base.TipSet.Key().String(), "took", dur)
	if dur > time.Second*time.Duration(build.BlockDelaySecs) || time.Now().Compare(time.Unix(int64(minedBlock.Header.Timestamp), 0)) >= 0 {
		log.Warnw("CAUTION: block production took us past the block time. Your computer may not be fast enough to keep up",
			"tPowercheck ", tPowercheck.Sub(tStart),
			"tTicket ", tTicket.Sub(tPowercheck),
			"tSeed ", tSeed.Sub(tTicket),
			"tProof ", tProof.Sub(tSeed),
			"tPending ", tPending.Sub(tProof),
			"tEquivocateWait ", tEquivocateWait.Sub(tPending),
			"tIntersectAndRefresh ", tIntersectAndRefresh.Sub(tEquivocateWait),
			"tCreateBlock ", tCreateBlock.Sub(tIntersectAndRefresh))
	}

	return minedBlock, nil
}

func (m *PqcMiner) computeTicket(ctx context.Context, brand *types.BeaconEntry, round abi.ChainEpoch, chainRand *types.Ticket, mbi *api.MiningBaseInfo) (*types.Ticket, error) {
	buf := new(bytes.Buffer)
	if err := m.address.MarshalCBOR(buf); err != nil {
		return nil, xerrors.Errorf("failed to marshal address to cbor: %w", err)
	}

	if round > build.UpgradeSmokeHeight {
		buf.Write(chainRand.VRFProof)
	}

	input, err := lrand.DrawRandomnessFromBase(brand.Data, crypto.DomainSeparationTag_TicketProduction, round-build.TicketRandomnessLookback, buf.Bytes())
	if err != nil {
		return nil, err
	}

	// vrfOut, err := gen.ComputeVRF(ctx, m.api.WalletSign, mbi.WorkerKey, input)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Infow("computeTicket vrfOut", "len", len(vrfOut))
	// var sig crypto.Signature
	// if err := sig.UnmarshalCBOR(bytes.NewReader(vrfOut)); err != nil {
	// 	return nil, xerrors.Errorf("unmarshaling error data: %w", err)
	// }
	return &types.Ticket{
		VRFProof: input,
	}, nil
}

// func computePqcPowProofV0(ctx context.Context, seed []byte, nbit []byte, minedBlock *types.BlockMsg, p pqcpow.PqcPowAPI) (*types.PqcPowProof, error) {
// }

func computePqcPowProof(ctx context.Context, nbit []byte, minedBlock *types.BlockMsg, p pqcpow.PqcPowAPI) (*types.PqcPowProof, error) {
	var seedProof []byte
	nbitProof := nbit

	// lastTipSet, err := p.ChainGetTipSet(ctx, types.NewTipSetKey(minedBlock.Header.Parents...))
	// if err != nil {
	// 	return nil, err
	// }

	// if build.UpgradeYellowStoneHeight >= 0 && lastTipSet.Height() >= build.UpgradeYellowStoneHeight {
	seedProof, err := minedBlock.Header.PqcProofSeed() //TODO:upgrade v2.0
	if err != nil {
		err = xerrors.Errorf("failed get qpcseed: %w", err)
		return nil, err
	}
	tm := time.NewTicker(17 * time.Second)
	for {
		seedBuf := shake3.Shake256XOF(seedProof, 72)
		log.Infow("computePqcPowProof", " seedBuf:", seedBuf)

		qproof, err := pqcpow.PqcPowProof(ctx, seedBuf, nbitProof, p, tm)
		if err == pqcpow.ErrXNotFound {
			log.Infow("computePqcPowProof RefreshNbit")
			nbitProof = pqcpow.RefreshNbit(nbitProof) //

			minedBlock.Header.PqcPowProof.Nbit = nbitProof

			qpcseed, err := minedBlock.Header.PqcProofSeed() //TODO:upgrade v2.0
			if err != nil {
				err = xerrors.Errorf("failed get qpcseed: %w", err)
				return nil, err
			}

			seedProof = qpcseed
			continue
		}

		if err != nil {
			return nil, err
		}

		fmt.Println("computePqcPowProof nbitProof:  qproof: len(qproof):", nbitProof, qproof, len(qproof))
		return &types.PqcPowProof{
			Nbit:     nbitProof,
			PowProof: qproof,
		}, nil
	}
	// }
	// seedProof, err = minedBlock.Header.PqcSeed()

	// seedBuf := shake3.Shake256XOF(seedProof, 72)
	// log.Infow("computePqcPowProof", " seedBuf:", seedBuf)
	// qproof, err := pqcpow.PqcPowProof(ctx, seedBuf, nbitProof, p)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("computePqcPowProof nbitProof:  qproof: len(qproof):", nbitProof, qproof, len(qproof))
	// return &types.PqcPowProof{
	// 	Nbit:     nbitProof,
	// 	PowProof: qproof,
	// }, nil
}

func (m *PqcMiner) createBlock(base *MiningBase, addr address.Address, ticket *types.Ticket,
	eproof *types.ElectionProof, bvals []types.BeaconEntry, wpostProof []proof.PoStProof,
	msgs []*types.SignedMessage, pproof *types.PqcPowProof, mts uint64) (*types.BlockMsg, error) {
	uts := base.TipSet.MinTimestamp() + build.BlockDelaySecs*(uint64(base.NullRounds)+1)
	// log.Infow("createBlock:", "mts", time.Unix(int64(mts), 0))
	nheight := base.TipSet.Height() + base.NullRounds + 1

	// why even return this? that api call could just submit it for us
	return m.api.MinerCreateBlock(context.TODO(), &api.BlockTemplate{
		Miner:            addr,
		Parents:          base.TipSet.Key(),
		Ticket:           ticket,
		Eproof:           eproof,
		BeaconValues:     bvals,
		Messages:         msgs,
		Epoch:            nheight,
		Timestamp:        uts,
		MiningTimestamp:  mts,
		WinningPoStProof: wpostProof,
		Pproof:           pproof,
	})
}

func (m *PqcMiner) getX(devID int, startSMCount int) ([]byte, error) {
	return nil, nil
}

// async getX(devID: number = 0, startSMCount: number = 0, cb?: Function): Promise<any> {
// 	this.setDiviceID(devID);
// 	this.startSMCount = startSMCount;
// 	let fix = this.minerController.getNextFixStr();

// 	while (true) {
// 		let x;
// 		if (this.minerController.fixNumber > 0) { //do fix
// 			if (fix) {
// 				x = await this.cal(fix);
// 			} else {
// 				if (!this.minerController.otherMachineRunning(devID)) {
// 					if (cb) cb(false);
// 				}
// 				this.minerController.setMachineStoping(devID);
// 				return false;
// 			}

// 			if (x.errCode === -101 || x.err || !x.result) {
// 				console.log(`Fix str '${fix}' not found.`);
// 				this.startSMCount = 0;
// 				fix = this.minerController.getNextFixStr();
// 				continue;
// 			}

// 			if (!this.mqphash.checkIsSolution(x.result.xBuf.subarray(0, this.mqphash.MQP.variablesByte))) {
// 				console.log(`Fix str '${fix}' check solution failed.`);
// 				this.startSMCount = 0;
// 				fix = this.minerController.getNextFixStr();
// 				continue;
// 			}
// 		}
// 		else { //no fix
// 			x = await this.cal();

// 			if (x.errCode === -101 || x.err || !x.result) {
// 				console.log('Get solution failed!');
// 				if (cb) cb(false);
// 				return false;
// 			}

// 			if (!this.checkSolution()) {
// 				console.log('Check solution failed!');
// 				if (cb) cb(false);
// 				return false;
// 			}
// 		}

// 		if (verifyPoW(this.mqphash.MQP.seed, this.nbit, x.result.xBuf)) {
// 			if (cb) cb(x.result.xBuf);
// 			return x.result.xBuf;
// 		}
// 		this.startSMCount = Number(x.result.smCount) + 1 || 0;
// 	}
// }

func (m *PqcMiner) cal(fix []byte) ([]byte, error) {
	if len(fix) != 0 {

	} else {

	}

	return nil, nil
}

// async cal(fix?: string): Promise<{ err?: string, errCode?: number, result?: any }> {
// 	return new Promise(async (r) => {
// 		if (this._minerController.minerBinPath == undefined) {
// 			console.log('no path!');
// 			r({ err: 'no path!' });
// 			return;
// 		}

// 		let minerBin = await fs.readFile(this._minerController.minerBinPath);
// 		let mbSha = shake256(minerBin);
// 		if (!mbSha.equals(this._minerController.minerBinHash)) {
// 			console.log('minerBinHash was change!');
// 			r({ err: 'minerBinHash was change!' });
// 			return;
// 		}

// 		let equations;
// 		let opts;
// 		let mf: MinerFix;
// 		if (fix) {
// 			mf = new MinerFix(this.mqphash, fix.length);
// 			equations = this.mqphash.MQP.equations.map(x => mf.fixOneEquation(fix, x.toString('hex'), this.mqphash.MQP.unwantedCoefficientBit).newCoeBuf.toString('hex'));
// 			opts = ['2', `${this.deviceID}`, `${this.m}`, `${mf.newN}`, `${this.whichXWidth}`, `${this.startSMCount}`, `${mf.newCoe}`];
// 		}
// 		else {
// 			equations = this.mqphash.MQP.equations.map(x => x.toString('hex'));
// 			opts = ['2', `${this.deviceID}`, `${this.m}`, `${this.n}`, `${this.whichXWidth}`, `${this.startSMCount}`, `${this.mqphash.MQP.coefficient}`];
// 		}
// 		let str = '';
// 		this.child = spawn(this._minerController.minerBinPath, opts);
// 		this.child.stdin.write(equations.join('\n'));
// 		this.child.stdin.write('\nend\n');

// 		this.child.stdout.on('data', (data) => {
// 			str += data.toString();
// 		});

// 		this.child.on('exit', (code) => {
// 			if (code == 0) {
// 				let words = str.toString().split('x found:');
// 				if (words[1] == undefined) {
// 					r({ errCode: -101 });
// 				}
// 				else {
// 					//fix back
// 					this.x_data = JSON.parse(words[1]);
// 					if (fix) {
// 						this.x_data.xBuf = mf.fixBack(this.x_data.x, fix);
// 					}
// 					r({ result: this.x_data });
// 				}
// 			}
// 		});

// 		this.child.on('close', (code, signal) => {
// 			if (signal === 'SIGTERM') {
// 				r({ err: 'child process terminated due to receipt of signal SIGTERM!' });
// 			} else {
// 				r({ err: 'child process terminated!' });
// 			}
// 		});

// 		this.child.stderr.on('data', (data) => {
// 			console.error(`stderr: ${data}`);
// 		});
// 	});
// }

// private checkSolution(): boolean {
// 	let solution = this.x_data.x;
// 	let x = solution.slice(solution.length - this.n, solution.length);

// 	for (let index = 0; index < this.mqphash.MQP.unwantedVariablesBit; index++) {
// 		x += '0';
// 	}

// 	let xBuf = Buffer.alloc(32);
// 	let index = 0;

// 	for (let i = 0; i < x.length; i += 8) {
// 		xBuf[index++] = parseInt(x.slice(i, i + 8), 2);
// 	}
// 	this.x_data.xBuf = xBuf;

// 	return this.mqphash.checkIsSolution(xBuf.subarray(0, this.mqphash.MQP.variablesByte));
// }
