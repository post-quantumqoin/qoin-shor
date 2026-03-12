package cron

import (
	"fmt"

	"github.com/ipfs/go-cid"

	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	cron0 "github.com/post-quantumqoin/specs-contracts/contracts/builtin/cron"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state0)(nil)

func load0(store adt.Store, root cid.Cid) (State, error) {
	out := state0{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make0(store adt.Store) (State, error) {
	out := state0{store: store}
	out.State = *cron0.ConstructState(cron0.BuiltInEntries())
	return &out, nil
}

type state0 struct {
	cron0.State
	store adt.Store
}

func (s *state0) GetState() interface{} {
	return &s.State
}

func (s *state0) ActorKey() string {
	return manifest.CronKey
}

func (s *state0) ActorVersion() actorstypes.Version {
	return actorstypes.Version0
}

func (s *state0) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
