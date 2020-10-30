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

	"arhat.dev/libext"
	"arhat.dev/libext/codecpb"
	"arhat.dev/libext/extperipheral"
	"arhat.dev/pkg/log"
	"github.com/spf13/cobra"

	"arhat.dev/template-arhat-ext-go/pkg/conf"
	"arhat.dev/template-arhat-ext-go/pkg/constant"
	"arhat.dev/template-arhat-ext-go/pkg/peripheral"
)

func NewtemplateArhatExtCmd() *cobra.Command {
	var (
		appCtx       context.Context
		configFile   string
		config       = new(conf.Config)
		cliLogConfig = new(log.Config)
	)

	templateArhatExtCmd := &cobra.Command{
		Use:           "template-arhat-ext-go",
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

	flags := templateArhatExtCmd.PersistentFlags()

	flags.StringVarP(&configFile, "config", "c", constant.DefaultTemplateArhatExtConfigFile,
		"path to the templateArhatExt config file")
	flags.AddFlagSet(conf.FlagsForTemplateArhatExt("", &config.TemplateArhatExt))

	return templateArhatExtCmd
}

func run(appCtx context.Context, config *conf.Config) error {
	logger := log.Log.WithName("TemplateArhatExt")

	endpoint := config.TemplateArhatExt.Endpoint

	tlsConfig, err := config.TemplateArhatExt.TLS.GetTLSConfig(false)
	if err != nil {
		return fmt.Errorf("failed to create tls config: %w", err)
	}

	client, err := libext.NewClient(appCtx, endpoint, tlsConfig, libext.ExtensionPeripheral, &codecpb.Codec{})
	if err != nil {
		return fmt.Errorf("failed to create extension client: %w", err)
	}

	ctrl, err := libext.NewController(appCtx, "my-extension-name", log.Log.WithName("controller"),
		extperipheral.NewHandler(log.Log.WithName("handler"), &peripheral.SamplePeripheralConnector{}),
	)
	if err != nil {
		return fmt.Errorf("failed to create extension controller: %w", err)
	}

	err = ctrl.Start()
	if err != nil {
		return fmt.Errorf("failed to start controller: %w", err)
	}

	logger.I("running")
	for {
		select {
		case <-appCtx.Done():
			return nil
		default:
			err = client.ProcessNewStream(ctrl.RefreshChannels())
			if err != nil {
				logger.I("error happened when processing data stream", log.Error(err))
			}
		}
	}
}
