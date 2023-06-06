package command

import (
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/alibaba/kt-connect/pkg/kt/service/cluster"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

// NewBirdseyeCommand show a summary of cluster service network
func NewVenvBirdseyeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "venv-birdseye",
		Short: "Show venv label in cluster",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("too many options specified (%s)", strings.Join(args, ","))
			}
			return general.Prepare()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return VenvBirdseye()
		},
		Example: "et venv-birdseye [command options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().VenvBirdseye, opt.VenvBirdseyeFlags())
	return cmd
}

func VenvBirdseye() error {
	t := table.Table{}
	header := table.Row{"Pod Name", "Virtual Env Label"}
	t.AppendHeader(header)
	t.SetAutoIndex(true)

	pods, err := cluster.Ins().GetAllPodInNamespace(opt.Get().Global.Namespace)
	if err != nil {
		return err
	}

	for _, p := range pods.Items {
		podName := p.Name
		venv, _ := p.Labels["virtual-env"]
		if venv == "" && !opt.Get().VenvBirdseye.ShowAll {
			continue
		}
		t.AppendRow(table.Row{podName, venv})
	}
	log.Info().Msgf("---- Venv label pods in namespace [%s] ---- \n%s", opt.Get().Global.Namespace, t.Render())

	return nil
}
