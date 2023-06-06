package options

func VenvInstallFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "VenvVersion",
			DefaultValue: "v0.6.1",
			Description:  "default v0.6.1, The venv version to choose, v0.6.0 or v0.6.1 can be support",
		},
	}
	return flags
}
