package filcns

import (
	"context"
	"fmt"
	"math/big"

	bstore "github.com/post-quantumqoin/qoin-shor/blockstore"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

var zero = types.NewInt(0)

func Weight(ctx context.Context, stateBs bstore.Blockstore, ts *types.TipSet) (types.BigInt, error) {
	if ts == nil {
		fmt.Printf("Weight ts == nil\n")
		return types.NewInt(0), nil
	}
	// >>> w[r] <<< + wFunction(totalPowerAtTipset(ts)) * 2^8 + (wFunction(totalPowerAtTipset(ts)) * sum(ts.blocks[].ElectionProof.WinCount) * wRatio_num * 2^8) / (e * wRatio_den)

	var out = new(big.Int).Set(ts.ParentWeight().Int)

	// >>> wFunction(totalPowerAtTipset(ts)) * 2^8 <<< + (wFunction(totalPowerAtTipset(ts)) * sum(ts.blocks[].ElectionProof.WinCount) * wRatio_num * 2^8) / (e * wRatio_den)

	// var tpow big2.Int
	// {
	// 	cst := cbor.NewCborStore(stateBs)
	// 	state, err := state.LoadStateTree(cst, ts.ParentState())
	// 	if err != nil {
	// 		return types.NewInt(0), xerrors.Errorf("load state tree: %w", err)
	// 	}

	// 	act, err := state.GetActor(power.Address)
	// 	if err != nil {
	// 		return types.NewInt(0), xerrors.Errorf("get power actor: %w", err)
	// 	}

	// 	powState, err := power.Load(store.ActorStore(ctx, stateBs), act)
	// 	if err != nil {
	// 		return types.NewInt(0), xerrors.Errorf("failed to load power actor state: %w", err)
	// 	}

	// 	claim, err := powState.TotalPower()
	// 	if err != nil {
	// 		return types.NewInt(0), xerrors.Errorf("failed to get total power: %w", err)
	// 	}

	// 	tpow = claim.QualityAdjPower // TODO: REVIEW: Is this correct?
	// 	fmt.Printf("Weight tpow:%d\n", claim.QualityAdjPower)
	// }

	// log2P := int64(0)
	// if tpow.GreaterThan(zero) {
	// 	log2P = int64(tpow.BitLen() - 1)
	// } else {
	// 	// Not really expect to be here ...
	// 	return types.EmptyInt, xerrors.Errorf("All power in the net is gone. You network might be disconnected, or the net is dead!")
	// }

	// out.Add(out, big.NewInt(log2P<<8))

	// (wFunction(totalPowerAtTipset(ts)) * sum(ts.blocks[].ElectionProof.WinCount) * wRatio_num * 2^8) / (e * wRatio_den)

	// totalJ := int64(0)
	// for _, b := range ts.Blocks() {
	// 	totalJ += b.ElectionProof.WinCount
	// }

	// eWeight := big.NewInt((log2P * build.WRatioNum))
	// eWeight = eWeight.Lsh(eWeight, 8)
	// eWeight = eWeight.Mul(eWeight, new(big.Int).SetInt64(totalJ))
	// eWeight = eWeight.Div(eWeight, big.NewInt(int64(build.BlocksPerEpoch*build.WRatioDen)))

	// out = out.Add(out, eWeight)
	//weight = ParentWeight + sum(headers)
	out = out.Add(out, new(big.Int).SetInt64(int64(len(ts.Blocks()))))
	fmt.Printf("Weight out:%d\n", out)
	return types.BigInt{Int: out}, nil
}
