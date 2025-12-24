package genesis

import (
	"context"

	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	actorstypes "github.com/post-quantumqoin/core-types/actors"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/manifest"

	bstore "github.com/post-quantumqoin/qoin-shor/blockstore"
	"github.com/post-quantumqoin/qoin-shor/core/actors"
	"github.com/post-quantumqoin/qoin-shor/core/actors/adt"
	"github.com/post-quantumqoin/qoin-shor/core/actors/builtin/cron"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

func SetupCronActor(ctx context.Context, bs bstore.Blockstore, av actorstypes.Version) (*types.Actor, error) {
	cst := cbor.NewCborStore(bs)
	st, err := cron.MakeState(adt.WrapStore(ctx, cbor.NewCborStore(bs)), av)
	if err != nil {
		return nil, err
	}

	statecid, err := cst.Put(ctx, st.GetState())
	if err != nil {
		return nil, err
	}

	actcid, ok := actors.GetActorCodeID(av, manifest.CronKey)
	if !ok {
		return nil, xerrors.Errorf("failed to get cron actor code ID for actors version %d", av)
	}

	act := &types.Actor{
		Code:    actcid,
		Head:    statecid,
		Balance: big.Zero(),
	}

	return act, nil
}
