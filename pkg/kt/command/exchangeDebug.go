package command

import (
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/command/connect"
	"github.com/alibaba/kt-connect/pkg/kt/command/exchange"
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/alibaba/kt-connect/pkg/kt/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

// NewExchangeCommand return new exchange command
func NewExchangeDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edebug",
		Short: "Combines connect and exchange together",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("name of service to exchange is required")
			} else if len(args) > 1 {
				return fmt.Errorf("too many service name are spcified (%s), should be one", strings.Join(args, ","))
			}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExchangeDebug(args[0])
		},
		Example: "et edebug <service-name> [command options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(true))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().ExchangeDebug, opt.ExchangeDebugFlags())
	return cmd
}

//Exchange exchange kubernetes workload
func ExchangeDebug(resourceName string) error {
	*opt.Get().Connect = opt.Get().ExchangeDebug.ConnectOptions
	*opt.Get().Exchange = opt.Get().ExchangeDebug.ExchangeOptions
	ch, err := general.SetupProcess(util.ComponentExchangeDebug)
	if err != nil {
		return err
	}
	if !opt.Get().Connect.SkipCleanup {
		go silenceCleanup()
	}

	log.Info().Msgf("Using %s mode", opt.Get().Connect.ConnectMode)
	if opt.Get().Connect.ConnectMode == util.ConnectModeTun2Socks {
		err = connect.ByTun2Socks()
	} else if opt.Get().Connect.ConnectMode == util.ConnectModeShuttle {
		err = connect.BySshuttle()
	} else {
		err = fmt.Errorf("invalid connect mode: '%s', supportted mode are %s, %s", opt.Get().Connect.ConnectMode,
			util.ConnectModeTun2Socks, util.ConnectModeShuttle)
	}
	if err != nil {
		return err
	}
	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf(" All looks good, now you can access to resources in the kubernetes cluster")
	log.Info().Msg("---------------------------------------------------------------")

	if opt.Get().Exchange.SkipPortChecking {
		if port := util.FindBrokenLocalPort(opt.Get().Exchange.Expose); port != "" {
			return fmt.Errorf("no application is running on port %s", port)
		}
	}

	log.Info().Msgf("Using %s mode", opt.Get().Exchange.ExchangeMode)
	if opt.Get().Exchange.ExchangeMode == util.ExchangeModeScale {
		err = exchange.ByScale(resourceName)
	} else if opt.Get().Exchange.ExchangeMode == util.ExchangeModeEphemeral {
		err = exchange.ByEphemeralContainer(resourceName)
	} else if opt.Get().Exchange.ExchangeMode == util.ExchangeModeSelector {
		err = exchange.BySelector(resourceName)
	} else {
		err = fmt.Errorf("invalid exchange method '%s', supportted are %s, %s, %s", opt.Get().Exchange.ExchangeMode,
			util.ExchangeModeSelector, util.ExchangeModeScale, util.ExchangeModeEphemeral)
	}
	if err != nil {
		return err
	}
	resourceType, realName := toTypeAndName(resourceName)
	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf(" Now all request to %s '%s' will be redirected to local", resourceType, realName)
	log.Info().Msg("---------------------------------------------------------------")

	// watch background process, clean the workspace and exit if background process occur exception
	s := <-ch
	log.Info().Msgf("Terminal Signal is %s", s)
	return nil
}
