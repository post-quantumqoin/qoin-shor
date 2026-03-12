package account

import (
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/post-quantumqoin/address"
	account12 "github.com/post-quantumqoin/core-types/builtin/v12/account"
	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state12)(nil)

func load12(store adt.Store, root cid.Cid) (State, error) {
	out := state12{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make12(store adt.Store, addr address.Address) (State, error) {
	out := state12{store: store}
	out.State = account12.State{Address: addr}
	return &out, nil
}

type state12 struct {
	account12.State
	store adt.Store
}

func (s *state12) PubkeyAddress() (address.Address, error) {
	return s.Address, nil
}

func (s *state12) GetState() interface{} {
	return &s.State
}

func (s *state12) ActorKey() string {
	return manifest.AccountKey
}

func (s *state12) ActorVersion() actorstypes.Version {
	return actorstypes.Version12
}

func (s *state12) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
