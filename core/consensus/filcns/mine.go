package filcns

import (
	"context"
	"fmt"

	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/core/consensus"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

func (filec *FilecoinEC) CreateBlock(ctx context.Context, w api.Wallet, bt *api.BlockTemplate) (*types.FullBlock, error) {
	fmt.Println("QoinEC CreateBlock")
	pts, err := filec.sm.ChainStore().LoadTipSet(ctx, bt.Parents)
	if err != nil {
		return nil, xerrors.Errorf("failed to load parent tipset: %w", err)
	}

	// _, lbst, err := stmgr.GetLookbackTipSetForRound(ctx, filec.sm, pts, bt.Epoch)
	// if err != nil {
	// 	return nil, xerrors.Errorf("getting lookback miner actor state: %w", err)
	// }

	// worker, err := stmgr.GetMinerWorkerRaw(ctx, filec.sm, lbst, bt.Miner)
	// if err != nil {
	// 	return nil, xerrors.Errorf("failed to get miner worker: %w", err)
	// }
	fmt.Println("QoinEC CreateBlock CreateBlockHeader")
	next, blsMessages, secpkMessages, err := consensus.CreateBlockHeader(ctx, filec.sm, pts, bt)
	if err != nil {
		return nil, xerrors.Errorf("failed to process messages from block template: %w", err)
	}
	fmt.Println("QoinEC CreateBlock signBlock")
	// if err := signBlock(ctx, w, worker, next); err != nil {
	// 	return nil, xerrors.Errorf("failed to sign new block: %w", err)
	// }

	fullBlock := &types.FullBlock{
		Header:        next,
		BlsMessages:   blsMessages,
		SecpkMessages: secpkMessages,
	}

	return fullBlock, nil
}
