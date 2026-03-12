package cron

import (
	"fmt"

	"github.com/ipfs/go-cid"

	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	cron2 "github.com/post-quantumqoin/specs-contracts/contracts/builtin/cron"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state2)(nil)

func load2(store adt.Store, root cid.Cid) (State, error) {
	out := state2{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make2(store adt.Store) (State, error) {
	out := state2{store: store}
	out.State = *cron2.ConstructState(cron2.BuiltInEntries())
	return &out, nil
}

type state2 struct {
	cron2.State
	store adt.Store
}

func (s *state2) GetState() interface{} {
	return &s.State
}

func (s *state2) ActorKey() string {
	return manifest.CronKey
}

func (s *state2) ActorVersion() actorstypes.Version {
	return actorstypes.Version2
}

func (s *state2) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
