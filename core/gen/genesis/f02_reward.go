package genesis

import (
	"context"

	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	actorstypes "github.com/post-quantumqoin/core-types/actors"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/manifest"

	bstore "github.com/post-quantumqoin/qoin-shor/blockstore"
	"github.com/post-quantumqoin/qoin-shor/build"
	"github.com/post-quantumqoin/qoin-shor/core/actors"
	"github.com/post-quantumqoin/qoin-shor/core/actors/adt"
	"github.com/post-quantumqoin/qoin-shor/core/actors/builtin/reward"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

func SetupRewardActor(ctx context.Context, bs bstore.Blockstore, qaPower big.Int, av actorstypes.Version) (*types.Actor, error) {
	cst := cbor.NewCborStore(bs)
	rst, err := reward.MakeState(adt.WrapStore(ctx, cst), av, qaPower)
	if err != nil {
		return nil, err
	}

	statecid, err := cst.Put(ctx, rst.GetState())
	if err != nil {
		return nil, err
	}

	actcid, ok := actors.GetActorCodeID(av, manifest.RewardKey)
	if !ok {
		return nil, xerrors.Errorf("failed to get reward actor code ID for actors version %d", av)
	}

	act := &types.Actor{
		Code:    actcid,
		Balance: types.BigInt{Int: build.InitialRewardBalance},
		Head:    statecid,
	}

	return act, nil
}
