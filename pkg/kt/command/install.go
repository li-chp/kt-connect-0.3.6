package command

import (
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	"github.com/alibaba/kt-connect/pkg/kt/command/kubectl"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

// NewApplyCommand return new config command
func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install plugin to kubernetes cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("name of plugin to install is required")
			} else if len(args) > 1 {
				return fmt.Errorf("too many plugin names are spcified (%s), should be one", strings.Join(args, ","))
			} else if args[0] != "etck" {
				return fmt.Errorf("only 'etck' can be support!")
			}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.HideGlobalFlags(cmd)
			return Install(args[0])
		},
		Example: "et install <plugins> [options], etck/env",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().Install, opt.InstallFlags())
	return cmd
}

// Connect setup vpn to kubernetes cluster
func Install(plugin string) error {
	err := kubectl.InstallEtck()
	if err != nil {
		return err
	}

	log.Info().Msgf("etck-controller is ready")
	log.Info().Msgf("enable injection please label your namespace first")
	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf("kubectl label namespace your-namespace etck-injection=enabled --overwrite")
	log.Info().Msg("---------------------------------------------------------------")
	return nil
}
