package options

func VenvConfigFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "Label",
			DefaultValue: "",
			Description:  "the pod venv label to be edit, label key is virtual-env",
		},
	}
	return flags
}
