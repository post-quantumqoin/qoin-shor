package paych

import (
	"golang.org/x/xerrors"

	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/abi"
	builtin12 "github.com/post-quantumqoin/core-types/builtin"
	init12 "github.com/post-quantumqoin/core-types/builtin/v12/init"
	paych12 "github.com/post-quantumqoin/core-types/builtin/v12/paych"
	paychtypes "github.com/post-quantumqoin/core-types/builtin/v8/paych"
	actorstypes "github.com/post-quantumqoin/core-types/contracts"

	init_ "github.com/post-quantumqoin/qoin-shor/core/contracts/builtin/init"

	actors "github.com/post-quantumqoin/qoin-shor/core/contracts"
	"github.com/post-quantumqoin/qoin-shor/core/types"
)

type message12 struct{ from address.Address }

func (m message12) Create(to address.Address, initialAmount abi.TokenAmount) (*types.Message, error) {

	actorCodeID, ok := actors.GetActorCodeID(actorstypes.Version12, "paymentchannel")
	if !ok {
		return nil, xerrors.Errorf("error getting actor paymentchannel code id for actor version %d", 12)
	}

	params, aerr := actors.SerializeParams(&paych12.ConstructorParams{From: m.from, To: to})
	if aerr != nil {
		return nil, aerr
	}
	enc, aerr := actors.SerializeParams(&init12.ExecParams{
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
		Method: builtin12.MethodsInit.Exec,
		Params: enc,
	}, nil
}

func (m message12) Update(paych address.Address, sv *paychtypes.SignedVoucher, secret []byte) (*types.Message, error) {
	params, aerr := actors.SerializeParams(&paych12.UpdateChannelStateParams{

		Sv: toV12SignedVoucher(*sv),

		Secret: secret,
	})
	if aerr != nil {
		return nil, aerr
	}

	return &types.Message{
		To:     paych,
		From:   m.from,
		Value:  abi.NewTokenAmount(0),
		Method: builtin12.MethodsPaych.UpdateChannelState,
		Params: params,
	}, nil
}

func toV12SignedVoucher(sv paychtypes.SignedVoucher) paych12.SignedVoucher {
	merges := make([]paych12.Merge, len(sv.Merges))
	for i := range sv.Merges {
		merges[i] = paych12.Merge{
			Lane:  sv.Merges[i].Lane,
			Nonce: sv.Merges[i].Nonce,
		}
	}

	return paych12.SignedVoucher{
		ChannelAddr:     sv.ChannelAddr,
		TimeLockMin:     sv.TimeLockMin,
		TimeLockMax:     sv.TimeLockMax,
		SecretHash:      sv.SecretHash,
		Extra:           (*paych12.ModVerifyParams)(sv.Extra),
		Lane:            sv.Lane,
		Nonce:           sv.Nonce,
		Amount:          sv.Amount,
		MinSettleHeight: sv.MinSettleHeight,
		Merges:          merges,
		Signature:       sv.Signature,
	}
}

func (m message12) Settle(paych address.Address) (*types.Message, error) {
	return &types.Message{
		To:     paych,
		From:   m.from,
		Value:  abi.NewTokenAmount(0),
		Method: builtin12.MethodsPaych.Settle,
	}, nil
}

func (m message12) Collect(paych address.Address) (*types.Message, error) {
	return &types.Message{
		To:     paych,
		From:   m.from,
		Value:  abi.NewTokenAmount(0),
		Method: builtin12.MethodsPaych.Collect,
	}, nil
}
