package genesis

import (
	"context"

	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/core-types/big"
	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	"github.com/post-quantumqoin/specs-contracts/contracts/util/adt"

	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	"github.com/post-quantumqoin/qoin-shor/core/contracts/builtin/power"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	bstore "github.com/post-quantumqoin/qoin-shor/dbstore"
)

func SetupStoragePowerActor(ctx context.Context, bs bstore.Blockstore, av actorstypes.Version) (*types.Actor, error) {

	cst := cbor.NewCborStore(bs)
	pst, err := power.MakeState(adt.WrapStore(ctx, cbor.NewCborStore(bs)), av)
	if err != nil {
		return nil, err
	}

	statecid, err := cst.Put(ctx, pst.GetState())
	if err != nil {
		return nil, err
	}

	actcid, ok := actors.GetActorCodeID(av, manifest.PowerKey)
	if !ok {
		return nil, xerrors.Errorf("failed to get power actor code ID for actors version %d", av)
	}

	act := &types.Actor{
		Code:    actcid,
		Head:    statecid,
		Balance: big.Zero(),
	}

	return act, nil
}
