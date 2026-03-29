package pqcpow

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	// "math"
	"strconv"
	"strings"
	"sync"
	"time"
	"golang.org/x/xerrors"

	logging "github.com/ipfs/go-log/v2"
	"github.com/post-quantumqoin/qoin-shor/pqcpow/pqc"
	"github.com/post-quantumqoin/qoin-shor/pqcpow/kernel"
	"github.com/post-quantumqoin/core-types/abi"
	"github.com/post-quantumqoin/qoin-shor/api"
	"github.com/post-quantumqoin/qoin-shor/build"
	"github.com/post-quantumqoin/qoin-shor/pqccrypto/mqphash"
)
var log = logging.Logger("pqcpow")


const maxN = 63 //If bigger then fix it.
type Controller struct {
	size      int
	fixNumber int
	fixIndex  int
	fixStr    []string

	numOfEquations int
	numOfVariables int
	devs           []*dev

	// devslk []*sync.Mutex
	// fixlk sync.Mutex
}

func NewController() (*Controller, error) {
	c := &Controller{}
	c.fixIndex = 0
	c.size = int(kernel.GetDeviceCount()) // get Device number.
	log.Infof("Device count: %d", c.size)
	devs, err := c.getDevs() //Registered Device List.
	if err != nil {
		return nil, err
	}
	c.devs = devs
	return c, nil
}

func (c *Controller)getDevs() ([]*dev, error) {
	var devs []*dev
	for devID := 0; devID < c.size; devID++ {
		d := NewDev(c)
		devs = append(devs, d)
	}
	return devs, nil
}

func (c *Controller)GeneratePQCProof(ctx context.Context, seed []byte, nbit []byte, p pqc.PqcPowAPI, tm *time.Ticker) ([]byte, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	m := int(nbit[0]) + pqc.EquationsOffset
	n := m + pqc.VariablesN
	mh := mqphash.CreateMQP(seed, m, n)
	log.Infof("GeneratePQCProof seed len: nbit: m: n: len(mh.Seed):", len(seed), nbit, m, n, len(mh.Seed))
	whichXWidth := pqc.WhichXWidth
	err := c.initPowParams(mh, nbit, whichXWidth)
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

func (c *Controller)initPowParams(mqphash *mqphash.MQPHash, nbit []byte, whichXWidth  int) (error){
	if(mqphash == nil || len(nbit) == 0 || whichXWidth <= 0){
		return xerrors.New("invalid parameters for InitPowParams")
	}
	c.numOfEquations = int(nbit[0]) + pqc.EquationsOffset
	c.numOfVariables = c.numOfEquations + pqc.VariablesN

	if c.size <= 1 { // set fixnumber.
		c.fixNumber = 0
	} else if c.size <= 2 {
		c.fixNumber = 1
	} else if c.size <= 4 {
		c.fixNumber = 2
	} else if c.size <= 8 {
		c.fixNumber = 3
	} else if c.size <= 16 {
		c.fixNumber = 4
	}

	//The number of variables exceeds the number of countable variables. Need fix it.
	if maxN < c.numOfVariables {
		diffN := c.numOfVariables - maxN
		if diffN > c.fixNumber {
			c.fixNumber = diffN
		}
	}
	// fmt.Println("c.fixNumber:", c.fixNumber)

	if c.fixNumber > 0 { //create fix str Array.
		// fLen := math.Pow(float64(2), float64(c.fixNumber))
		fLen := 1 << uint(c.fixNumber)
        c.fixStr = make([]string, 0, fLen)
		for i := 0; i < int(fLen); i++ {
			str := strconv.FormatInt(int64(i), 2)
			for j := 0; c.fixNumber > len(str); j++ {
				str = "0" + str
			}
			c.fixStr = append(c.fixStr, str)
		}
	}
	log.Infof("InitPowParams numOfEquations: numOfVariables: fixNumber: fixStr len:", c.numOfEquations, c.numOfVariables, c.fixNumber, len(c.fixStr))
	for _, dev := range c.devs {
        if err := dev.initParams(mqphash, nbit, whichXWidth); err != nil {
            return err
        }		
	}
	return nil

}

func (c *Controller) Run(notifs <-chan []*api.HeadChange,hgt abi.ChainEpoch,tm *time.Ticker) ([]byte, error) {
	//Receive blocks generated from devices
	result := make(chan []byte)
	var wg sync.WaitGroup
	defer wg.Wait() 
	defer close(result) 
	for devID := 0; devID < c.size; devID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c.devs[id].GetX(id, 0, result)
		}(devID)
	}

	tickerC := tm.C
	for {
		select {
		case r := <-result:
			// if len(r) == 0 {
			// 	log.Warnf("run x not found")
			// 	return nil, pqc.ErrXNotFound
			// }
			kernel.AbortCalc()
			return r, nil
		case <-tickerC:
			// if build.UpgradeYellowStoneHeight >= 0 && hgt > build.UpgradeYellowStoneHeight {
			log.Warnf("Run out time")
			// stopch <- true
			kernel.AbortCalc()
			// er = pqc.ErrXFoundOutTime
			// wg.Wait() 
			// stop()
			return nil, pqc.ErrXFoundOutTime
			// }
		case n := <-notifs:
			for _, change := range n {
				//a head change notifs,if a new block header is generated
				// the miner stops the Proof-of-Work for this block
				if hgt < change.Val.Height() {
					//In order to maintain fairness for miners, each miner has a mining retention time(about:4s)
					retention := time.Unix(int64(change.Val.MinTimestamp()-(uint64(15)-build.MinerRetentionTimeSecs)), 0)
					build.Clock.Sleep(build.Clock.Until(retention))
					log.Infow("new chain notify ", "now time:", build.Clock.Now(), "head MinTimestamp:", time.Unix(int64(change.Val.MinTimestamp()), 0))
					// stopch <- true
					kernel.AbortCalc()
					// wg.Wait() 
					// stop()
					return nil, pqc.NewBlockheads
					// er = pqc.NewBlockheads
				}

				log.Infow("new chain ", "hgt:", hgt, "Height:", change.Val.Height())
			}
		}
	}
}

func (c *Controller) GetNextFixStr() string {
	if c.fixNumber == 0 ||
		len(c.fixStr) == 0 ||
		c.fixIndex >= len(c.fixStr) {
		return ""
	}

	fmt.Println("GetNextFixStr: c.fixIndex:", c.fixIndex)
	rt := c.fixStr[c.fixIndex]
	c.fixIndex++
	return rt
}

type dev struct {
	m int
	n int
	//  startSMCount: number;
	whichXWidth  int
	startSMCount int
	mqphash      *mqphash.MQPHash
	nbit         []byte
	//  child: ChildProcessWithoutNullStreams;
	deviceID   int
	controller *Controller
	xbuf       []byte
	smCount    int

	lk sync.Mutex
}

func NewDev(ctr *Controller) *dev {
	d := &dev{
		// mqphash:     mqphash,
		// nbit:        nbit,
		// whichXWidth: whichXWidth,
		controller:  ctr,
	}
	// d.m = int(nbit[0]) + pqc.EquationsOffset
	// d.n = d.m + pqc.VariablesN
	// d.startSMCount = 0
	// d.lk = new(sync.Mutex)
	return d
}

func (d *dev) initParams(mqphash *mqphash.MQPHash, nbit []byte, whichXWidth int) error{
	if mqphash == nil {
        return xerrors.New("mqphash is nil")
    }
    if len(nbit) == 0 {
        return xerrors.New("nbit is empty")
    }
    if whichXWidth <= 0 {
        return xerrors.New("whichXWidth must be > 0")
    }
	// optional: if already initialized, return or no-op
    // if d.mqphash != nil {
    //     return xerrors.New("dev already initialized")
    // }
	// defensive copy
    d.nbit = append([]byte(nil), nbit...)
	d.mqphash = mqphash
	d.whichXWidth = whichXWidth
	d.m = int(nbit[0]) + pqc.EquationsOffset
	d.n = d.m + pqc.VariablesN
	d.startSMCount = 0

	return nil
}

func (d *dev) GetX(devID int, startSMCount int, results chan []byte) {
	d.lk.Lock()
	defer d.lk.Unlock()
	// defer close(results)
	d.deviceID = devID
	d.startSMCount = startSMCount
	fix := d.controller.GetNextFixStr()
	fmt.Println("GetX devID: fix: ", devID, fix)
	var verify bool = false
	for {
		if d.controller.fixNumber > 0 { //do fix
			var x []byte
			var err error
			if len(fix) != 0 {
				x, _, err = d.calculate(fix) // return d.xbuf = mf.fixBack(rx, fix)
				if err ==  pqc.ErrAbort {
					return
				}
				if err ==  pqc.ErrXNotFound {
					// results <- nil
					return
				}
				if err != nil {
					fmt.Println("GetX calculate err:", err)
					// results <- nil
					return
				}
			} else {
				results <- nil
				return
			}

			if len(x) == 0 {
				fmt.Println("x not found fix:", fix)
				fix = d.controller.GetNextFixStr()
				d.startSMCount = 0
				continue
			}

			if !d.mqphash.CheckIsSolution(x[0:d.mqphash.VariablesByte]) {
				fmt.Println(`Fix str '${fix}' check solution failed.`, fix)
				d.startSMCount = 0
				fix = d.controller.GetNextFixStr()
				continue
			}
		} else { //no fix
			_, x, err := d.calculate(fix)
			if err ==  pqc.ErrAbort {
				return
			}
			if err ==  pqc.ErrXNotFound {
				// results <- nil
				return
			}
			if err != nil {
				fmt.Println("GetX calculate err:", err)
				// results <- nil
				return
			}
			fmt.Println("GetX calculate:", x)

			if len(x) == 0 {
				fmt.Println("Check solution failed!")
				results <- nil
				return
			}

			if !d.checkSolution(x) { //d.xbuf = xBuf
				fmt.Println("Check solution failed!")
				results <- nil
				return
			}
		}
		fmt.Println("VerifyPoW seed: nbit: ", len(d.mqphash.Seed), d.nbit)
		//Proof of generation passes validation
		if pqc.VerifyPoW(d.mqphash.Seed, d.nbit, d.xbuf) {
			verify = true
		}

		if verify {
			fmt.Println("VerifyPoW is ok  d.deviceID: fix: d.xbuf:s", d.deviceID, fix, d.xbuf)
			results <- d.xbuf
			return
		}
		d.startSMCount = d.smCount + 1
		fmt.Println("GetX d.startSMCount: d.smCount:", d.startSMCount, d.smCount)
	}

}

func (d *dev) checkSolution(solution string) bool {
	fmt.Println("checkSolution solution:", solution, len(solution))
	//
	x := solution[len(solution)-d.n : len(solution)]

	// var sf []string
	for index := 0; index < d.mqphash.UnwantedVariablesBit; index++ {
		// sf = append(sf, "0")
		x += "0"
	}
	// x = strings.Join(sf, "") + s
	fmt.Println("checkSolution x:", x, len(x))
	xBuf := make([]byte, 32)
	index := 0

	for i := 0; i < len(x); i += 8 {
		// xBuf[index] = parseInt(x.slice(i, i+8), 2)
		end := i + 8
		r, _ := strconv.ParseInt(x[i:end], 2, 32)
		// fmt.Println("checkSolution   r:%x", x[i:end], r)
		xBuf[index] = byte(r)
		index++
	}
	d.xbuf = xBuf

	fmt.Println("checkSolution xBuf:", hex.EncodeToString(xBuf), len(xBuf))
	return d.mqphash.CheckIsSolution(xBuf[0:d.mqphash.VariablesByte])
}

func (d *dev) calculate(fix string) ([]byte, string, error) {
	// d.lk.Lock()
	// defer d.lk.Unlock()

	var equations []string
	var coefficientBit int

	type rxresult struct {
		X       string `json:"x"`
		GpuTime string `json:"gpuTime"`
		Rate    string `json:"rate"`
		SmCount string `json:"smCount"`
		SmUse   string `json:"smUse"`
	}

	var rs rxresult
	if len(fix) != 0 {
		mf := pqc.NewFix(d.mqphash, len(fix))
		for _, equation := range d.mqphash.Equations {
			eq, _, _, _ := mf.FixOneEquation(fix, hex.EncodeToString(equation), d.mqphash.UnwantedCoefficientBit)
			equations = append(equations, hex.EncodeToString(eq))
		}
		// fmt.Println("calculate CudaGetX fix:", fix)

		fmt.Println("calculate  d.deviceID, fix d.m, mf.NewN(), d.whichXWidth, uint64(d.startSMCount), mf.NewCoe()", d.deviceID, fix, d.m, mf.NewN(), d.whichXWidth, uint64(d.startSMCount), mf.NewCoe())
		// for i, eq := range equations {
		// 	fmt.Printf("calculate fix:%s Equations:%s len:%d  index:%d \n", fix, eq, len(eq), i)
		// }
		rx := kernel.CudaGetX(d.deviceID, d.m, mf.NewN(), d.whichXWidth, uint64(d.startSMCount), mf.NewCoe(), equations)
		// CudaGetX(deviceID int, m int, n int, whichXWidth int, startSMCount uint64, coefficientBit int, xIn []string)
		srx := strings.Split(rx, "x found:")

		if strings.Contains(rx, "abort") {
			return nil, "", pqc.ErrAbort
		}
		if strings.Contains(rx, "x not found") {
			return nil, "", pqc.ErrXNotFound
		}
		// fmt.Println("calculate CudaGetX fix: srx[1]:", fix, srx[1])
		if err := json.Unmarshal([]byte(srx[1]), &rs); err != nil {
			return nil, rx, err
		}

		num, err := strconv.Atoi(rs.SmCount)
		if err != nil {
			return nil, "", err
		}
		d.smCount = num
		fmt.Println("calculate  d.deviceID: fix:  rs.SmCount:", d.deviceID, fix, rs.SmCount)
		d.xbuf = mf.FixBack(rs.X, fix)
		// fmt.Println("calculate fix: fixBack:", fix, hex.EncodeToString(d.xbuf))
		return d.xbuf, "", nil
	} else {
		for _, equation := range d.mqphash.Equations {
			equations = append(equations, hex.EncodeToString(equation))
		}
		coefficientBit = d.mqphash.Coefficient
	}
	fmt.Println("calculate  d.deviceID, d.m, mf.newN, d.whichXWidth, uint64(d.startSMCount), coefficientBit: len(equations)", d.deviceID, d.m, d.n, d.whichXWidth, uint64(d.startSMCount), coefficientBit, len(equations))
	// for i, equation := range equations {
	// 	fmt.Println("calculate  i: equation:", i, equation)
	// }
	rx := kernel.CudaGetX(d.deviceID, d.m, d.n, d.whichXWidth, uint64(d.startSMCount), coefficientBit, equations)
	srx := strings.Split(rx, "x found:")
	for _, val := range srx {
		fmt.Println(val)
	}
	if len(srx) <= 1 {
		return nil, "", pqc.ErrXNotFound
	}
	if err := json.Unmarshal([]byte(srx[1]), &rs); err != nil {
		return nil, rx, err
	}
	num, err := strconv.Atoi(rs.SmCount)
	if err != nil {
		return nil, "", err
	}
	d.smCount = num
	fmt.Println("calculate  d.deviceID: fix:  rs.SmCount:", d.deviceID, fix, rs.SmCount)

	return nil, rs.X, nil
}
