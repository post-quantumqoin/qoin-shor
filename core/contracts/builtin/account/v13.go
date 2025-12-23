package account

import (
	"fmt"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/lotus/chain/actors"
	"github.com/post-quantumqoin/address"
	actorstypes "github.com/post-quantumqoin/core-types/actors"
	account13 "github.com/post-quantumqoin/core-types/builtin/v13/account"
	"github.com/post-quantumqoin/core-types/manifest"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state13)(nil)

func load13(store adt.Store, root cid.Cid) (State, error) {
	out := state13{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make13(store adt.Store, addr address.Address) (State, error) {
	out := state13{store: store}
	out.State = account13.State{Address: addr, Cert: nil}
	return &out, nil
}

type state13 struct {
	account13.State
	store adt.Store
}

func (s *state13) PubkeyAddress() (address.Address, error) {
	return s.Address, nil
}

func (s *state13) GetState() interface{} {
	return &s.State
}

func (s *state13) ActorKey() string {
	return manifest.AccountKey
}

func (s *state13) ActorVersion() actorstypes.Version {
	return actorstypes.Version13
}

func (s *state13) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
