package genesis

import (
	"context"

	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	actorstypes "github.com/post-quantumqoin/core-types/actors"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/builtin"
	"github.com/post-quantumqoin/core-types/manifest"
	"github.com/post-quantumqoin/specs-contracts/contracts/util/adt"

	bstore "github.com/post-quantumqoin/qoin-shor/blockstore"
	"github.com/post-quantumqoin/qoin-shor/core/actors"
	"github.com/post-quantumqoin/qoin-shor/core/actors/builtin/datacap"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

var GovernorId address.Address

func init() {
	idk, err := address.NewFromString("t06")
	if err != nil {
		panic(err)
	}

	GovernorId = idk
}

func SetupDatacapActor(ctx context.Context, bs bstore.Blockstore, av actorstypes.Version) (*types.Actor, error) {
	cst := cbor.NewCborStore(bs)
	dst, err := datacap.MakeState(adt.WrapStore(ctx, cbor.NewCborStore(bs)), av, GovernorId, builtin.DefaultTokenActorBitwidth)
	if err != nil {
		return nil, err
	}

	statecid, err := cst.Put(ctx, dst.GetState())
	if err != nil {
		return nil, err
	}

	actcid, ok := actors.GetActorCodeID(av, manifest.DatacapKey)
	if !ok {
		return nil, xerrors.Errorf("failed to get datacap actor code ID for actors version %d", av)
	}

	act := &types.Actor{
		Code:    actcid,
		Head:    statecid,
		Balance: big.Zero(),
	}

	return act, nil
}
