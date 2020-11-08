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
	"fmt"

	"arhat.dev/arhat-proto/arhatgopb"
	"arhat.dev/libext"
	"arhat.dev/libext/codec"
	"arhat.dev/libext/extperipheral"
	"arhat.dev/libext/extruntime"
	"arhat.dev/pkg/log"
	"github.com/spf13/cobra"

	"ext.arhat.dev/template-go/pkg/conf"
	"ext.arhat.dev/template-go/pkg/constant"
	"ext.arhat.dev/template-go/pkg/peripheral"
	"ext.arhat.dev/template-go/pkg/runtime"

	// Add protobuf codec support
	_ "arhat.dev/libext/codec/codecpb"
)

func NewTemplateGoCmd() *cobra.Command {
	var (
		appCtx       context.Context
		configFile   string
		config       = new(conf.Config)
		cliLogConfig = new(log.Config)
	)

	templateGoCmd := &cobra.Command{
		Use:           "template-go",
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

	flags := templateGoCmd.PersistentFlags()

	flags.StringVarP(&configFile, "config", "c", constant.DefaultTemplateGoConfigFile,
		"path to the templateGo config file")
	flags.AddFlagSet(conf.FlagsForTemplateGo("", &config.TemplateGo))

	return templateGoCmd
}

func run(appCtx context.Context, config *conf.Config) error {
	logger := log.Log.WithName("TemplateGo")

	endpoint := config.TemplateGo.Endpoint

	tlsConfig, err := config.TemplateGo.TLS.GetTLSConfig(false)
	if err != nil {
		return fmt.Errorf("failed to create tls config: %w", err)
	}

	go func() {
		// sample peripheral controller

		c := codec.GetCodec(arhatgopb.CODEC_PROTOBUF)
		client, err := libext.NewClient(
			appCtx,
			arhatgopb.EXTENSION_PERIPHERAL,
			"my-extension-name",
			c,
			nil,
			endpoint,
			tlsConfig,
		)
		if err != nil {
			panic(fmt.Errorf("failed to create extension client: %w", err))
		}

		ctrl, err := libext.NewController(appCtx, log.Log.WithName("controller"), c.Marshal,
			extperipheral.NewHandler(
				log.Log.WithName("handler"),
				c.Unmarshal,
				&peripheral.SamplePeripheralConnector{},
			),
		)
		if err != nil {
			panic(fmt.Errorf("failed to create extension controller: %w", err))
		}

		err = ctrl.Start()
		if err != nil {
			panic(fmt.Errorf("failed to start controller: %w", err))
		}

		logger.I("running")
		for {
			select {
			case <-appCtx.Done():
				return
			default:
				err = client.ProcessNewStream(ctrl.RefreshChannels())
				if err != nil {
					logger.I("error happened when processing data stream", log.Error(err))
				}
			}
		}
	}()

	go func() {
		// sample runtime engine

		c := codec.GetCodec(arhatgopb.CODEC_PROTOBUF)
		client, err := libext.NewClient(
			appCtx,
			arhatgopb.EXTENSION_RUNTIME,
			"my-runtime-name",
			c,
			nil,
			endpoint,
			tlsConfig,
		)
		if err != nil {
			panic(fmt.Errorf("failed to create extension client: %w", err))
		}

		ctrl, err := libext.NewController(appCtx, log.Log.WithName("runtime"), c.Marshal,
			extruntime.NewHandler(
				log.Log.WithName("handler"),
				&runtime.SampleRuntime{},
			),
		)
		if err != nil {
			panic(fmt.Errorf("failed to create extension controller: %w", err))
		}

		err = ctrl.Start()
		if err != nil {
			panic(fmt.Errorf("failed to start controller: %w", err))
		}

		logger.I("running")
		for {
			select {
			case <-appCtx.Done():
				return
			default:
				err = client.ProcessNewStream(ctrl.RefreshChannels())
				if err != nil {
					logger.I("error happened when processing data stream", log.Error(err))
				}
			}
		}
	}()

	<-appCtx.Done()

	return nil
}
