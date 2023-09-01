package options

func VenvInstallFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "VenvVersion",
			DefaultValue: "v0.6.2",
			Description:  "default v0.6.2, The venv version to choose",
		},
		{
			Target:       "EnvHeader",
			DefaultValue: "cmss-env-mark",
			Description:  "The alternative transparent HTTP headers",
		},
		//{
		//	Target:       "EnvHeaderAlias",
		//	DefaultValue: "",
		//	Description:  "Add additional alternative transparent HTTP headers (usually not required)",
		//},
	}
	return flags
}
