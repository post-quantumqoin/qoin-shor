package events

import (
	"context"

	"github.com/ipfs/go-cid"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/abi"

	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

type uncachedAPI interface {
	ChainNotify(context.Context) (<-chan []*api.HeadChange, error)
	ChainGetPath(ctx context.Context, from, to types.TipSetKey) ([]*api.HeadChange, error)
	StateSearchMsg(ctx context.Context, from types.TipSetKey, msg cid.Cid, limit abi.ChainEpoch, allowReplaced bool) (*api.MsgLookup, error)

	StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error) // optional / for CalledMsg
}

type cache struct {
	*tipSetCache
	*messageCache
	uncachedAPI
}

func newCache(api EventHelperAPI, gcConfidence abi.ChainEpoch) *cache {
	return &cache{
		newTSCache(api, gcConfidence),
		newMessageCache(api),
		api,
	}
}
