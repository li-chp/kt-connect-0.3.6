package command

import (
	"github.com/alibaba/kt-connect/pkg/kt/command/general"
	opt "github.com/alibaba/kt-connect/pkg/kt/command/options"
	"github.com/alibaba/kt-connect/pkg/kt/command/upgrade"
	"github.com/spf13/cobra"
)

// NewConfigCommand return new config command
func NewUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade ET packages on your system",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.HideGlobalFlags(cmd)
			return Upgrade()
		},
		Example: "et upgrade <sub-command> [options]",
	}

	cmd.SetUsageTemplate(general.UsageTemplate(false))
	opt.SetOptions(cmd, cmd.Flags(), opt.Get().Upgrade, opt.UpgradeFlags())
	return cmd
}

// Connect setup vpn to kubernetes cluster
func Upgrade() error {
	currentVersion := opt.Store.Version
	serverUrl := opt.Get().Upgrade.ServerUrl
	return upgrade.Handle(serverUrl, currentVersion)
}
