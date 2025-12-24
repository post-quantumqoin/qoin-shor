package genesis

import (
	"context"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	actorstypes "github.com/post-quantumqoin/core-types/actors"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/manifest"

	bstore "github.com/post-quantumqoin/qoin-shor/blockstore"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
	"github.com/post-quantumqoin/qoin-shor/core/contracts/builtin/system"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

func SetupSystemActor(ctx context.Context, bs bstore.Blockstore, av actorstypes.Version) (*types.Actor, error) {

	cst := cbor.NewCborStore(bs)
	// TODO pass in built-in actors cid for V8 and later
	st, err := system.MakeState(adt.WrapStore(ctx, cst), av, cid.Undef)
	if err != nil {
		return nil, err
	}

	// if av >= actorstypes.Version8 {
	mfCid, ok := actors.GetManifest(av)
	if !ok {
		return nil, xerrors.Errorf("missing manifest for actors version %d", av)
	}

	mf := manifest.Manifest{}
	if err := cst.Get(ctx, mfCid, &mf); err != nil {
		return nil, xerrors.Errorf("loading manifest for actors version %d: %w", av, err)
	}

	if err := st.SetBuiltinActors(mf.Data); err != nil {
		return nil, xerrors.Errorf("failed to set manifest data: %w", err)
	}
	// }

	statecid, err := cst.Put(ctx, st.GetState())
	if err != nil {
		return nil, err
	}

	actcid, ok := actors.GetActorCodeID(av, manifest.SystemKey)
	if !ok {
		return nil, xerrors.Errorf("failed to get system actor code ID for actors version %d", av)
	}

	act := &types.Actor{
		Code:    actcid,
		Head:    statecid,
		Balance: big.Zero(),
	}

	return act, nil
}
