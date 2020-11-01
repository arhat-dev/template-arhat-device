package types

import (
	"arhat.dev/arhat-proto/arhatgopb"
	"github.com/gogo/protobuf/proto"
)

type Handler interface {
	HandleCmd(id uint64, kind arhatgopb.CmdType, payload []byte) (proto.Marshaler, error)
}
