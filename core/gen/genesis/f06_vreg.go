package genesis

import (
	"context"

	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	actorstypes "github.com/post-quantumqoin/core-types/actors"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/manifest"
	"github.com/post-quantumqoin/specs-contracts/contracts/util/adt"

	bstore "github.com/post-quantumqoin/qoin-shor/blockstore"
	"github.com/post-quantumqoin/qoin-shor/core/actors"
	"github.com/post-quantumqoin/qoin-shor/core/actors/builtin/verifreg"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

var RootVerifierID address.Address

func init() {

	idk, err := address.NewFromString("t080")
	if err != nil {
		panic(err)
	}

	RootVerifierID = idk
}

func SetupVerifiedRegistryActor(ctx context.Context, bs bstore.Blockstore, av actorstypes.Version) (*types.Actor, error) {
	cst := cbor.NewCborStore(bs)
	vst, err := verifreg.MakeState(adt.WrapStore(ctx, cbor.NewCborStore(bs)), av, RootVerifierID)
	if err != nil {
		return nil, err
	}

	statecid, err := cst.Put(ctx, vst.GetState())
	if err != nil {
		return nil, err
	}

	actcid, ok := actors.GetActorCodeID(av, manifest.VerifregKey)
	if !ok {
		return nil, xerrors.Errorf("failed to get verifreg actor code ID for actors version %d", av)
	}

	act := &types.Actor{
		Code:    actcid,
		Head:    statecid,
		Balance: big.Zero(),
	}

	return act, nil
}
