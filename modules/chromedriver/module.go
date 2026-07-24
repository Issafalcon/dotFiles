package chromedriver

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "chromedriver",
		Icon:        "",
		Description: "ChromeDriver for browser automation",
		Category:    "Utility",
		Website:     "https://chromedriver.chromium.org/",
		Repo:        "https://github.com/nicedoc/chromium",
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "chromedriver --version",
	})
}
