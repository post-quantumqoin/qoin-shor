//go:build debug || 2k
// +build debug 2k

package build

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ipfs/go-cid"

	"github.com/post-quantumqoin/core-types/abi"
	actorstypes "github.com/post-quantumqoin/core-types/contracts"
	"github.com/post-quantumqoin/core-types/network"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/policy"
)

const BootstrappersFile = ""
const GenesisFile = ""

var NetworkBundle = "devnet"
var BundleOverrides map[actorstypes.Version]string
var ActorDebugging = false

var GenesisNetworkVersion = network.Version21

var UpgradeBreezeHeight = abi.ChainEpoch(-1)

const BreezeGasTampingDuration = 0

var UpgradeSmokeHeight = abi.ChainEpoch(-1)
var UpgradeIgnitionHeight = abi.ChainEpoch(-2)
var UpgradeRefuelHeight = abi.ChainEpoch(-3)
var UpgradeTapeHeight = abi.ChainEpoch(-4)

var UpgradeAssemblyHeight = abi.ChainEpoch(-5)
var UpgradeLiftoffHeight = abi.ChainEpoch(-6)

var UpgradeKumquatHeight = abi.ChainEpoch(-7)
var UpgradeCalicoHeight = abi.ChainEpoch(-9)
var UpgradePersianHeight = abi.ChainEpoch(-10)
var UpgradeOrangeHeight = abi.ChainEpoch(-11)
var UpgradeClausHeight = abi.ChainEpoch(-12)

var UpgradeTrustHeight = abi.ChainEpoch(-13)

var UpgradeNorwegianHeight = abi.ChainEpoch(-14)

var UpgradeTurboHeight = abi.ChainEpoch(-15)

var UpgradeHyperdriveHeight = abi.ChainEpoch(-16)

var UpgradeChocolateHeight = abi.ChainEpoch(-17)

var UpgradeOhSnapHeight = abi.ChainEpoch(-18)

var UpgradeSkyrHeight = abi.ChainEpoch(-19)

var UpgradeSharkHeight = abi.ChainEpoch(-20)

var UpgradeHyggeHeight = abi.ChainEpoch(-21)

var UpgradeLightningHeight = abi.ChainEpoch(-22)

var UpgradeThunderHeight = abi.ChainEpoch(-23)

var UpgradeWatermelonHeight = abi.ChainEpoch(-24)

var UpgradeDragonHeight = abi.ChainEpoch(20)

var UpgradePhoenixHeight = UpgradeDragonHeight + 120

// This fix upgrade only ran on calibrationnet
const UpgradeWatermelonFixHeight = -100

// This fix upgrade only ran on calibrationnet
const UpgradeWatermelonFix2Height = -101

var DrandSchedule = map[abi.ChainEpoch]DrandEnum{
	0:                    DrandMainnet,
	UpgradePhoenixHeight: DrandQuicknet,
}

var SupportedProofTypes = []abi.RegisteredSealProof{
	abi.RegisteredSealProof_StackedDrg2KiBV1,
	abi.RegisteredSealProof_StackedDrg8MiBV1,
}
var ConsensusMinerMinPower = abi.NewStoragePower(2048)
var MinVerifiedDealSize = abi.NewStoragePower(256)
var PreCommitChallengeDelay = abi.ChainEpoch(10)

func init() {
	policy.SetSupportedProofTypes(SupportedProofTypes...)
	policy.SetConsensusMinerMinPower(ConsensusMinerMinPower)
	policy.SetMinVerifiedDealSize(MinVerifiedDealSize)
	policy.SetPreCommitChallengeDelay(PreCommitChallengeDelay)

	getGenesisNetworkVersion := func(ev string, def network.Version) network.Version {
		hs, found := os.LookupEnv(ev)
		if found {
			h, err := strconv.Atoi(hs)
			if err != nil {
				log.Panicf("failed to parse %s env var", ev)
			}

			return network.Version(h)
		}

		return def
	}
	fmt.Println("params_2k.go getGenesisNetworkVersion:%s\n", getGenesisNetworkVersion)
	GenesisNetworkVersion = getGenesisNetworkVersion("QOIN_GENESIS_NETWORK_VERSION", GenesisNetworkVersion)
	fmt.Println("params_2k.go GenesisNetworkVersion:%s\n", GenesisNetworkVersion)
	getUpgradeHeight := func(ev string, def abi.ChainEpoch) abi.ChainEpoch {
		hs, found := os.LookupEnv(ev)
		if found {
			h, err := strconv.Atoi(hs)
			if err != nil {
				log.Panicf("failed to parse %s env var", ev)
			}

			return abi.ChainEpoch(h)
		}

		return def
	}

	UpgradeBreezeHeight = getUpgradeHeight("QOIN_BREEZE_HEIGHT", UpgradeBreezeHeight)
	UpgradeSmokeHeight = getUpgradeHeight("QOIN_SMOKE_HEIGHT", UpgradeSmokeHeight)
	UpgradeIgnitionHeight = getUpgradeHeight("QOIN_IGNITION_HEIGHT", UpgradeIgnitionHeight)
	UpgradeRefuelHeight = getUpgradeHeight("QOIN_REFUEL_HEIGHT", UpgradeRefuelHeight)
	UpgradeTapeHeight = getUpgradeHeight("QOIN_TAPE_HEIGHT", UpgradeTapeHeight)
	UpgradeAssemblyHeight = getUpgradeHeight("QOIN_ACTORSV2_HEIGHT", UpgradeAssemblyHeight)
	UpgradeLiftoffHeight = getUpgradeHeight("QOIN_LIFTOFF_HEIGHT", UpgradeLiftoffHeight)
	UpgradeKumquatHeight = getUpgradeHeight("QOIN_KUMQUAT_HEIGHT", UpgradeKumquatHeight)
	UpgradeCalicoHeight = getUpgradeHeight("QOIN_CALICO_HEIGHT", UpgradeCalicoHeight)
	UpgradePersianHeight = getUpgradeHeight("QOIN_PERSIAN_HEIGHT", UpgradePersianHeight)
	UpgradeOrangeHeight = getUpgradeHeight("QOIN_ORANGE_HEIGHT", UpgradeOrangeHeight)
	UpgradeClausHeight = getUpgradeHeight("QOIN_CLAUS_HEIGHT", UpgradeClausHeight)
	UpgradeTrustHeight = getUpgradeHeight("QOIN_ACTORSV3_HEIGHT", UpgradeTrustHeight)
	UpgradeNorwegianHeight = getUpgradeHeight("QOIN_NORWEGIAN_HEIGHT", UpgradeNorwegianHeight)
	UpgradeTurboHeight = getUpgradeHeight("QOIN_ACTORSV4_HEIGHT", UpgradeTurboHeight)
	UpgradeHyperdriveHeight = getUpgradeHeight("QOIN_HYPERDRIVE_HEIGHT", UpgradeHyperdriveHeight)
	UpgradeChocolateHeight = getUpgradeHeight("QOIN_CHOCOLATE_HEIGHT", UpgradeChocolateHeight)
	UpgradeOhSnapHeight = getUpgradeHeight("QOIN_OHSNAP_HEIGHT", UpgradeOhSnapHeight)
	UpgradeSkyrHeight = getUpgradeHeight("QOIN_SKYR_HEIGHT", UpgradeSkyrHeight)
	UpgradeSharkHeight = getUpgradeHeight("QOIN_SHARK_HEIGHT", UpgradeSharkHeight)
	UpgradeHyggeHeight = getUpgradeHeight("QOIN_HYGGE_HEIGHT", UpgradeHyggeHeight)
	UpgradeLightningHeight = getUpgradeHeight("QOIN_LIGHTNING_HEIGHT", UpgradeLightningHeight)
	UpgradeThunderHeight = getUpgradeHeight("QOIN_THUNDER_HEIGHT", UpgradeThunderHeight)
	UpgradeWatermelonHeight = getUpgradeHeight("QOIN_WATERMELON_HEIGHT", UpgradeWatermelonHeight)
	UpgradeDragonHeight = getUpgradeHeight("QOIN_DRAGON_HEIGHT", UpgradeDragonHeight)

	UpgradePhoenixHeight = getUpgradeHeight("QOIN_PHOENIX_HEIGHT", UpgradePhoenixHeight)
	DrandSchedule = map[abi.ChainEpoch]DrandEnum{
		0:                    DrandMainnet,
		UpgradePhoenixHeight: DrandQuicknet,
	}

	BuildType |= Build2k

}

const BlockDelaySecs = uint64(4)

const PropagationDelaySecs = uint64(1)

var EquivocationDelaySecs = uint64(0)

// SlashablePowerDelay is the number of epochs after ElectionPeriodStart, after
// which the miner is slashed
//
// Epochs
const SlashablePowerDelay = 20

// Epochs
const InteractivePoRepConfidence = 6

const BootstrapPeerThreshold = 1

// ChainId defines the chain ID used in the Ethereum JSON-RPC endpoint.
// As per https://github.com/ethereum-lists/chains
const Eip155ChainId = 31415926

var WhitelistedBlock = cid.Undef
