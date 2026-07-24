package googlechrome

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "googlechrome",
		Icon:        "",
		Description: "Google Chrome web browser",
		Category:    "Application",
		Website:     "https://www.google.com/chrome/",
		Repo:        "https://github.com/nicedoc/chromium",
		StowEnabled:   false,
		EstimatedTime: "2m",
		EstimatedSize: "300MB",
		CheckCommand:  "google-chrome --version",
	})
}
