package command

import (
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	"github.com/alibaba/kt-connect/pkg/kt/command/kubectl"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/alibaba/kt-connect/pkg/kt/service/cluster"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

// NewBirdseyeCommand show a summary of cluster service network
func NewVenvUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "venv-upgrade",
		Short: "Upgrade venv webhook version in cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("too many options specified (%s)", strings.Join(args, ","))
			}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return VenvUpgrade()
		},
		Example: "et venv-upgrade [command options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().Install, opt.InstallFlags())
	return cmd
}

func VenvUpgrade() error {
	ns, err := cluster.Ins().GetNamespace("et-virtual-environment")
	if err != nil {
		return fmt.Errorf("Virtual env webhook does not exist, please install first!")
	}
	if ns.Labels["venvVersion"] == "v0.6.1" {
		log.Info().Msgf("Webhook venv already installed and latest version, Nothing to do")
		return nil
	}
	err = kubectl.UninstallVenvController()
	if err != nil {
		return fmt.Errorf("Webhook venv uninstall failed!")
	}
	err = kubectl.UninstallVenvOperator()
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)
	err = kubectl.InstallVenv()
	if err != nil {
		return fmt.Errorf("Webhook venv install failed!")
	}

	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf(" All looks good, venv upgrade success, v0.6.0 -> v0.6.1")
	log.Info().Msg("---------------------------------------------------------------")

	return nil
}
