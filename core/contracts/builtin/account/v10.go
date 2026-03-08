package account

import (
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/post-quantumqoin/address"
	account10 "github.com/post-quantumqoin/core-types/builtin/v10/account"
	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state10)(nil)

func load10(store adt.Store, root cid.Cid) (State, error) {
	out := state10{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make10(store adt.Store, addr address.Address) (State, error) {
	out := state10{store: store}
	out.State = account10.State{Address: addr}
	return &out, nil
}

type state10 struct {
	account10.State
	store adt.Store
}

func (s *state10) PubkeyAddress() (address.Address, error) {
	return s.Address, nil
}

func (s *state10) GetState() interface{} {
	return &s.State
}

func (s *state10) ActorKey() string {
	return manifest.AccountKey
}

func (s *state10) ActorVersion() actorstypes.Version {
	return actorstypes.Version10
}

func (s *state10) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
