package modules

import (
	"go.uber.org/fx"

	"github.com/post-quantumqoin/qoin-shor/core/beacon"
	"github.com/post-quantumqoin/qoin-shor/core/index"
	"github.com/post-quantumqoin/qoin-shor/core/stmgr"
	"github.com/post-quantumqoin/qoin-shor/core/store"
	"github.com/post-quantumqoin/qoin-shor/core/vm"
	"github.com/post-quantumqoin/qoin-shor/node/modules/dtypes"
)

func StateManager(lc fx.Lifecycle, cs *store.ChainStore, exec stmgr.Executor, sys vm.SyscallBuilder, us stmgr.UpgradeSchedule, b beacon.Schedule, metadataDs dtypes.MetadataDS, msgIndex index.MsgIndex) (*stmgr.StateManager, error) {
	sm, err := stmgr.NewStateManager(cs, exec, sys, us, b, metadataDs, msgIndex)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: sm.Start,
		OnStop:  sm.Stop,
	})
	return sm, nil
}
