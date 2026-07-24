package powershell

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "powershell",
		Icon:        "󰨊",
		Description: "Microsoft PowerShell for Linux",
		Category:    "Shell",
		Website:     "https://docs.microsoft.com/en-us/powershell/",
		Repo:        "https://github.com/PowerShell/PowerShell",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "pwsh --version",
		RequiresInput: true,
	})
}
