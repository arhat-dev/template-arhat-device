/*
Copyright 2020 The arhat.dev Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package extruntime

import (
	"context"
	"fmt"
	"sync"

	"arhat.dev/aranya-proto/aranyagopb"
	"arhat.dev/aranya-proto/aranyagopb/runtimepb"
	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/log"

	"arhat.dev/libext/types"
	"arhat.dev/libext/util"
)

type cmdHandleFunc func(ctx context.Context, sid uint64, payload []byte) (*runtimepb.Packet, error)

func NewHandler(logger log.Interface, impl RuntimeEngine) types.Handler {
	mu := new(sync.RWMutex)

	h := &Handler{
		BaseHandler: util.NewBaseHandler(mu),

		logger: logger,
		impl:   impl,

		funcMap: nil,

		streams: util.NewStreamManager(),

		mu: mu,
	}

	h.funcMap = map[runtimepb.PacketType]cmdHandleFunc{
		runtimepb.CMD_GET_INFO: h.handleGetInfo,

		runtimepb.CMD_EXEC:         h.handleExec,
		runtimepb.CMD_ATTACH:       h.handleAttach,
		runtimepb.CMD_LOGS:         h.handleLogs,
		runtimepb.CMD_TTY_RESIZE:   h.handleTtyResize,
		runtimepb.CMD_PORT_FORWARD: h.handlePortForward,

		runtimepb.CMD_IMAGE_LIST:   h.handleImageList,
		runtimepb.CMD_IMAGE_ENSURE: h.handleImageEnsure,
		runtimepb.CMD_IMAGE_DELETE: h.handleImageDelete,

		runtimepb.CMD_POD_LIST:   h.handlePodList,
		runtimepb.CMD_POD_ENSURE: h.handlePodEnsure,
		runtimepb.CMD_POD_DELETE: h.handlePodDelete,
	}

	return h
}

type Handler struct {
	*util.BaseHandler

	logger log.Interface
	impl   RuntimeEngine

	funcMap map[runtimepb.PacketType]cmdHandleFunc

	streams *util.StreamManager

	mu *sync.RWMutex
}

var _zero = make([]byte, 0)

func (h *Handler) HandleCmd(
	ctx context.Context,
	sid, seq uint64,
	kind arhatgopb.CmdType,
	payload []byte,
) (interface{}, error) {
	switch kind {
	case arhatgopb.CMD_DATA_INPUT:
		if payload == nil {
			payload = _zero
		}
		h.streams.Write(sid, seq, payload)
		return nil, nil
	case arhatgopb.CMD_DATA_CLOSE:
		h.streams.Write(sid, seq, nil)
		return nil, nil
	case arhatgopb.CMD_RUNTIME_ARANYA_PROTO:
	default:
		return nil, fmt.Errorf("unknown cmd type %d", kind)
	}

	// is normal runtime cmd, decode as aranya-proto runtime packet

	pkt := new(runtimepb.Packet)
	err := pkt.Unmarshal(payload)
	if err != nil {
		return nil, err
	}

	do, ok := h.funcMap[pkt.Kind]
	if !ok {
		return nil, fmt.Errorf("unknown runtime cmd type %d", kind)
	}

	go func() {
		pkt, err = do(ctx, sid, pkt.Payload)
		if err != nil {
			var errData []byte
			errData, _ = (&aranyagopb.ErrorMsg{
				Kind:        aranyagopb.ERR_COMMON,
				Description: err.Error(),
				Code:        0,
			}).Marshal()

			err = h.SendMsg(&arhatgopb.Msg{
				Kind:    arhatgopb.MSG_RUNTIME_ARANYA_PROTO,
				Id:      sid,
				Ack:     seq,
				Payload: errData,
			})
			if err != nil {
				h.logger.I("failed to send error message")
			}

			return
		}

		if pkt == nil {
			// ignore empty packets
			return
		}

		data, err := pkt.Marshal()
		if err != nil {
			h.logger.I("failed to marshal result message", log.Error(err))
			return
		}

		err = h.SendMsg(&arhatgopb.Msg{
			Kind:    arhatgopb.MSG_RUNTIME_ARANYA_PROTO,
			Id:      sid,
			Ack:     seq,
			Payload: data,
		})
		if err != nil {
			h.logger.I("failed to send result message", log.Error(err))
		}
	}()

	return nil, nil
}

func (h *Handler) handleGetInfo(_ context.Context, _ uint64, _ []byte) (*runtimepb.Packet, error) {
	info := &runtimepb.RuntimeInfo{
		Name:          h.impl.Name(),
		Version:       h.impl.Version(),
		Os:            h.impl.OS(),
		OsImage:       h.impl.OSImage(),
		Arch:          h.impl.Arch(),
		KernelVersion: h.impl.KernelVersion(),
	}

	data, err := info.Marshal()
	if err != nil {
		return nil, err
	}

	return &runtimepb.Packet{
		Kind:    runtimepb.MSG_RUNTIME_INFO,
		Payload: data,
	}, nil
}
