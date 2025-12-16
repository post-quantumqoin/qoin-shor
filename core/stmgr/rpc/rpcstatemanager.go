package rpcstmgr

import (
	"context"

	cbor "github.com/ipfs/go-ipld-cbor"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"

	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/blockstore"
	"github.com/post-quantumqoin/qoin-shor/core/contracts/adt"
	"github.com/post-quantumqoin/qoin-shor/core/contracts/builtin/paych"
	"github.com/post-quantumqoin/qoin-shor/core/stmgr"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

type RPCStateManager struct {
	gapi   api.Gateway
	cstore *cbor.BasicIpldStore
}

func NewRPCStateManager(api api.Gateway) *RPCStateManager {
	cstore := cbor.NewCborStore(blockstore.NewAPIBlockstore(api))
	return &RPCStateManager{gapi: api, cstore: cstore}
}

func (s *RPCStateManager) GetPaychState(ctx context.Context, addr address.Address, ts *types.TipSet) (*types.Actor, paych.State, error) {
	act, err := s.gapi.StateGetActor(ctx, addr, ts.Key())
	if err != nil {
		return nil, nil, err
	}

	actState, err := paych.Load(adt.WrapStore(ctx, s.cstore), act)
	if err != nil {
		return nil, nil, err
	}
	return act, actState, nil

}

func (s *RPCStateManager) LoadActorTsk(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*types.Actor, error) {
	return s.gapi.StateGetActor(ctx, addr, tsk)
}

func (s *RPCStateManager) LookupID(ctx context.Context, addr address.Address, ts *types.TipSet) (address.Address, error) {
	return s.gapi.StateLookupID(ctx, addr, ts.Key())
}

func (s *RPCStateManager) ResolveToDeterministicAddress(ctx context.Context, addr address.Address, ts *types.TipSet) (address.Address, error) {
	return s.gapi.StateAccountKey(ctx, addr, ts.Key())
}

func (s *RPCStateManager) Call(ctx context.Context, msg *types.Message, ts *types.TipSet) (*api.InvocResult, error) {
	return nil, xerrors.Errorf("RPCStateManager does not implement StateManager.Call")
}

var _ stmgr.StateManagerAPI = (*RPCStateManager)(nil)
