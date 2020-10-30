package extperipheral

import (
	"fmt"
	"sync"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/log"
	"github.com/gogo/protobuf/proto"
)

type cmdHandleFunc func(p Peripheral, payload []byte) (proto.Marshaler, error)

func NewHandler(logger log.Interface, impl PeripheralConnector) *Handler {
	return &Handler{
		logger:      logger,
		impl:        impl,
		peripherals: new(sync.Map),
	}
}

type Handler struct {
	logger log.Interface

	impl        PeripheralConnector
	peripherals *sync.Map
}

func (c *Handler) HandleCmd(id uint64, kind arhatgopb.CmdType, payload []byte) (proto.Marshaler, error) {
	handlerMap := map[arhatgopb.CmdType]cmdHandleFunc{
		arhatgopb.CMD_PERIPHERAL_OPERATE:         c.handlePeripheralOperate,
		arhatgopb.CMD_PERIPHERAL_COLLECT_METRICS: c.handlePeripheralMetricsCollect,
	}

	switch kind {
	case arhatgopb.CMD_PERIPHERAL_CLOSE:
		c.logger.D("removing peripheral")
		c.removePeripheral(id)
		return &arhatgopb.DoneMsg{}, nil
	case arhatgopb.CMD_PERIPHERAL_CONNECT:
		c.logger.D("connecting peripheral")
		err := c.handlePeripheralConnect(id, payload)
		if err != nil {
			return nil, err
		}
		return &arhatgopb.DoneMsg{}, nil
	default:
		c.logger.D("working on peripheral specific operation")
		// requires peripheral
		handle, ok := handlerMap[kind]
		if !ok {
			c.logger.I("unknown peripheral cmd type", log.Int32("kind", int32(kind)))
			return nil, fmt.Errorf("unknown cmd")
		}

		p, ok := c.getPeripheral(id)
		if !ok {
			return nil, fmt.Errorf("peripheral not found")
		}

		ret, err := handle(p, payload)
		if err != nil {
			return nil, err
		}

		return ret, nil
	}
}

func (c *Handler) handlePeripheralConnect(peripheralID uint64, data []byte) (err error) {
	if _, loaded := c.peripherals.Load(peripheralID); loaded {
		return fmt.Errorf("invalid duplicate peripheral id")
	}

	spec := new(arhatgopb.PeripheralConnectCmd)
	err = spec.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal PeripheralConnectCmd: %w", err)
	}

	p, err := c.impl.Connect(spec.Target, spec.Params, spec.Tls)
	if err != nil {
		return fmt.Errorf("failed to establish connection to peripheral: %w", err)
	}

	defer func() {
		if err != nil {
			p.Close()
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

func (c *Handler) removePeripheral(peripheralID uint64) {
	p, ok := c.getPeripheral(peripheralID)
	if ok {
		p.Close()
	}

	c.peripherals.Delete(peripheralID)
}

func (c *Handler) handlePeripheralOperate(p Peripheral, payload []byte) (proto.Marshaler, error) {
	spec := new(arhatgopb.PeripheralOperateCmd)
	err := spec.Unmarshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal PeripheralOperateCmd: %w", err)
	}

	ret, err := p.Operate(spec.Params, spec.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute operation: %w", err)
	}

	return &arhatgopb.PeripheralOperationResultMsg{Result: ret}, nil
}

func (c *Handler) handlePeripheralMetricsCollect(p Peripheral, payload []byte) (proto.Marshaler, error) {
	spec := new(arhatgopb.PeripheralMetricsCollectCmd)
	err := spec.Unmarshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal PeripheralMetricsCollectCmd: %w", err)
	}

	ret, err := p.CollectMetrics(spec.Params)
	if err != nil {
		return nil, fmt.Errorf("failed to collect peripheral metrics: %w", err)
	}

	return &arhatgopb.PeripheralMetricsMsg{Values: ret}, nil
}
