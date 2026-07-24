package python

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "python",
		Icon:        "",
		Description: "Python 3 with pip and venv support",
		Category:    "Language",
		Website:     "https://www.python.org/",
		Repo:        "https://github.com/python/cpython",
		StowEnabled:   false,
		EstimatedTime: "1m",
		EstimatedSize: "150MB",
		CheckCommand:  "python3 --version",
	})
}
