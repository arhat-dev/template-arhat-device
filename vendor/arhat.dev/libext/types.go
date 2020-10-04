package libext

import (
	"io"

	"arhat.dev/arhat-proto/arhatgopb"
	"github.com/gogo/protobuf/proto"
)

type Encoder interface {
	Encode(any interface{}) error
}

type Decoder interface {
	Decode(out interface{}) error
}

type Codec interface {
	ContentType() string
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

type Handler interface {
	HandleCmd(id uint64, kind arhatgopb.CmdType, payload []byte) (proto.Marshaler, error)
}
