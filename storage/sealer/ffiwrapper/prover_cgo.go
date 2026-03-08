//go:build cgo
// +build cgo

package ffiwrapper

import (
	"github.com/post-quantumqoin/core-types/proof"
	ffi "github.com/post-quantumqoin/qvm"

	"github.com/post-quantumqoin/qoin-shor/storage/sealer/storiface"
)

var ProofProver = proofProver{}

var _ storiface.Prover = ProofProver

type proofProver struct{}

func (v proofProver) AggregateSealProofs(aggregateInfo proof.AggregateSealVerifyProofAndInfos, proofs [][]byte) ([]byte, error) {
	return ffi.AggregateSealProofs(aggregateInfo, proofs)
}
