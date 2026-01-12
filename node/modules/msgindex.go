package modules

import (
	"context"

	"go.uber.org/fx"

	"github.com/post-quantumqoin/qoin-shor/core/index"
	"github.com/post-quantumqoin/qoin-shor/core/store"
	"github.com/post-quantumqoin/qoin-shor/node/modules/helpers"
	"github.com/post-quantumqoin/qoin-shor/node/repo"
)

func MsgIndex(lc fx.Lifecycle, mctx helpers.MetricsCtx, cs *store.ChainStore, r repo.LockedRepo) (index.MsgIndex, error) {
	basePath, err := r.SqlitePath()
	if err != nil {
		return nil, err
	}

	msgIndex, err := index.NewMsgIndex(helpers.LifecycleCtx(mctx, lc), basePath, cs)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return msgIndex.Close()
		},
	})

	return msgIndex, nil
}

func DummyMsgIndex() index.MsgIndex {
	return index.DummyMsgIndex
}
