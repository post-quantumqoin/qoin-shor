package build

import (
	"github.com/post-quantumqoin/core-types/abi"
)

func IsNearUpgrade(epoch, upgradeEpoch abi.ChainEpoch) bool {
	if upgradeEpoch < 0 {
		return false
	}
	return epoch > upgradeEpoch-Finality && epoch < upgradeEpoch+Finality
}
