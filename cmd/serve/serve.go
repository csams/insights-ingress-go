package serve

import (
	"github.com/spf13/cobra"

	"github.com/redhatinsights/insights-ingress-go/internal/common"
	"github.com/redhatinsights/insights-ingress-go/internal/errors"
	"github.com/redhatinsights/insights-ingress-go/internal/logging"
	"github.com/redhatinsights/insights-ingress-go/internal/server"
)

func NewCommand(
	commonOptions *common.Options,
	loggerOptions *logging.Options,
	serverOptions *server.Options) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the ingress server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// common config
			if errs := commonOptions.Complete(); errs != nil {
				return errors.NewAggregate(errs)
			}

			if errs := commonOptions.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			commonConfig, err := common.NewConfig(commonOptions).Complete()
			if err != nil {
				return err
			}

			// logging config
			if errs := loggerOptions.Complete(); errs != nil {
				return errors.NewAggregate(errs)
			}

			if errs := loggerOptions.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			// logging setup
			loggingConfig := logging.NewConfig(loggerOptions, commonConfig).Complete()
			logging.Setup(loggingConfig)

			// server config
			if errs := serverOptions.Complete(); errs != nil {
				return errors.NewAggregate(errs)
			}

			if errs := serverOptions.Validate(); errs != nil {
				return errors.NewAggregate(errs)
			}

			config, err := server.NewConfig(serverOptions, commonConfig).Complete()
			if err != nil {
				return err
			}

			// server startup
			server, err := server.New(config)
			if err != nil {
				return err
			}

			return server.PrepareRun().Run()
		},
	}

	commonOptions.AddFlags(cmd.Flags(), "")
	loggerOptions.AddFlags(cmd.Flags(), "")
	serverOptions.AddFlags(cmd.Flags(), "")

	return cmd
}
