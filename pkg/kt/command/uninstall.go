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
func NewUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall plugin from kubernetes cluster",
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
			return Uninstall(args[0])
		},
		Example: "et uninstall <plugins> [options], etck/env",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().Uninstall, opt.UninstallFlags())
	return cmd
}

// Connect setup vpn to kubernetes cluster
func Uninstall(plugin string) error {
	err := kubectl.UninstallEtck()
	if err != nil {
		return err
	}

	log.Info().Msgf("etck-controller is removed")
	return nil
}
