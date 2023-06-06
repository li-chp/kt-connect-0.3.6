package options

func VenvUninstallFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "RemoveAll",
			DefaultValue: false,
			Description:  "Uninstall all resource include webhook and operator, otherwise operator only",
		},
	}
	return flags
}
