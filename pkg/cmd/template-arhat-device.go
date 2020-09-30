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

package cmd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/pkg/iohelper"
	"arhat.dev/pkg/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"arhat.dev/template-arhat-device/pkg/conf"
	"arhat.dev/template-arhat-device/pkg/constant"
	"arhat.dev/template-arhat-device/pkg/controller"
)

func NewTemplateArhatDeviceCmd() *cobra.Command {
	var (
		appCtx       context.Context
		configFile   string
		config       = new(conf.TemplateArhatDeviceConfig)
		cliLogConfig = new(log.Config)
	)

	templateArhatDeviceCmd := &cobra.Command{
		Use:           "template-arhat-device",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use == "version" {
				return nil
			}

			var err error
			appCtx, err = conf.ReadConfig(cmd, &configFile, cliLogConfig, config)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(appCtx, config)
		},
	}

	flags := templateArhatDeviceCmd.PersistentFlags()

	flags.StringVarP(&configFile, "config", "c", constant.DefaultTemplateArhatDeviceConfigFile,
		"path to the templateArhatDevice config file")
	flags.AddFlagSet(conf.FlagsForTemplateArhatDevice("", &config.TemplateArhatDevice))

	return templateArhatDeviceCmd
}

func run(appCtx context.Context, config *conf.TemplateArhatDeviceConfig) error {
	logger := log.Log.WithName("TemplateArhatDevice")

	endpoint := config.TemplateArhatDevice.Endpoint

	tlsConfig, err := config.TemplateArhatDevice.TLS.GetTLSConfig(false)
	if err != nil {
		return fmt.Errorf("failed to create tls config: %w", err)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint url: %w", err)
	}

	ctrl, err := controller.NewController(appCtx)
	if err != nil {
		return fmt.Errorf("failed to create device controller: %w", err)
	}

	err = ctrl.Start()
	if err != nil {
		return fmt.Errorf("failed to start controller: %w", err)
	}

	for {
		select {
		case <-appCtx.Done():
			return nil
		default:
			cmdCh, msgCh := ctrl.RefreshChannels()
			err = processNewJSONStream(appCtx, u.Scheme, u.Host, tlsConfig, cmdCh, msgCh)
			if err != nil {
				logger.I("error happened when processing json stream", log.Error(err))
			}
		}
	}
}

func processNewJSONStream(
	appCtx context.Context, network, addr string, tlsConfig *tls.Config,
	cmdCh chan<- *arhatgopb.Cmd, msgCh <-chan *arhatgopb.Msg,
) error {
	client, cleanup, err := createHTTPClient(appCtx, network, addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial arhat extension endpoint for the first time: %w", err)
	}

	pr, pw := iohelper.Pipe()

	resp, err := client.Post("/ext/devices", "application/json", pr)
	if err != nil {
		return fmt.Errorf("failed to start sync loop for the first time: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	wg, ctx := errgroup.WithContext(appCtx)

	wg.Go(func() error {
		for msg := range msgCh {
			data, err2 := json.Marshal(msg)
			if err2 != nil {
				return fmt.Errorf("failed to marshal json data: %w", err2)
			}

			_, err2 = pw.Write(data)
			if err2 != nil {
				return fmt.Errorf("failed to send msg: %w", err2)
			}
		}

		_ = pw.Close()

		return io.EOF
	})

	wg.Go(func() error {
		defer func() {
			close(cmdCh)
		}()

		for {
			dec := json.NewDecoder(resp.Body)
			cmd := new(arhatgopb.Cmd)
			err2 := dec.Decode(cmd)
			if err2 != nil {
				if err2 != io.EOF {
					return fmt.Errorf("failed to decode cmd: %w", err2)
				}
				return nil
			}

			select {
			case cmdCh <- cmd:
			case <-ctx.Done():
				return nil
			}
		}
	})

	err = wg.Wait()
	if err != nil && err != io.EOF {
		return err
	}

	cleanup()

	return nil
}

func createHTTPClient(
	parent context.Context,
	network, addr string,
	tlsConfig *tls.Config,
) (_ *http.Client, cleanup func(), _ error) {
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(parent, network, addr)
	if err != nil {
		return nil, nil, err
	}

	return &http.Client{
			Transport: &http.Transport{
				Proxy: nil,
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return conn, err
				},
				TLSClientConfig:   tlsConfig,
				ForceAttemptHTTP2: true,
			},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}, func() {
			_ = conn.Close()
		}, nil
}
