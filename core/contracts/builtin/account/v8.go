package account

import (
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/post-quantumqoin/address"
	account8 "github.com/post-quantumqoin/core-types/builtin/v8/account"
	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state8)(nil)

func load8(store adt.Store, root cid.Cid) (State, error) {
	out := state8{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make8(store adt.Store, addr address.Address) (State, error) {
	out := state8{store: store}
	out.State = account8.State{Address: addr}
	return &out, nil
}

type state8 struct {
	account8.State
	store adt.Store
}

func (s *state8) PubkeyAddress() (address.Address, error) {
	return s.Address, nil
}

func (s *state8) GetState() interface{} {
	return &s.State
}

func (s *state8) ActorKey() string {
	return manifest.AccountKey
}

func (s *state8) ActorVersion() actorstypes.Version {
	return actorstypes.Version8
}

func (s *state8) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
