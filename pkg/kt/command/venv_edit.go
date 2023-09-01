package command

import (
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/alibaba/kt-connect/pkg/kt/service/cluster"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewBirdseyeCommand show a summary of cluster service network
func NewVenvEditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "venv-edit",
		Short: "Edit deployments templates label virtual-env in cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("name of pod must be spcified")
			}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return VenvEdit(args[0])
		},
		Example: "et venv-edit <deployments> [command options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().VenvEdit, opt.VenvEditFlags())
	return cmd
}

func VenvEdit(deploy string) error {
	if "" == opt.Get().VenvEdit.Label {
		return fmt.Errorf("The parameter label must be spcified!")
	}

	apps, err := cluster.Ins().GetDeployment(deploy, opt.Get().Global.Namespace)
	if err != nil {
		return err
	}
	apps.Spec.Template.Labels["virtual-env"] = opt.Get().VenvEdit.Label
	cluster.Ins().UpdateDeployment(apps)
	log.Info().Msgf("Start to update deployments %s", deploy)

	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf("All looks good, Edit deployment template virtual env label success!")
	log.Info().Msg("---------------------------------------------------------------")
	return nil
}
