package command

import (
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	"github.com/alibaba/kt-connect/pkg/kt/command/kubectl"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewApplyCommand return new config command
func NewVenvUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "venv-uninstall",
		Short: "Uninstall virtual environment from kubernetes cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			//if len(args) == 0 {
			//	return fmt.Errorf("name of plugin to install is required")
			//} else if len(args) > 1 {
			//	return fmt.Errorf("too many plugin names are spcified (%s), should be one", strings.Join(args, ","))
			//} else if args[0] != "etck" {
			//	return fmt.Errorf("only 'etck' can be support!")
			//}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.HideGlobalFlags(cmd)
			return UninstallVenv()
		},
		Example: "et uninstall <plugins> [options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().VenvUninstall, opt.VenvUninstallFlags())
	return cmd
}

// Connect setup vpn to kubernetes cluster
func UninstallVenv() error {
	err := kubectl.UninstallVenvOperator()
	if err != nil {
		return err
	}
	if true == opt.Get().VenvUninstall.RemoveAll {
		kubectl.UninstallVenvController()
	}

	log.Info().Msgf("Virtual environment is removed")
	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf("All looks good, now virtual environment is removed from the kubernetes cluster")
	log.Info().Msg("---------------------------------------------------------------")
	return nil
}
