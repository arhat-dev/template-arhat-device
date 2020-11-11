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

package conf

import (
	"arhat.dev/pkg/log"
	"arhat.dev/pkg/tlshelper"
	"github.com/spf13/pflag"

	"ext.arhat.dev/runtimeutil/storageutil"
	"ext.arhat.dev/template-go/pkg/constant"
)

// nolint:lll
type Config struct {
	App AppConfig `json:"app" yaml:"app"`

	Runtime RuntimeConfig            `json:"runtime" yaml:"runtime"`
	Storage storageutil.ClientConfig `json:"storage" yaml:"storage"`
}

type AppConfig struct {
	Log log.ConfigSet `json:"log" yaml:"log"`

	// Endpoint url for the extension server of arhat
	Endpoint string `json:"endpoint" yaml:"endpoint"`

	// affects udp packaet payload size
	MaxDataMessagePayload int `json:"maxDataMessagePayload" yaml:"maxDataMessagePayload"`

	// TLS Client config for the endpoint
	TLS tlshelper.TLSConfig `json:"tls" yaml:"tls"`
}

func FlagsForApp(prefix string, config *AppConfig) *pflag.FlagSet {
	fs := pflag.NewFlagSet("app", pflag.ExitOnError)

	fs.StringVar(&config.Endpoint, prefix+"endpoint",
		constant.DefaultArhatExtensionEndpoint, "set arhat listen address")
	fs.IntVar(&config.MaxDataMessagePayload, prefix+"maxDataMessagePayload", 65535, "")
	fs.AddFlagSet(tlshelper.FlagsForTLSConfig(prefix+"tls.", &config.TLS))
	return fs
}
