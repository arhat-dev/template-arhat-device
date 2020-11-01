package types

import (
	"io"

	"arhat.dev/arhat-proto/arhatgopb"
)

type Encoder interface {
	Encode(any interface{}) error
}

type Decoder interface {
	Decode(out interface{}) error
}

type Codec interface {
	Type() arhatgopb.CodecType
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}
