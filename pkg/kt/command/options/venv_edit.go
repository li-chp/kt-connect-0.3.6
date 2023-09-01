package options

func VenvEditFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "Label",
			DefaultValue: "",
			Description:  "The deployments template virtual-env to be edit for pod label",
		},
	}
	return flags
}
