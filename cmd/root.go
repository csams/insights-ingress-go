package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/redhatinsights/insights-ingress-go/cmd/serve"
	"github.com/redhatinsights/insights-ingress-go/internal/common"
	"github.com/redhatinsights/insights-ingress-go/internal/logging"
	"github.com/redhatinsights/insights-ingress-go/internal/server"
)

var (
	rootCmd = &cobra.Command{
		Use: "ingress",
		// this will run before the run func of every sub command
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig()
		},
	}

	options = struct {
		Common *common.Options
		Logger *logging.Options
		Server *server.Options
	}{
		common.NewOptions(),
		logging.NewOptions(),
		server.NewOptions(),
	}
)

func init() {
	viper.SetEnvPrefix("INGRESS")
	viper.AutomaticEnv()

	// allow users to pass a config file as the first argument in the CLI
	rootCmd.PersistentFlags().String("config", "", "config file")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	serveCmd := serve.NewCommand(options.Common, options.Logger, options.Server)
	rootCmd.AddCommand(serveCmd)
	viper.BindPFlags(serveCmd.Flags())
}

func initConfig() error {
	cfgFile := viper.GetString("config")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		return nil
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// unmarshal the config file if it was passed in
	return viper.Unmarshal(options)
}

// Execute runs the root command and provided a catch-all for any unhandled errors.
func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Printf("Unhandled error %v", err)
		os.Exit(1)
	}
}
