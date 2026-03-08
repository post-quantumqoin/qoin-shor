package ffiwrapper

import (
	"golang.org/x/xerrors"

	rlepluslazy "github.com/post-quantumqoin/bitset/rle"
	"github.com/post-quantumqoin/core-types/abi"

	"github.com/post-quantumqoin/qoin-shor/storage/sealer/partialfile"
	"github.com/post-quantumqoin/qoin-shor/storage/sealer/storiface"
)

// merge gaps between ranges which are close to each other
//
//	TODO: more benchmarking to come up with more optimal number
const mergeGaps = 32 << 20

// TODO const expandRuns = 16 << 20 // unseal more than requested for future requests

func computeUnsealRanges(unsealed rlepluslazy.RunIterator, offset storiface.UnpaddedByteIndex, size abi.UnpaddedPieceSize) (rlepluslazy.RunIterator, error) {
	todo := partialfile.PieceRun(offset.Padded(), size.Padded())
	todo, err := rlepluslazy.Subtract(todo, unsealed)
	if err != nil {
		return nil, xerrors.Errorf("compute todo-unsealed: %w", err)
	}

	return rlepluslazy.JoinClose(todo, mergeGaps)
}
