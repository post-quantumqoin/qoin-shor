package modules

import (
	"go.uber.org/fx"

	"github.com/post-quantumqoin/qoin-shor/core/gen/slashfilter/slashsvc"
	"github.com/post-quantumqoin/qoin-shor/node/config"
	"github.com/post-quantumqoin/qoin-shor/node/impl/full"
	"github.com/post-quantumqoin/qoin-shor/node/modules/helpers"
)

type consensusReporterModules struct {
	fx.In

	full.WalletAPI
	full.ChainAPI
	full.MpoolAPI
	full.SyncAPI
}

func RunConsensusFaultReporter(config config.FaultReporterConfig) func(mctx helpers.MetricsCtx, lc fx.Lifecycle, mod consensusReporterModules) error {
	return func(mctx helpers.MetricsCtx, lc fx.Lifecycle, mod consensusReporterModules) error {
		ctx := helpers.LifecycleCtx(mctx, lc)

		return slashsvc.SlashConsensus(ctx, &mod, config.ConsensusFaultReporterDataDir, config.ConsensusFaultReporterAddress)
	}
}
