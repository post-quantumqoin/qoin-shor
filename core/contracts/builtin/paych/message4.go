package paych

import (
	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/abi"
	paychtypes "github.com/post-quantumqoin/core-types/builtin/v8/paych"
	builtin4 "github.com/post-quantumqoin/specs-contracts/contracts/builtin"
	init4 "github.com/post-quantumqoin/specs-contracts/contracts/builtin/init"
	paych4 "github.com/post-quantumqoin/specs-contracts/contracts/builtin/paych"

	init_ "github.com/post-quantumqoin/qoin-shor/core/contracts/builtin/init"

	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

type message4 struct{ from address.Address }

func (m message4) Create(to address.Address, initialAmount abi.TokenAmount) (*types.Message, error) {

	actorCodeID := builtin4.PaymentChannelActorCodeID

	params, aerr := actors.SerializeParams(&paych4.ConstructorParams{From: m.from, To: to})
	if aerr != nil {
		return nil, aerr
	}
	enc, aerr := actors.SerializeParams(&init4.ExecParams{
		CodeCID:           actorCodeID,
		ConstructorParams: params,
	})
	if aerr != nil {
		return nil, aerr
	}

	return &types.Message{
		To:     init_.Address,
		From:   m.from,
		Value:  initialAmount,
		Method: builtin4.MethodsInit.Exec,
		Params: enc,
	}, nil
}

func (m message4) Update(paych address.Address, sv *paychtypes.SignedVoucher, secret []byte) (*types.Message, error) {
	params, aerr := actors.SerializeParams(&paych4.UpdateChannelStateParams{

		Sv: toV0SignedVoucher(*sv),

		Secret: secret,
	})
	if aerr != nil {
		return nil, aerr
	}

	return &types.Message{
		To:     paych,
		From:   m.from,
		Value:  abi.NewTokenAmount(0),
		Method: builtin4.MethodsPaych.UpdateChannelState,
		Params: params,
	}, nil
}

func (m message4) Settle(paych address.Address) (*types.Message, error) {
	return &types.Message{
		To:     paych,
		From:   m.from,
		Value:  abi.NewTokenAmount(0),
		Method: builtin4.MethodsPaych.Settle,
	}, nil
}

func (m message4) Collect(paych address.Address) (*types.Message, error) {
	return &types.Message{
		To:     paych,
		From:   m.from,
		Value:  abi.NewTokenAmount(0),
		Method: builtin4.MethodsPaych.Collect,
	}, nil
}
