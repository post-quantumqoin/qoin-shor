package actors

import (
	"bytes"

	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/post-quantumqoin/core-types/exitcode"

	"github.com/post-quantumqoin/qoin-shor/core/contracts/aerrors"
)

func SerializeParams(i cbg.CBORMarshaler) ([]byte, aerrors.ActorError) {
	buf := new(bytes.Buffer)
	if err := i.MarshalCBOR(buf); err != nil {
		// TODO: shouldn't this be a fatal error?
		return nil, aerrors.Absorb(err, exitcode.ErrSerialization, "failed to encode parameter")
	}
	// fmt.Println("SerializeParams MarshalCBOR")
	return buf.Bytes(), nil
}
