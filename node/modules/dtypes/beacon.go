package dtypes

import "github.com/post-quantumqoin/core-types/abi"

type DrandSchedule []DrandPoint

type DrandPoint struct {
	Start  abi.ChainEpoch
	Config DrandConfig
}

type DrandConfig struct {
	Servers       []string
	Relays        []string
	ChainInfoJSON string
	IsChained     bool // Prior to Drand quicknet, beacons form a chain, post quicknet they do not (FIP-0063)
}
