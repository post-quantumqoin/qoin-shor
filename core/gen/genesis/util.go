package genesis

import (
	"context"

	cbg "github.com/whyrusleeping/cbor-gen"
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/abi"

	"github.com/post-quantumqoin/qoin-shor/core/actors"
	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/vm"
)

func mustEnc(i cbg.CBORMarshaler) []byte {
	enc, err := actors.SerializeParams(i)
	if err != nil {
		panic(err) // ok
	}
	return enc
}

func doExecValue(ctx context.Context, vm vm.Interface, to, from address.Address, value types.BigInt, method abi.MethodNum, params []byte) ([]byte, error) {
	ret, err := vm.ApplyImplicitMessage(ctx, &types.Message{
		To:       to,
		From:     from,
		Method:   method,
		Params:   params,
		GasLimit: 1_000_000_000_000_000,
		Value:    value,
		Nonce:    0,
	})
	if err != nil {
		return nil, xerrors.Errorf("doExec apply message failed: %w", err)
	}

	if ret.ExitCode != 0 {
		return nil, xerrors.Errorf("failed to call method: %w", ret.ActorErr)
	}

	return ret.Return, nil
}
