package modules

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/fx"

	"github.com/post-quantumqoin/core-types/abi"

	"github.com/post-quantumqoin/qoin-shor/core/ethhashlookup"
	"github.com/post-quantumqoin/qoin-shor/core/events"
	"github.com/post-quantumqoin/qoin-shor/core/messagepool"
	"github.com/post-quantumqoin/qoin-shor/core/stmgr"
	"github.com/post-quantumqoin/qoin-shor/core/store"
	"github.com/post-quantumqoin/qoin-shor/node/config"
	"github.com/post-quantumqoin/qoin-shor/node/impl/full"
	"github.com/post-quantumqoin/qoin-shor/node/modules/helpers"
	"github.com/post-quantumqoin/qoin-shor/node/repo"
)

func EthModuleAPI(cfg config.FevmConfig) func(helpers.MetricsCtx, repo.LockedRepo, fx.Lifecycle, *store.ChainStore, *stmgr.StateManager, EventHelperAPI, *messagepool.MessagePool, full.StateAPI, full.ChainAPI, full.MpoolAPI, full.SyncAPI) (*full.EthModule, error) {
	return func(mctx helpers.MetricsCtx, r repo.LockedRepo, lc fx.Lifecycle, cs *store.ChainStore, sm *stmgr.StateManager, evapi EventHelperAPI, mp *messagepool.MessagePool, stateapi full.StateAPI, chainapi full.ChainAPI, mpoolapi full.MpoolAPI, syncapi full.SyncAPI) (*full.EthModule, error) {
		sqlitePath, err := r.SqlitePath()
		if err != nil {
			return nil, err
		}

		dbPath := filepath.Join(sqlitePath, "txhash.db")

		// Check if the db exists, if not, we'll back-fill some entries
		_, err = os.Stat(dbPath)
		dbAlreadyExists := err == nil

		transactionHashLookup, err := ethhashlookup.NewTransactionHashLookup(dbPath)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return transactionHashLookup.Close()
			},
		})

		ethTxHashManager := full.EthTxHashManager{
			StateAPI:              stateapi,
			TransactionHashLookup: transactionHashLookup,
		}

		if !dbAlreadyExists {
			err = ethTxHashManager.PopulateExistingMappings(mctx, 0)
			if err != nil {
				return nil, err
			}
		}

		// prefill the whole skiplist cache maintained internally by the GetTipsetByHeight
		go func() {
			start := time.Now()
			log.Infoln("Start prefilling GetTipsetByHeight cache")
			_, err := cs.GetTipsetByHeight(mctx, abi.ChainEpoch(0), cs.GetHeaviestTipSet(), false)
			if err != nil {
				log.Warnf("error when prefilling GetTipsetByHeight cache: %w", err)
			}
			log.Infof("Prefilling GetTipsetByHeight done in %s", time.Since(start))
		}()

		ctx := helpers.LifecycleCtx(mctx, lc)
		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				ev, err := events.NewEvents(ctx, &evapi)
				if err != nil {
					return err
				}

				// Tipset listener
				_ = ev.Observe(&ethTxHashManager)

				ch, err := mp.Updates(ctx)
				if err != nil {
					return err
				}
				go full.WaitForMpoolUpdates(ctx, ch, &ethTxHashManager)
				go full.EthTxHashGC(ctx, cfg.EthTxHashMappingLifetimeDays, &ethTxHashManager)

				return nil
			},
		})

		return &full.EthModule{
			Chain:        cs,
			Mpool:        mp,
			StateManager: sm,

			ChainAPI: chainapi,
			MpoolAPI: mpoolapi,
			StateAPI: stateapi,
			SyncAPI:  syncapi,

			EthTxHashManager: &ethTxHashManager,
		}, nil
	}
}
