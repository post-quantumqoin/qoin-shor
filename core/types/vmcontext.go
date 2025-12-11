package types

import (
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/post-quantumqoin/address"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/aerrors"
)

type Storage interface {
	Put(cbg.CBORMarshaler) (cid.Cid, aerrors.ActorError)
	Get(cid.Cid, cbg.CBORUnmarshaler) aerrors.ActorError

	GetHead() cid.Cid

	// Commit sets the new head of the actors state as long as the current
	// state matches 'oldh'
	Commit(oldh cid.Cid, newh cid.Cid) aerrors.ActorError
}

type StateTree interface {
	SetActor(addr address.Address, act *Actor) error
	// GetActor returns the actor from any type of `addr` provided.
	GetActor(addr address.Address) (*Actor, error)

	Version() StateTreeVersion
}
