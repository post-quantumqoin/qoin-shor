package qoincns

import (
	"context"
	"fmt"
	"math/big"

	"github.com/post-quantumqoin/qoin-shor/core/types"
	bstore "github.com/post-quantumqoin/qoin-shor/dbstore"
)

var zero = types.NewInt(0)

func Weight(ctx context.Context, stateBs bstore.Blockstore, ts *types.TipSet) (types.BigInt, error) {
	if ts == nil {
		fmt.Printf("Weight ts == nil\n")
		return types.NewInt(0), nil
	}

	var out = new(big.Int).Set(ts.ParentWeight().Int)

 	out = out.Add(out, new(big.Int).SetInt64(int64(len(ts.Blocks()))))
	fmt.Printf("Weight out:%d\n", out)
	return types.BigInt{Int: out}, nil
}
