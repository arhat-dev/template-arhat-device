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

package peripheral

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/libext/extperipheral"
)

var _ extperipheral.PeripheralConnector = (*SamplePeripheralConnector)(nil)

type SamplePeripheralConnector struct{}

func (c *SamplePeripheralConnector) Connect(
	ctx context.Context,
	target string,
	params map[string]string,
	tlsConfig *arhatgopb.TLSConfig,
) (extperipheral.Peripheral, error) {
	config, err := resolvePeripheralConfig(params, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve peripheral config")
	}

	return &SamplePeripheral{
		target: target,
		config: config,
	}, nil
}

var _ extperipheral.Peripheral = (*SamplePeripheral)(nil)

type SamplePeripheral struct {
	target string
	config *SampleConfig
}

func (d *SamplePeripheral) Operate(
	ctx context.Context,
	params map[string]string,
	data []byte,
) ([][]byte, error) {
	var ret [][]byte
	for k, v := range params {
		ret = append(ret, []byte(k))
		ret = append(ret, []byte(v))
	}
	ret = append(ret, data)
	return ret, nil
}

func (d *SamplePeripheral) CollectMetrics(
	ctx context.Context,
	param map[string]string,
) ([]*arhatgopb.PeripheralMetricsMsg_Value, error) {
	_ = param

	return []*arhatgopb.PeripheralMetricsMsg_Value{
		{Value: float64(d.config.Bar - 1), Timestamp: time.Now().Add(-time.Second).UnixNano()},
		{Value: float64(d.config.Bar + 1), Timestamp: time.Now().Add(time.Second).UnixNano()},
	}, nil
}

func (d *SamplePeripheral) Close(ctx context.Context) {}

type SampleConfig struct {
	Foo string
	Bar int32
	TLS *tls.Config
}

func resolvePeripheralConfig(
	params map[string]string, tlsConfig *arhatgopb.TLSConfig,
) (*SampleConfig, error) {
	ret := &SampleConfig{
		Foo: "<nil>",
		Bar: 0,
		TLS: nil,
	}

	if len(params) != 0 {
		ret.Foo = params["foo"]
	}

	for k, v := range params {
		switch strings.ToLower(k) {
		case "foo":
			ret.Foo = v
		case "bar":
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return nil, err
			}
			ret.Bar = int32(i)
		}
	}

	_ = tlsConfig
	return ret, nil
}
