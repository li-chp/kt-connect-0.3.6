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
func NewVenvCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "venv-clean",
		Short: "Clean one or more pod venv label in cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("name of pod must be spcified")
			}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return VenvClean(args)
		},
		Example: "et venv-config [pod..]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().VenvConfig, opt.InstallFlags())
	return cmd
}

func VenvClean(pods []string) error {
	for _, name := range pods {
		pod, err := cluster.Ins().GetPod(name, opt.Get().Global.Namespace)
		if err != nil {
			return err
		}
		delete(pod.Labels, "virtual-env")
		log.Info().Msgf("Start to update pod %s", name)
		cluster.Ins().UpdatePod(pod)
	}

	log.Info().Msg("---------------------------------------------------------------")
	log.Info().Msgf("All looks good, Clean pod venv label success!")
	log.Info().Msg("---------------------------------------------------------------")
	return nil
}
