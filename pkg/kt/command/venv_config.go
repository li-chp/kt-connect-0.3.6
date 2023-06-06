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
func NewVenvConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "venv-config",
		Short: "Config one or more pod venv label in cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("name of pod must be spcified")
			}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return VenvConfig(args)
		},
		Example: "et venv-config [command options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().VenvConfig, opt.VenvConfigFlags())
	return cmd
}

func VenvConfig(pods []string) error {
	if "" == opt.Get().VenvConfig.Label {
		return fmt.Errorf("The parameter label must be spcified!")
	}

	for _, name := range pods {
		pod, err := cluster.Ins().GetPod(name, opt.Get().Global.Namespace)
		if err != nil {
			return err
		}
		pod.Labels["virtual-env"] = opt.Get().VenvConfig.Label
		log.Info().Msgf("Start to update pod %s", name)
		cluster.Ins().UpdatePod(pod)
	}

	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf("All looks good, Config pod virtual env label success!")
	log.Info().Msg("---------------------------------------------------------------")
	return nil
}
