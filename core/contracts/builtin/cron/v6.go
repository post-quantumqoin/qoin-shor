package cron

import (
	"fmt"

	"github.com/ipfs/go-cid"

	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	cron6 "github.com/post-quantumqoin/specs-contracts/contracts/builtin/cron"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state6)(nil)

func load6(store adt.Store, root cid.Cid) (State, error) {
	out := state6{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make6(store adt.Store) (State, error) {
	out := state6{store: store}
	out.State = *cron6.ConstructState(cron6.BuiltInEntries())
	return &out, nil
}

type state6 struct {
	cron6.State
	store adt.Store
}

func (s *state6) GetState() interface{} {
	return &s.State
}

func (s *state6) ActorKey() string {
	return manifest.CronKey
}

func (s *state6) ActorVersion() actorstypes.Version {
	return actorstypes.Version6
}

func (s *state6) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
