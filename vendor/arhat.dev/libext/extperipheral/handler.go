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

package extperipheral

import (
	"context"
	"fmt"
	"sync"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/log"
	"arhat.dev/pkg/wellknownerrors"

	"arhat.dev/libext/types"
	"arhat.dev/libext/util"
)

type cmdHandleFunc func(ctx context.Context, p Peripheral, payload []byte) (interface{}, error)

func NewHandler(logger log.Interface, unmarshal types.UnmarshalFunc, impl PeripheralConnector) types.Handler {
	mu := new(sync.RWMutex)
	h := &Handler{
		BaseHandler: util.NewBaseHandler(mu),

		logger: logger,

		unmarshal:   unmarshal,
		impl:        impl,
		peripherals: new(sync.Map),

		funcMap: nil,
	}

	h.funcMap = map[arhatgopb.CmdType]cmdHandleFunc{
		arhatgopb.CMD_PERIPHERAL_OPERATE:         h.handlePeripheralOperate,
		arhatgopb.CMD_PERIPHERAL_COLLECT_METRICS: h.handlePeripheralMetricsCollect,
	}

	return h
}

type Handler struct {
	*util.BaseHandler

	logger log.Interface

	unmarshal   types.UnmarshalFunc
	impl        PeripheralConnector
	peripherals *sync.Map

	funcMap map[arhatgopb.CmdType]cmdHandleFunc
}

func (c *Handler) HandleCmd(
	ctx context.Context,
	id, seq uint64,
	kind arhatgopb.CmdType,
	payload []byte,
) (interface{}, error) {
	switch kind {
	case arhatgopb.CMD_PERIPHERAL_CLOSE:
		c.logger.D("removing peripheral")
		c.removePeripheral(ctx, id)
		return &arhatgopb.DoneMsg{}, nil
	case arhatgopb.CMD_PERIPHERAL_CONNECT:
		c.logger.D("connecting peripheral")
		err := c.handlePeripheralConnect(ctx, id, payload)
		if err != nil {
			return nil, err
		}
		return &arhatgopb.DoneMsg{}, nil
	default:
	}

	c.logger.D("working on peripheral specific operation")
	// requires peripheral
	handle, ok := c.funcMap[kind]
	if !ok {
		c.logger.I("unknown peripheral cmd type", log.Int32("kind", int32(kind)))
		return nil, fmt.Errorf("unknown cmd")
	}

	p, ok := c.getPeripheral(id)
	if !ok {
		return nil, wellknownerrors.ErrNotFound
	}

	ret, err := handle(ctx, p, payload)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (c *Handler) handlePeripheralConnect(ctx context.Context, peripheralID uint64, payload []byte) (err error) {
	if _, loaded := c.peripherals.Load(peripheralID); loaded {
		return fmt.Errorf("invalid duplicate peripheral id")
	}

	spec := new(arhatgopb.PeripheralConnectCmd)
	err = c.unmarshal(payload, spec)
	if err != nil {
		return fmt.Errorf("failed to unmarshal PeripheralConnectCmd: %w", err)
	}

	p, err := c.impl.Connect(ctx, spec.Target, spec.Params, spec.Tls)
	if err != nil {
		return fmt.Errorf("failed to establish connection to peripheral: %w", err)
	}

	defer func() {
		if err != nil {
			p.Close(ctx)
		}
	}()

	if _, loaded := c.peripherals.LoadOrStore(peripheralID, p); loaded {
		return fmt.Errorf("invalid duplicate peripheral")
	}

	return nil
}

func (c *Handler) getPeripheral(peripheralID uint64) (Peripheral, bool) {
	i, ok := c.peripherals.Load(peripheralID)
	if !ok {
		return nil, false
	}

	p, ok := i.(Peripheral)
	if !ok {
		c.peripherals.Delete(peripheralID)
		return nil, false
	}

	return p, true
}

func (c *Handler) removePeripheral(ctx context.Context, peripheralID uint64) {
	p, ok := c.getPeripheral(peripheralID)
	if ok {
		p.Close(ctx)
	}

	c.peripherals.Delete(peripheralID)
}

func (c *Handler) handlePeripheralOperate(
	ctx context.Context, p Peripheral, payload []byte,
) (interface{}, error) {
	spec := new(arhatgopb.PeripheralOperateCmd)
	err := c.unmarshal(payload, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal PeripheralOperateCmd: %w", err)
	}

	ret, err := p.Operate(ctx, spec.Params, spec.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute operation: %w", err)
	}

	return &arhatgopb.PeripheralOperationResultMsg{Result: ret}, nil
}

func (c *Handler) handlePeripheralMetricsCollect(
	ctx context.Context, p Peripheral, payload []byte,
) (interface{}, error) {
	spec := new(arhatgopb.PeripheralMetricsCollectCmd)
	err := c.unmarshal(payload, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal PeripheralMetricsCollectCmd: %w", err)
	}

	ret, err := p.CollectMetrics(ctx, spec.Params)
	if err != nil {
		return nil, fmt.Errorf("failed to collect peripheral metrics: %w", err)
	}

	return &arhatgopb.PeripheralMetricsMsg{Values: ret}, nil
}
