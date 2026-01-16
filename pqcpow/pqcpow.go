package pqcpow

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/post-quantumqoin/core-types/abi"
	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/chain/types"
	"github.com/post-quantumqoin/qoin-shor/pqccrypto/mqphash"
	"github.com/post-quantumqoin/qoin-shor/pqccrypto/shake3"
)

var log = logging.Logger("pqcpow")

var (
	ErrXNotFound     = errors.New("x not found")
	ErrXFoundOutTime = errors.New("found x run out time")
	NewBlockheads    = errors.New("new blockheads have been created to follow the new blockheads")
)

const (
	EquationsOffset  = 31
	VariablesN       = 5
	WhichXWidth      = 1000
	NbitSampleRate   = 1
	ReferenceSeconds = 8 * 1000 //8s
	MaxThreshold     = 480      //480
	StepSize         = 2
)

//	const powParameter = {
//		equationsOffset: 31,
//		nbitSampleRate: 200,
//		referenceSeconds: 600
//	}
type PqcPowAPI interface {
	ChainGetTipSetByHeight(context.Context, abi.ChainEpoch, types.TipSetKey) (*types.TipSet, error)
	ChainHead(context.Context) (*types.TipSet, error) //perm:read
	ChainNotify(context.Context) (<-chan []*api.HeadChange, error)
	ChainGetTipSet(context.Context, types.TipSetKey) (*types.TipSet, error)
}

/**
 * @param {Buffer} seed - Block Header, Does not contain
 * @param {Buffer} nbit - 2 byte; First is equations (Rough adjustment), last is threshold (Fine-tune)
 * @param {Buffer} x - 32byte; Hash input
 */
func VerifyPoW(seed []byte, nbit []byte, x []byte) bool {
	equationsN := int(nbit[0]) + EquationsOffset
	variablesN := equationsN + VariablesN
	var threshold int

	//before UpgradeYellowStoneHeight height nbit is defined:|equations[0]|threshold[1]|
	//After UpgradeYellowStoneHeight height nbit is defined:|equations[0]|threshold[1:2]|seed[2:34]|
	if len(nbit) > 2 { //TODO:upgrade v2.0
		threshold = int(nbit[2])<<8 | int(nbit[1])
	} else {
		threshold = int(nbit[1])
	}

	fmt.Println("verifyPoW threshold:", threshold)
	// fmt.Println("verifyPoW seed: nbit: x:", seed, nbit, x)
	// fmt.Println("verifyPoW equationsN: variablesN: threshold:", equationsN, variablesN, threshold)
	if len(x) != 32 {
		return false
	}

	//------- step 1 create mpqhash function -------
	powHash := mqphash.CreateMQP(seed, equationsN, variablesN)

	for i := powHash.VariablesByte; i < 32; i++ {
		if x[i] != 0 {
			return false
		}
	}

	if (x[powHash.VariablesByte-1] & (^(0xff << powHash.UnwantedVariablesBit))) != 0 {
		return false
	}

	x = x[0:powHash.VariablesByte]
	//------- step 2 mpqhash(x) = zero, Confirm whether x is the solution of the system of equations -------
	isSolution := powHash.CheckIsSolution(x)
	if !isSolution {
		return false
	}
	// fmt.Println("verifyPoW isSolution:", isSolution)
	//------- step 3 Let x2 = shake256(x), and set the extra bits to 0 -------
	x2 := shake3.Shake256XOF(x, powHash.VariablesByte)
	x2[len(x2)-1] >>= powHash.UnwantedVariablesBit
	x2[len(x2)-1] <<= powHash.UnwantedVariablesBit

	// fmt.Println("verifyPoW x2:", x2)
	//------- step 4 Get the hashvalue of mpqhash(x2) -------
	hashVal := powHash.Update(x2)
	if len(hashVal) == 0 {
		return false
	}
	// fmt.Println("verifyPoW hashVal:", hashVal)
	//------- step 5 Let hashVal = shake256(hashVal), and take the first 9 bits as the integer value -------
	hashVal2 := shake3.Shake256XOF(hashVal, 2)                // 2 byte = 16 bit
	Val2 := (int(hashVal2[0]) << 1) | (int(hashVal2[1]) >> 7) //first 9 bits to integer value (big-endian)
	fmt.Println("verifyPoW hashVal2:", Val2)
	//------- step 6 Compare with threshold -------
	if threshold < Val2 {
		return true
	} else {
		return false
	}
}

func CalculateNbit(targetTime uint64, lastNbit []byte, windowSize uint64, windowTimeStart uint64, windowTimeEnd uint64) []byte {
	// fmt.Println("CalculateNbit targetTime: lastNbit: windowSize: windowTimeStart: windowTimeEnd:", targetTime, hex.EncodeToString(lastNbit), windowSize, windowTimeStart, windowTimeEnd)
	exponent := uint8(lastNbit[0])
	threshold := uint8(lastNbit[1])
	actualTime := (windowTimeEnd - windowTimeStart) * 1000 / windowSize
	timeRatio := float64(actualTime) / float64(targetTime)
	speedRatio := 1 / timeRatio
	fmt.Println("CalculateNbit actualTime: timeRatio: speedRatio:", actualTime, timeRatio, speedRatio)
	nbitDigitize := float64(exponent) + math.Log2(512/(512-float64(threshold)))
	newNbitDigitize := nbitDigitize + math.Log2(float64(speedRatio))
	var newExponent, newThreshold float64
	// fmt.Println("CalculateNbit nbitDigitize: newNbitDigitize: ", nbitDigitize, newNbitDigitize)
	if newNbitDigitize < 1 {
		newNbitDigitize = 1
		newExponent = 1
		newThreshold = 0
	} else {
		newExponent = math.Floor(newNbitDigitize)
		newThreshold = 512 / math.Pow(2, (newNbitDigitize-newExponent))
		newThreshold = math.Floor(512 - newThreshold)
	}

	if newThreshold > 255 {
		newThreshold = 255
	} else if newThreshold < 0 {
		newThreshold = 0
	}
	fmt.Println("CalculateNbit newExponent: newThreshold: ", newExponent, newThreshold)
	newNbit := make([]byte, 2)
	newNbit[0] = uint8(newExponent)
	newNbit[1] = uint8(newThreshold)
	fmt.Println("CalculateNbit newNbit: ", newNbit)
	return newNbit
}

func FixCalculateNbit(targetTime uint64, lastNbit []byte, windowSize uint64, windowTimeStart uint64, windowTimeEnd uint64, nulls int) []byte {
	// fmt.Println("CalculateNbit targetTime: lastNbit: windowSize: windowTimeStart: windowTimeEnd:", targetTime, hex.EncodeToString(lastNbit), windowSize, windowTimeStart, windowTimeEnd)
	exponent := uint8(lastNbit[0])
	// threshold := uint8(lastNbit[1])
	var threshold int

	threshold = int(lastNbit[2])<<8 | int(lastNbit[1])

	actualTime := (windowTimeEnd - windowTimeStart) * 1000 / windowSize
	if actualTime == 0 {
		actualTime = 1000
	}

	timeRatio := float64(actualTime) / float64(targetTime)
	speedRatio := 1 / timeRatio
	fmt.Println("FixCalculateNbit actualTime: timeRatio: exponent:", actualTime, timeRatio, exponent)
	var newExponent uint8
	var newThreshold int

	if timeRatio > 1 {
		newThreshold = threshold - StepSize + StepSize*int(math.Floor(math.Log2(float64(speedRatio)))/2) - nulls*StepSize
		fmt.Println("FixCalculateNbit newThreshold: threshold: math.Floor(math.Log2(float64(speedRatio))):", newThreshold, threshold, math.Floor(math.Log2(float64(speedRatio))))
	} else {
		newThreshold = threshold + StepSize + StepSize*int(math.Floor(math.Log2(float64(speedRatio)))/2)
		fmt.Println("FixCalculateNbit newThreshold: threshold: math.Floor(math.Log2(float64(speedRatio))):", newThreshold, threshold, math.Floor(math.Log2(float64(speedRatio))))
	}
	newExponent = exponent
	if newThreshold > (MaxThreshold) {
		newExponent = newExponent + 1
		fmt.Println("FixCalculateNbit newExponent: exponent:", newExponent, exponent)
		newThreshold = newThreshold - MaxThreshold
	} else if newThreshold < 0 {
		newExponent = newExponent - 1
		fmt.Println("FixCalculateNbit newExponent: exponent:", newExponent, exponent)
		newThreshold = MaxThreshold + newThreshold
	}
	if newExponent == 0 {
		newExponent = 2
	}

	fmt.Println("FixCalculateNbit newExponent: newThreshold: ", newExponent, newThreshold)
	newNbit := make([]byte, 3)
	// newExponent = 2
	// newThreshold = 220
	newNbit[0] = uint8(newExponent)
	newNbit[1] = byte(uint16(newThreshold) & 0xFF)
	newNbit[2] = byte((uint16(newThreshold) >> 8) & 0xFF)
	// newNbit[0] = 6
	// newNbit[1] = 144
	fmt.Println("FixCalculateNbit newNbit: ", newNbit)
	return newNbit
}

func RefreshNbit(lastNbit []byte) []byte {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	newNbit := lastNbit[:3]
	newNbit = append(newNbit, bytes...)

	fmt.Println("RefreshNbit newNbit: ", newNbit)
	return newNbit
}

func GetNbit(ctx context.Context, p PqcPowAPI, minedBlock *types.BlockMsg) ([]byte, error) {
	var newNbit []byte //TODO:upgrade v2.0
	var lastNbit []byte

	lastTipSet, err := p.ChainGetTipSet(ctx, types.NewTipSetKey(minedBlock.Header.Parents...))
	if err != nil {
		return nil, err
	}
	lastNbit = lastTipSet.PqcPowProof().Nbit
	if lastTipSet.Height() != 0 && lastTipSet.Height()%NbitSampleRate == 0 {
		startTime := lastTipSet.MinTimestamp()
		endTime := lastTipSet.MinTimestamp() + lastTipSet.MinMiningTimestamp()

		// if build.UpgradeYellowStoneHeight >= 0 && lastTipSet.Height() > build.UpgradeYellowStoneHeight {
		nulls := minedBlock.Header.Height - (lastTipSet.Height() + 1) //TODO:upgrade v2.0
		if nulls > 0 {
			endTime = startTime + uint64(nulls)*15 //15s
		}

		log.Infow("GetNbit ", "lastNbit:", lastNbit, "lastTipSet.Height: ",
			lastTipSet.Height(), "startTime:", time.Unix(int64(startTime), 0),
			"endTime:", time.Unix(int64(endTime), 0))
		// newNbit = CalculateNbit(ReferenceSeconds, lastNbit, NbitSampleRate, startTime, endTime)
		newNbit = FixCalculateNbit(ReferenceSeconds, lastNbit, NbitSampleRate, startTime, endTime, int(nulls))
		// } else if build.UpgradeYellowStoneHeight >= 0 && lastTipSet.Height() == build.UpgradeYellowStoneHeight {
		// 	newNbit = append(newNbit, 0x04)
		// 	newNbit = append(newNbit, 0x00)
		// 	newNbit = append(newNbit, 0x00)
		// } else {
		// 	log.Infow("GetNbit ", "lastNbit:", lastNbit, "lastTipSet.Height: ",
		// 		lastTipSet.Height(), "startTime:", time.Unix(int64(startTime), 0),
		// 		"endTime:", time.Unix(int64(endTime), 0))
		// 	var referenceSeconds uint64 = 15000 //15s
		// 	newNbit = CalculateNbit(referenceSeconds, lastNbit, NbitSampleRate, startTime, endTime)
		// }

	} else {
		fmt.Println("GetNbit ", "lastNbit:", lastNbit, "lastTipSet.Height: ", lastTipSet.Height())
		newNbit = lastNbit
	}
	return newNbit, nil
}

func PqcPowProof(ctx context.Context, seed []byte, nbit []byte, p PqcPowAPI, tm *time.Ticker) ([]byte, error) {
	m := int(nbit[0]) + EquationsOffset
	n := m + VariablesN
	mh := mqphash.CreateMQP(seed, m, n)
	fmt.Println("PqcPowProof seed len: nbit: m: n: len(mh.Seed):", len(seed), nbit, m, n, len(mh.Seed))
	whichXWidth := WhichXWidth
	c, err := NewController(mh, nbit, whichXWidth)
	if err != nil {
		return nil, err
	}
	ts, err := p.ChainHead(ctx)
	if err != nil {
		return nil, err
	}
	notifs, err := p.ChainNotify(ctx)
	if err != nil {
		return nil, err
	}
	x, err := c.Run(notifs, ts.Height(), tm)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func getDifficultyByNbit(nbit []byte) float64 {
	exponent := float64(nbit[0])
	// threshold := float64(nbit[1])
	threshold := int(nbit[2])<<8 | int(nbit[1])
	return exponent + math.Log2(512/(512-float64(threshold)))
}
