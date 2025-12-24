package cron

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/post-quantumqoin/qoin-shor/core/actors"

	actorstypes "github.com/post-quantumqoin/core-types/actors"
	cron13 "github.com/post-quantumqoin/core-types/builtin/v13/cron"
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

func make13(store adt.Store) (State, error) {
	out := state13{store: store}
	out.State = *cron13.ConstructState(cron13.BuiltInEntries())
	return &out, nil
}

type state13 struct {
	cron13.State
	store adt.Store
}

func (s *state13) GetState() interface{} {
	return &s.State
}

func (s *state13) ActorKey() string {
	return manifest.CronKey
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
