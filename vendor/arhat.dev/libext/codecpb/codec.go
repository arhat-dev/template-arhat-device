package codecpb

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"arhat.dev/arhat-proto/arhatgopb"
	"github.com/gogo/protobuf/proto"

	"arhat.dev/libext/types"
)

type Codec struct{}

func (c *Codec) Type() arhatgopb.CodecType {
	return arhatgopb.CODEC_PROTOBUF
}

func (c *Codec) NewEncoder(w io.Writer) types.Encoder {
	return &Encoder{w}
}

func (c *Codec) NewDecoder(r io.Reader) types.Decoder {
	return &Decoder{bufio.NewReader(r)}
}

type Encoder struct {
	w io.Writer
}

func (enc *Encoder) Encode(any interface{}) error {
	var (
		data []byte
		err  error
	)
	switch t := any.(type) {
	case proto.Marshaler:
		data, err = t.Marshal()
	case proto.Message:
		data, err = proto.Marshal(t)
	default:
		return fmt.Errorf("invalid not protobuf message")
	}
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	sizeBuf := make([]byte, 10)
	i := binary.PutUvarint(sizeBuf, uint64(len(data)))
	_, err = enc.w.Write(sizeBuf[:i])
	if err != nil {
		return fmt.Errorf("failed to write message size: %w", err)
	}

	_, err = enc.w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write message body: %w", err)
	}

	return nil
}

type Decoder struct {
	r *bufio.Reader
}

func (dec *Decoder) Decode(out interface{}) error {
	size, err := binary.ReadUvarint(dec.r)
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("failed to read size of the message: %w", err)
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(dec.r, buf)
	if err != nil {
		return fmt.Errorf("failed to read message body: %w", err)
	}

	o, ok := out.(proto.Message)
	if !ok {
		return fmt.Errorf("invalid not protobuf message")
	}

	return proto.Unmarshal(buf, o)
}
