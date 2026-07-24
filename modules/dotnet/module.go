package dotnet

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "dotnet",
		Icon:        "󰪮",
		Description: ".NET SDK with global tools",
		Category:    "Language",
		Website:     "https://dotnet.microsoft.com/",
		Repo:        "https://github.com/dotnet/sdk",
		StowEnabled:   true,
		EstimatedTime: "3m",
		EstimatedSize: "2GB",
		CheckCommand:  "dotnet --version",
		RequiresInput: true,
		ConfigOptions: []module.ConfigOption{
			{
				Name:        "dotnet_versions",
				Description: "Which .NET SDK versions to install",
				Default:     "8.0,9.0,10.0",
				Choices:     []string{"8.0", "9.0", "10.0"},
			},
		},
	})
}
