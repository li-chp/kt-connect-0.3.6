package options

func VenvBirdseyeFlags() []OptionConfig {
	flags := []OptionConfig{
		{
			Target:       "ShowAll",
			DefaultValue: false,
			Description:  "Still show all pods in the namespace",
		},
	}
	return flags
}
