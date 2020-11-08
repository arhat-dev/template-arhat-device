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

package libext

import (
	"context"
	"fmt"
	"sync"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/log"

	"arhat.dev/libext/types"
	"arhat.dev/libext/util"
)

// NewController creates a hub for message send/receive
func NewController(
	ctx context.Context,
	logger log.Interface,
	marshal types.MarshalFunc,
	h types.Handler,
) (*Controller, error) {
	return &Controller{
		ctx:    ctx,
		logger: logger,

		marshal:     marshal,
		handler:     h,
		currentCB:   nil,
		chRefreshed: make(chan *channelBundle, 1),

		closed: false,
		mu:     new(sync.RWMutex),
	}, nil
}

type Controller struct {
	ctx    context.Context
	logger log.Interface

	marshal     types.MarshalFunc
	handler     types.Handler
	currentCB   *channelBundle
	chRefreshed chan *channelBundle

	closed bool
	mu     *sync.RWMutex
}

func (c *Controller) Start() error {
	go c.handleSession()

	return nil
}

func (c *Controller) handleSession() {
	for {
		var (
			cb   *channelBundle
			more bool
		)

		select {
		case <-c.ctx.Done():
			return
		case cb, more = <-c.chRefreshed:
			if !more {
				return
			}
		}

		// new session, register first

		sendMsg := func(msg *arhatgopb.Msg) error {
			c.logger.V("sending msg")
			select {
			case <-cb.closed:
				return nil
			case cb.msgCh <- msg:
				return nil
			case <-c.ctx.Done():
				return c.ctx.Err()
			}
		}

		c.logger.D("receiving commands")
		err := func() error {
			ctx, cancel := context.WithCancel(c.ctx)
			defer cancel()

			c.handler.SetMsgSendFunc(sendMsg)

			// cmdCh will be closed once RefreshChannels called
			for cmd := range cb.cmdCh {
				ret, err := c.handler.HandleCmd(ctx, cmd.Id, cmd.Seq, cmd.Kind, cmd.Payload)
				if err != nil {
					ret = &arhatgopb.ErrorMsg{Description: err.Error()}
				}

				// no return message means async operation, will send message later
				if ret == nil {
					continue
				}

				kind := util.GetMsgType(ret)
				if kind == 0 {
					return fmt.Errorf("unknown response msg")
				}

				msg, err := util.NewMsg(c.marshal, kind, cmd.Id, cmd.Seq, ret)
				if err != nil {
					return fmt.Errorf("failed to marshal response msg")
				}

				err = sendMsg(msg)
				if err != nil {
					return err
				}
			}
			return nil
		}()

		if err != nil {
			cb.Close()
		}
	}
}

// Close controller, will not handle incoming commands anymore
func (c *Controller) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		c.closed = true
		close(c.chRefreshed)
	}
}

// RefreshChannels creates a new cmd and msg channel pair for new connection
// usually this function is called in conjunction with Client.ProcessNewStream
func (c *Controller) RefreshChannels() (cmdCh chan<- *arhatgopb.Cmd, msgCh <-chan *arhatgopb.Msg) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cb := newChannelBundle()
	if c.closed {
		// return a closed msg channel since we have been closed before
		cb.Close()
		return cb.cmdCh, cb.msgCh
	}

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
