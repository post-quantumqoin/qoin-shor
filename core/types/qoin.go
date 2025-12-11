package types

import (
	"encoding"
	"fmt"
	"math/big"
	"strings"

	"github.com/post-quantumqoin/qoin-shor/build"
)

type QOIN BigInt

func (q QOIN) String() string {
	if q.Int == nil {
		return "0 Q"
	}
	return q.Unitless() + " Q"
}

func (q QOIN) Unitless() string {
	r := new(big.Rat).SetFrac(q.Int, big.NewInt(int64(build.FilecoinPrecision)))
	if r.Sign() == 0 {
		return "0"
	}
	return strings.TrimRight(strings.TrimRight(r.FloatString(18), "0"), ".")
}

var AttoQ = NewInt(1)
var FemtoQ = BigMul(AttoQ, NewInt(1000))
var PicoQ = BigMul(FemtoQ, NewInt(1000))
var NanoQ = BigMul(PicoQ, NewInt(1000))

var qunitPrefixes = []string{"a", "f", "p", "n", "μ", "m"}

func (q QOIN) Short() string {
	n := BigInt(q).Abs()

	dn := uint64(1)
	var prefix string
	for _, p := range qunitPrefixes {
		if n.LessThan(NewInt(dn * 1000)) {
			prefix = p
			break
		}
		dn *= 1000
	}

	r := new(big.Rat).SetFrac(q.Int, big.NewInt(int64(dn)))
	if r.Sign() == 0 {
		return "0"
	}

	return strings.TrimRight(strings.TrimRight(r.FloatString(3), "0"), ".") + " " + prefix + "Q"
}

func (q QOIN) Nano() string {
	r := new(big.Rat).SetFrac(q.Int, big.NewInt(int64(1e9)))
	if r.Sign() == 0 {
		return "0"
	}

	return strings.TrimRight(strings.TrimRight(r.FloatString(9), "0"), ".") + " nQ"
}

func (q QOIN) Format(s fmt.State, ch rune) {
	switch ch {
	case 's', 'v':
		fmt.Fprint(s, q.String())
	default:
		q.Int.Format(s, ch)
	}
}

func (q QOIN) MarshalText() (text []byte, err error) {
	return []byte(q.String()), nil
}

func (q QOIN) UnmarshalText(text []byte) error {
	p, err := ParseQ(string(text))
	if err != nil {
		return err
	}

	q.Int.Set(p.Int)
	return nil
}

func ParseQ(s string) (QOIN, error) {
	suffix := strings.TrimLeft(s, "-.1234567890")
	s = s[:len(s)-len(suffix)]
	var attoq bool
	if suffix != "" {
		norm := strings.ToLower(strings.TrimSpace(suffix))
		switch norm {
		case "", "q":
		case "attoq", "aq":
			attoq = true
		default:
			return QOIN{}, fmt.Errorf("unrecognized suffix: %q", suffix)
		}
	}

	if len(s) > 50 {
		return QOIN{}, fmt.Errorf("string length too large: %d", len(s))
	}

	r, ok := new(big.Rat).SetString(s) //nolint:gosec
	if !ok {
		return QOIN{}, fmt.Errorf("failed to parse %q as a decimal number", s)
	}

	if !attoq {
		r = r.Mul(r, big.NewRat(int64(build.FilecoinPrecision), 1))
	}

	if !r.IsInt() {
		var pref string
		if attoq {
			pref = "atto"
		}
		return QOIN{}, fmt.Errorf("invalid %sQ value: %q", pref, s)
	}

	return QOIN{r.Num()}, nil
}

func MustParseQ(s string) QOIN {
	n, err := ParseQ(s)
	if err != nil {
		panic(err)
	}

	return n
}

var _ encoding.TextMarshaler = (*QOIN)(nil)
var _ encoding.TextUnmarshaler = (*QOIN)(nil)
