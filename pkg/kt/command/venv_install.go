package command

import (
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	"github.com/alibaba/kt-connect/pkg/kt/command/kubectl"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewApplyCommand return new config command
func NewVenvInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "venv-install",
		Short: "Install virtual environment to kubernetes cluster",
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
			return InstallVenv()
		},
		Example: "et venv-install [options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().VenvInstall, opt.VenvInstallFlags())
	return cmd
}

// Connect setup vpn to kubernetes cluster
func InstallVenv() error {
	if opt.Get().VenvInstall.VenvVersion != "v0.6.0" && opt.Get().VenvInstall.VenvVersion != "v0.6.1" {
		return fmt.Errorf("Wrong venv version you prompt!")
	}
	err := kubectl.InstallVenv()
	if err != nil {
		return err
	}

	log.Info().Msgf("virtual environment is ready")
	//log.Info().Msgf("enable injection please label your namespace first")
	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf(" All looks good, now you can access to virtual environment in the kubernetes cluster")
	log.Info().Msg("---------------------------------------------------------------")
	return nil
}
