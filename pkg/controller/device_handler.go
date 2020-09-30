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
	"fmt"

	"arhat.dev/arhat-proto/arhatgopb"
	"github.com/gogo/protobuf/proto"

	"arhat.dev/template-arhat-device/pkg/device"
)

func (c *Controller) handleDeviceConnect(deviceID uint64, data []byte) (err error) {
	if _, loaded := c.devices.Load(deviceID); loaded {
		return fmt.Errorf("invalid duplicate device id")
	}

	spec := new(arhatgopb.DeviceConnectCmd)
	err = spec.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal DeviceConnectCmd: %w", err)
	}

	cfg, err := device.ResolveDeviceConfig(spec.Params, spec.Tls)
	if err != nil {
		return fmt.Errorf("failed to resolve device config: %w", err)
	}

	dev, err := device.ConnectDevice(spec.Target, cfg)
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

func (c *Controller) getDevice(deviceID uint64) (*device.Device, bool) {
	i, ok := c.devices.Load(deviceID)
	if !ok {
		return nil, false
	}

	dev, ok := i.(*device.Device)
	if !ok {
		c.devices.Delete(deviceID)
	}

	return dev, true
}

func (c *Controller) removeDevice(deviceID uint64) {
	dev, ok := c.getDevice(deviceID)
	if ok {
		dev.Close()
	}

	c.devices.Delete(deviceID)
}

func (c *Controller) handleDeviceOperate(dev *device.Device, payload []byte) (proto.Marshaler, error) {
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

func (c *Controller) handleDeviceMetricsCollect(dev *device.Device, payload []byte) (proto.Marshaler, error) {
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
