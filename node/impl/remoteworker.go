package impl

import (
	"context"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/core-types/abi"
	"github.com/post-quantumqoin/go-jsonrpc"
	"github.com/post-quantumqoin/go-jsonrpc/auth"

	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/api/client"
	"github.com/post-quantumqoin/qoin-shor/storage/sealer"
)

type remoteWorker struct {
	api.Worker
	closer jsonrpc.ClientCloser
}

func (r *remoteWorker) NewSector(ctx context.Context, sector abi.SectorID) error {
	return xerrors.New("unsupported")
}

func connectRemoteWorker(ctx context.Context, fa api.Common, url string) (*remoteWorker, error) {
	token, err := fa.AuthNew(ctx, []auth.Permission{"admin"})
	if err != nil {
		return nil, xerrors.Errorf("creating auth token for remote connection: %w", err)
	}

	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+string(token))

	wapi, closer, err := client.NewWorkerRPCV0(context.TODO(), url, headers)
	if err != nil {
		return nil, xerrors.Errorf("creating jsonrpc client: %w", err)
	}

	wver, err := wapi.Version(ctx)
	if err != nil {
		closer()
		return nil, err
	}

	if !wver.EqMajorMinor(api.WorkerAPIVersion0) {
		return nil, xerrors.Errorf("unsupported worker api version: %s (expected %s)", wver, api.WorkerAPIVersion0)
	}

	return &remoteWorker{wapi, closer}, nil
}

func (r *remoteWorker) Close() error {
	r.closer()
	return nil
}

var _ sealer.Worker = &remoteWorker{}
