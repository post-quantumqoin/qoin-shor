package mock

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-commp-utils/zerocomm"
	commcid "github.com/filecoin-project/go-fil-commcid"
	"github.com/post-quantumqoin/address"
	"github.com/post-quantumqoin/core-types/abi"
	"github.com/post-quantumqoin/core-types/big"
	"github.com/post-quantumqoin/core-types/builtin/v9/market"

	"github.com/post-quantumqoin/qoin-shor/core/types"
	"github.com/post-quantumqoin/qoin-shor/core/wallet/key"
	"github.com/post-quantumqoin/qoin-shor/genesis"
)

func CommDR(in []byte) (out [32]byte) {
	for i, b := range in {
		out[i] = ^b
	}

	return out
}

func PreSeal(spt abi.RegisteredSealProof, maddr address.Address, sectors int) (*genesis.Miner, *types.KeyInfo, error) {
	k, err := key.GenerateKey(types.KTBLS)
	if err != nil {
		return nil, nil, err
	}

	pqckey, err := key.PqcGenerateKey(types.KTPqc)
	if err != nil {
		return nil, nil, err
	}

	pqcWrKey, err := key.PqcGenerateKey(types.KTPqc)
	if err != nil {
		return nil, nil, err
	}

	ssize, err := spt.SectorSize()
	if err != nil {
		return nil, nil, err
	}

	genm := &genesis.Miner{
		ID:            maddr,
		Owner:         pqcWrKey.Address,
		Worker:        k.Address,
		MarketBalance: big.NewInt(0),
		PowerBalance:  big.NewInt(0),
		SectorSize:    ssize,
		Sectors:       make([]*genesis.PreSeal, sectors),
	}

	for i := range genm.Sectors {
		label, err := market.NewLabelFromString(fmt.Sprintf("%d", i))
		if err != nil {
			return nil, nil, xerrors.Errorf("failed to create label: %w", err)
		}

		preseal := &genesis.PreSeal{}

		preseal.ProofType = spt
		preseal.CommD = zerocomm.ZeroPieceCommitment(abi.PaddedPieceSize(ssize).Unpadded())
		d, _ := commcid.CIDToPieceCommitmentV1(preseal.CommD)
		r := CommDR(d)
		preseal.CommR, _ = commcid.ReplicaCommitmentV1ToCID(r[:])
		preseal.SectorID = abi.SectorNumber(i + 1)
		preseal.Deal = market.DealProposal{
			PieceCID:             preseal.CommD,
			PieceSize:            abi.PaddedPieceSize(ssize),
			Client:               pqckey.Address,
			Provider:             maddr,
			Label:                label,
			StartEpoch:           1,
			EndEpoch:             10000,
			StoragePricePerEpoch: big.Zero(),
			ProviderCollateral:   big.Zero(),
			ClientCollateral:     big.Zero(),
		}
		preseal.DealClientKey = pqckey.KeyInfo

		genm.Sectors[i] = preseal
	}

	return genm, &pqcWrKey.KeyInfo, nil
}
