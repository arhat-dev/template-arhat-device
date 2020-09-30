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

package device

import (
	"crypto/tls"
	"strconv"
	"strings"
	"time"

	"arhat.dev/arhat-proto/arhatgopb"
)

func ConnectDevice(target string, config *Config) (*Device, error) {
	return &Device{target: target, config: config}, nil
}

type Device struct {
	target string
	config *Config
}

func (d *Device) Operate(params map[string]string, data []byte) ([][]byte, error) {
	var ret [][]byte
	for k, v := range params {
		ret = append(ret, []byte(k))
		ret = append(ret, []byte(v))
	}
	ret = append(ret, data)
	return ret, nil
}

func (d *Device) CollectMetrics(param map[string]string) ([]*arhatgopb.DeviceMetricsMsg_Value, error) {
	_ = param

	return []*arhatgopb.DeviceMetricsMsg_Value{
		{Value: float64(d.config.Bar - 1), Timestamp: time.Now().Add(-time.Second).UnixNano()},
		{Value: float64(d.config.Bar + 1), Timestamp: time.Now().Add(time.Second).UnixNano()},
	}, nil
}

func (d *Device) Close() {

}

type Config struct {
	Foo string
	Bar int32
	TLS *tls.Config
}

func ResolveDeviceConfig(
	params map[string]string, tlsConfig *arhatgopb.TLSConfig,
) (*Config, error) {
	ret := &Config{
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

	if tlsConfig != nil {
		var err error
		ret.TLS, err = tlsConfig.GetTLSConfig()
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
