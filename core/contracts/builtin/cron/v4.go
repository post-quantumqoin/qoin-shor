package cron

import (
	"fmt"

	"github.com/ipfs/go-cid"

	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/manifest"
	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	cron4 "github.com/post-quantumqoin/specs-contracts/contracts/builtin/cron"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
)

var _ State = (*state4)(nil)

func load4(store adt.Store, root cid.Cid) (State, error) {
	out := state4{store: store}
	err := store.Get(store.Context(), root, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func make4(store adt.Store) (State, error) {
	out := state4{store: store}
	out.State = *cron4.ConstructState(cron4.BuiltInEntries())
	return &out, nil
}

type state4 struct {
	cron4.State
	store adt.Store
}

func (s *state4) GetState() interface{} {
	return &s.State
}

func (s *state4) ActorKey() string {
	return manifest.CronKey
}

func (s *state4) ActorVersion() actorstypes.Version {
	return actorstypes.Version4
}

func (s *state4) Code() cid.Cid {
	code, ok := actors.GetActorCodeID(s.ActorVersion(), s.ActorKey())
	if !ok {
		panic(fmt.Errorf("didn't find actor %v code id for actor version %d", s.ActorKey(), s.ActorVersion()))
	}

	return code
}
