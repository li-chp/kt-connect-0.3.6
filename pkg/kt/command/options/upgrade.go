package options

import "github.com/alibaba/kt-connect/pkg/kt/util"

func UpgradeFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "ServerUrl",
			DefaultValue: util.UpgradeServerUrl,
			Description:  "Address for ET release server url",
		},
	}
	return flags
}
