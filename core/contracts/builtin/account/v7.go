package account

import (
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/post-quantumqoin/address"
	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	account7 "github.com/post-quantumqoin/specs-contracts/contracts/builtin/account"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state7)(nil)

func load7(store adt.Store, root cid.Cid) (State, error) {
	out := state7{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make7(store adt.Store, addr address.Address) (State, error) {
	out := state7{store: store}
	out.State = account7.State{Address: addr}
	return &out, nil
}

type state7 struct {
	account7.State
	store adt.Store
}

func (s *state7) PubkeyAddress() (address.Address, error) {
	return s.Address, nil
}

func (s *state7) GetState() interface{} {
	return &s.State
}

func (s *state7) ActorKey() string {
	return manifest.AccountKey
}

func (s *state7) ActorVersion() actorstypes.Version {
	return actorstypes.Version7
}

func (s *state7) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
