package pqcminer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	graphsync "github.com/ipfs/go-graphsync/impl"
	gsnet "github.com/ipfs/go-graphsync/network"
	"github.com/ipfs/go-graphsync/storeutil"
	provider "github.com/ipni/index-provider"
	"github.com/libp2p/go-libp2p/core/host"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"golang.org/x/xerrors"

	dtimpl "github.com/filecoin-project/go-data-transfer/v2/impl"
	dtnet "github.com/filecoin-project/go-data-transfer/v2/network"
	dtgstransport "github.com/filecoin-project/go-data-transfer/v2/transport/graphsync"
	piecefilestore "github.com/post-quantumqoin/go-qoin-markets/filestore"
	piecestoreimpl "github.com/post-quantumqoin/go-qoin-markets/piecestore/impl"
	"github.com/post-quantumqoin/go-qoin-markets/retrievalmarket"
	retrievalimpl "github.com/post-quantumqoin/go-qoin-markets/retrievalmarket/impl"
	rmnet "github.com/post-quantumqoin/go-qoin-markets/retrievalmarket/network"
	"github.com/post-quantumqoin/go-qoin-markets/shared"
	"github.com/post-quantumqoin/go-qoin-markets/storagemarket"
	storageimpl "github.com/post-quantumqoin/go-qoin-markets/storagemarket/impl"
	"github.com/post-quantumqoin/go-qoin-markets/storagemarket/impl/storedask"
	smnet "github.com/post-quantumqoin/go-qoin-markets/storagemarket/network"
	"github.com/filecoin-project/go-paramfetch"
	"github.com/filecoin-project/go-statestore"
	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/abi"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/go-jsonrpc/auth"

	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/api/v0api"
	"github.com/post-quantumqoin/qoin-shor/api/v1api"
	"github.com/post-quantumqoin/qoin-shor/build"
	"github.com/post-quantumqoin/qoin-shor/core/contracts/builtin/miner"
	"github.com/post-quantumqoin/qoin-shor/core/events"
	"github.com/post-quantumqoin/qoin-shor/core/gen"
	"github.com/post-quantumqoin/qoin-shor/core/gen/slashfilter"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	blockstore "github.com/post-quantumqoin/qoin-shor/dbstore"
	"github.com/post-quantumqoin/qoin-shor/journal"
	"github.com/post-quantumqoin/qoin-shor/markets"
	"github.com/post-quantumqoin/qoin-shor/markets/dagstore"
	"github.com/post-quantumqoin/qoin-shor/markets/idxprov"
	marketevents "github.com/post-quantumqoin/qoin-shor/markets/loggers"
	"github.com/post-quantumqoin/qoin-shor/markets/pricing"
	lotusminer "github.com/post-quantumqoin/qoin-shor/miner"
	"github.com/post-quantumqoin/qoin-shor/node/config"
	"github.com/post-quantumqoin/qoin-shor/node/modules/dtypes"
	"github.com/post-quantumqoin/qoin-shor/node/modules/helpers"
	"github.com/post-quantumqoin/qoin-shor/node/repo"
	"github.com/post-quantumqoin/qoin-shor/storage/ctladdr"
	"github.com/post-quantumqoin/qoin-shor/storage/paths"
	sealing "github.com/post-quantumqoin/qoin-shor/storage/pipeline"
	"github.com/post-quantumqoin/qoin-shor/storage/pipeline/sealiface"
	"github.com/post-quantumqoin/qoin-shor/storage/sealer"
	"github.com/post-quantumqoin/qoin-shor/storage/sealer/storiface"
	"github.com/post-quantumqoin/qoin-shor/storage/wdpost"
)

var (
	StagingAreaDirName = "deal-staging"
)




func SetupBlockPqcProducer(lc fx.Lifecycle, mctx helpers.MetricsCtx, ds dtypes.MetadataDS, api v1api.FullNode, sf *slashfilter.SlashFilter, j journal.Journal) (*lotusminer.PqcMiner, error) {
	ctx := helpers.LifecycleCtx(mctx, lc)
	// minerAddr, err := api.WalletDefaultAddress(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// var maddr address.Address
	addrbye, err := ds.Get(ctx, datastore.NewKey("miner-address"))
	if err != nil {
		return nil, err
	}

	minerAddr, err := address.NewFromBytes(addrbye)
	if err != nil {
		return nil, err
	}
	fmt.Printf("SetupBlockPqcProducer minerAddr:", minerAddr)
	m := lotusminer.NewPqcMiner(api, minerAddr, sf, j)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := m.Start(ctx); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return m.Stop(ctx)
		},
	})

	return m, nil
}

