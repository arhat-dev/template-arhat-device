package extdevice

import (
	"fmt"
	"sync"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/log"
	"github.com/gogo/protobuf/proto"
)

type cmdHandleFunc func(dev Device, payload []byte) (proto.Marshaler, error)

func NewHandler(logger log.Interface, impl DeviceConnector) *Handler {
	return &Handler{
		logger:  logger,
		impl:    impl,
		devices: new(sync.Map),
	}
}

type Handler struct {
	logger log.Interface

	impl    DeviceConnector
	devices *sync.Map
}

func (c *Handler) HandleCmd(id uint64, kind arhatgopb.CmdType, payload []byte) (proto.Marshaler, error) {
	handlerMap := map[arhatgopb.CmdType]cmdHandleFunc{
		arhatgopb.CMD_DEV_OPERATE:         c.handleDeviceOperate,
		arhatgopb.CMD_DEV_COLLECT_METRICS: c.handleDeviceMetricsCollect,
	}

	var ret proto.Marshaler
	switch kind {
	case arhatgopb.CMD_DEV_CLOSE:
		c.logger.D("removing device")
		c.removeDevice(id)
		ret = &arhatgopb.DoneMsg{}
	case arhatgopb.CMD_DEV_CONNECT:
		c.logger.D("connecting device")
		err := c.handleDeviceConnect(id, payload)
		if err != nil {
			ret = &arhatgopb.ErrorMsg{Description: err.Error()}
		} else {
			ret = &arhatgopb.DoneMsg{}
		}
	default:
		c.logger.D("working on device specific operation")
		// requires device
		handle, ok := handlerMap[kind]
		if !ok {
			c.logger.I("unknown device cmd type", log.Int32("kind", int32(kind)))
			ret = &arhatgopb.ErrorMsg{Description: "unknown cmd"}
			break
		}

		dev, ok := c.getDevice(id)
		if !ok {
			ret = &arhatgopb.ErrorMsg{Description: "device not found"}
			break
		}

		var err error
		ret, err = handle(dev, payload)
		if err != nil {
			ret = &arhatgopb.ErrorMsg{Description: err.Error()}
			break
		}
	}

	return ret, nil
}

func (c *Handler) handleDeviceConnect(deviceID uint64, data []byte) (err error) {
	if _, loaded := c.devices.Load(deviceID); loaded {
		return fmt.Errorf("invalid duplicate device id")
	}

	spec := new(arhatgopb.DeviceConnectCmd)
	err = spec.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal DeviceConnectCmd: %w", err)
	}

	dev, err := c.impl.Connect(spec.Target, spec.Params, spec.Tls)
	if err != nil {
		return fmt.Errorf("failed to establish connection to device: %w", err)
	}

	defer func() {
		if err != nil {
			dev.Close()
		}
	}()

	if _, loaded := c.devices.LoadOrStore(deviceID, dev); loaded {
		return fmt.Errorf("invalid duplicate device")
	}

	return nil
}

func (c *Handler) getDevice(deviceID uint64) (Device, bool) {
	i, ok := c.devices.Load(deviceID)
	if !ok {
		return nil, false
	}

	dev, ok := i.(Device)
	if !ok {
		c.devices.Delete(deviceID)
	}

	return dev, true
}

func (c *Handler) removeDevice(deviceID uint64) {
	dev, ok := c.getDevice(deviceID)
	if ok {
		dev.Close()
	}

	c.devices.Delete(deviceID)
}

func (c *Handler) handleDeviceOperate(dev Device, payload []byte) (proto.Marshaler, error) {
	spec := new(arhatgopb.DeviceOperateCmd)
	err := spec.Unmarshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DeviceOperateCmd: %w", err)
	}

	ret, err := dev.Operate(spec.Params, spec.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute operation: %w", err)
	}

	return &arhatgopb.DeviceOperationResultMsg{Result: ret}, nil
}

func (c *Handler) handleDeviceMetricsCollect(dev Device, payload []byte) (proto.Marshaler, error) {
	spec := new(arhatgopb.DeviceMetricsCollectCmd)
	err := spec.Unmarshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DeviceMetricsCollectCmd: %w", err)
	}

	ret, err := dev.CollectMetrics(spec.Params)
	if err != nil {
		return nil, fmt.Errorf("failed to collect device metrics: %w", err)
	}

	return &arhatgopb.DeviceMetricsMsg{Values: ret}, nil
}
