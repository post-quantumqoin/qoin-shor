package remotewallet

import (
	"context"

	"go.uber.org/fx"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/api/client"
	cliutil "github.com/post-quantumqoin/qoin-shor/cli/util"
	"github.com/post-quantumqoin/qoin-shor/node/modules/helpers"
)

type RemoteWallet struct {
	api.Wallet
}

func SetupRemoteWallet(info string) func(mctx helpers.MetricsCtx, lc fx.Lifecycle) (*RemoteWallet, error) {
	return func(mctx helpers.MetricsCtx, lc fx.Lifecycle) (*RemoteWallet, error) {
		ai := cliutil.ParseApiInfo(info)

		url, err := ai.DialArgs("v0")
		if err != nil {
			return nil, err
		}

		wapi, closer, err := client.NewWalletRPCV0(mctx, url, ai.AuthHeader())
		if err != nil {
			return nil, xerrors.Errorf("creating jsonrpc client: %w", err)
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				closer()
				return nil
			},
		})

		return &RemoteWallet{wapi}, nil
	}
}

func (w *RemoteWallet) Get() api.Wallet {
	if w == nil {
		return nil
	}

	return w
}
