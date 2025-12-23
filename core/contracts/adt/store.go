package adt

import (
	"context"

	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/post-quantumqoin/specs-contracts/contracts/util/adt"
)

type Store interface {
	Context() context.Context
	cbor.IpldStore
}

func WrapStore(ctx context.Context, store cbor.IpldStore) Store {
	return adt.WrapStore(ctx, store)
}
