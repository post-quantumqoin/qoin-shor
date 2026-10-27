//go:build debug || 2k || testground || calibnet || butterflynet || interopnet || quantumresistnet
// +build debug 2k testground calibnet butterflynet interopnet quantumresistnet

package build

import "github.com/post-quantumqoin/core-types/abi"

// Non-mainnet networks require address binding from genesis.
const UpgradePqcAddressBindingHeight = abi.ChainEpoch(-1)
