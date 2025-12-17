package vm

import (
	"context"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/post-quantumqoin/address"
	actorstypes "github.com/post-quantumqoin/core-types/actors"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/exitcode"
	"github.com/post-quantumqoin/core-types/network"
	builtin0 "github.com/post-quantumqoin/specs-contracts/contracts/builtin"
	builtin2 "github.com/post-quantumqoin/specs-contracts/contracts/builtin"
	builtin3 "github.com/post-quantumqoin/specs-contracts/contracts/builtin"
	builtin4 "github.com/post-quantumqoin/specs-contracts/contracts/builtin"
	builtin5 "github.com/post-quantumqoin/specs-contracts/contracts/builtin"
	builtin6 "github.com/post-quantumqoin/specs-contracts/contracts/builtin"
	builtin7 "github.com/post-quantumqoin/specs-contracts/contracts/builtin"

	"github.com/post-quantumqoin/qoin-shor/build"
	"github.com/post-quantumqoin/qoin-shor/core/actors"
	"github.com/post-quantumqoin/qoin-shor/core/actors/aerrors"
	"github.com/post-quantumqoin/qoin-shor/core/actors/builtin"
	"github.com/post-quantumqoin/qoin-shor/core/actors/builtin/account"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

func init() {
	cst := cbor.NewMemCborStore()
	emptyobject, err := cst.Put(context.TODO(), []struct{}{})
	if err != nil {
		panic(err)
	}

	EmptyObjectCid = emptyobject
}

var EmptyObjectCid cid.Cid

// TryCreateAccountActor creates account actors from only BLS/SECP256K1 addresses.
func TryCreateAccountActor(rt *Runtime, addr address.Address) (*types.Actor, address.Address, aerrors.ActorError) {
	if err := rt.chargeGasSafe(PricelistByEpoch(rt.height).OnCreateActor()); err != nil {
		return nil, address.Undef, err
	}

	if addr == build.ZeroAddress && rt.NetworkVersion() >= network.Version10 {
		return nil, address.Undef, aerrors.New(exitcode.ErrIllegalArgument, "cannot create the zero bls actor")
	}

	addrID, err := rt.state.RegisterNewAddress(addr)
	if err != nil {
		return nil, address.Undef, aerrors.Escalate(err, "registering actor address")
	}

	av, err := actorstypes.VersionForNetwork(rt.NetworkVersion())
	if err != nil {
		return nil, address.Undef, aerrors.Escalate(err, "unsupported network version")
	}

	act, aerr := makeAccountActor(av, addr)
	if aerr != nil {
		return nil, address.Undef, aerr
	}

	if err := rt.state.SetActor(addrID, act); err != nil {
		return nil, address.Undef, aerrors.Escalate(err, "creating new actor failed")
	}

	p, err := actors.SerializeParams(&addr)
	if err != nil {
		return nil, address.Undef, aerrors.Escalate(err, "couldn't serialize params for actor construction")
	}
	// call constructor on account

	_, aerr = rt.internalSend(builtin.SystemActorAddr, addrID, account.Methods.Constructor, big.Zero(), p)
	if aerr != nil {
		return nil, address.Undef, aerrors.Wrap(aerr, "failed to invoke account constructor")
	}

	act, err = rt.state.GetActor(addrID)
	if err != nil {
		return nil, address.Undef, aerrors.Escalate(err, "loading newly created actor failed")
	}
	return act, addrID, nil
}

func makeAccountActor(ver actorstypes.Version, addr address.Address) (*types.Actor, aerrors.ActorError) {
	switch addr.Protocol() {
	case address.BLS, address.SECP256K1, address.PQC:
		return newAccountActor(ver, addr), nil
	case address.ID:
		return nil, aerrors.Newf(exitcode.SysErrInvalidReceiver, "no actor with given ID: %s", addr)
	case address.Actor:
		return nil, aerrors.Newf(exitcode.SysErrInvalidReceiver, "no such actor: %s", addr)
	default:
		return nil, aerrors.Newf(exitcode.SysErrInvalidReceiver, "address has unsupported protocol: %d", addr.Protocol())
	}
}

func newAccountActor(ver actorstypes.Version, addr address.Address) *types.Actor {
	// TODO: ActorsUpgrade use a global actor registry?
	var code cid.Cid
	switch ver {
	case actorstypes.Version0:
		code = builtin0.AccountActorCodeID
	case actorstypes.Version2:
		code = builtin2.AccountActorCodeID
	case actorstypes.Version3:
		code = builtin3.AccountActorCodeID
	case actorstypes.Version4:
		code = builtin4.AccountActorCodeID
	case actorstypes.Version5:
		code = builtin5.AccountActorCodeID
	case actorstypes.Version6:
		code = builtin6.AccountActorCodeID
	case actorstypes.Version7:
		code = builtin7.AccountActorCodeID
	default:
		panic("unsupported actors version")
	}
	nact := &types.Actor{
		Code:    code,
		Balance: types.NewInt(0),
		Head:    EmptyObjectCid,
		Address: &addr,
	}

	return nact
}
