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

package controller

import (
	"context"
	"fmt"
	"sync"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/log"
	"github.com/gogo/protobuf/proto"

	"arhat.dev/template-arhat-device/pkg/device"
)

type CmdHandleFunc func(dev *device.Device, payload []byte) (proto.Marshaler, error)

func NewController(ctx context.Context) (*Controller, error) {
	reg := &arhatgopb.RegisterMsg{
		Name: "",
	}
	regMsg, err := arhatgopb.NewDeviceMsg(0, 0, reg)
	if err != nil {
		return nil, fmt.Errorf("failed to create register message: %w", err)
	}

	return &Controller{
		ctx:    ctx,
		logger: log.Log.WithName("controller"),

		regMsg:  regMsg,
		devices: new(sync.Map),

		chRefreshed: make(chan *channelBundle, 1),
		mu:          new(sync.RWMutex),
	}, nil
}

type Controller struct {
	ctx    context.Context
	logger log.Interface

	regMsg *arhatgopb.Msg

	devices     *sync.Map
	currentCB   *channelBundle
	chRefreshed chan *channelBundle

	mu *sync.RWMutex
}

func (c *Controller) Start() error {
	go c.handleSession()

	return nil
}

func (c *Controller) handleSession() {
	handlerMap := map[arhatgopb.CmdType]CmdHandleFunc{
		arhatgopb.CMD_DEV_OPERATE:         c.handleDeviceOperate,
		arhatgopb.CMD_DEV_COLLECT_METRICS: c.handleDeviceMetricsCollect,
	}

	for {
		var cb *channelBundle
		select {
		case <-c.ctx.Done():
			return
		case cb = <-c.chRefreshed:
		}

		// new session, register first

		sendMsg := func(msg *arhatgopb.Msg) (sent bool) {
			select {
			case <-cb.closed:
				return false
			case cb.msgCh <- msg:
				return true
			case <-c.ctx.Done():
				return false
			}
		}

		if !sendMsg(c.regMsg) {
			c.logger.I("failed to send register msg")
			continue
		}

	loop:
		for cmd := range cb.cmdCh {
			var ret proto.Marshaler
			switch cmd.Kind {
			case arhatgopb.CMD_DEV_CLOSE:
				c.removeDevice(cmd.DeviceId)
				ret = &arhatgopb.DoneMsg{}
			case arhatgopb.CMD_DEV_CONNECT:
				err := c.handleDeviceConnect(cmd.DeviceId, cmd.Payload)
				if err != nil {
					ret = &arhatgopb.ErrorMsg{Description: err.Error()}
				} else {
					ret = &arhatgopb.DoneMsg{}
				}
			default:
				// requires device
				handle, ok := handlerMap[cmd.Kind]
				if !ok {
					c.logger.I("unknown cmd type", log.Int32("kind", int32(cmd.Kind)))
					ret = &arhatgopb.ErrorMsg{Description: "unknown cmd"}
					break
				}

				dev, ok := c.getDevice(cmd.DeviceId)
				if !ok {
					ret = &arhatgopb.ErrorMsg{Description: "device not found"}
					break
				}

				var err error
				ret, err = handle(dev, cmd.Payload)
				if err != nil {
					ret = &arhatgopb.ErrorMsg{Description: err.Error()}
					break
				}
			}

			msg, err := arhatgopb.NewDeviceMsg(cmd.DeviceId, cmd.Seq, ret)
			if err != nil {
				c.logger.I("failed to marshal msg", log.Error(err))
				panic(err)
			}

			if !sendMsg(msg) {
				break loop
			}
		}
	}
}

func (c *Controller) RefreshChannels() (cmdCh chan<- *arhatgopb.Cmd, msgCh <-chan *arhatgopb.Msg) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cb := newChannelBundle()

	select {
	case <-c.ctx.Done():
		return nil, nil
	case c.chRefreshed <- cb:
		if c.currentCB != nil {
			c.currentCB.Close()
		}
	}

	c.currentCB = cb

	return cb.cmdCh, cb.msgCh
}
