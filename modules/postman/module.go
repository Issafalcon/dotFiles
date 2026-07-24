package postman

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "postman",
		Icon:        "󰛳",
		Description: "API development and testing platform",
		Category:    "Application",
		Website:     "https://www.postman.com/",
		Repo:        "https://github.com/postmanlabs/postman-app-support",
		StowEnabled:   false,
		EstimatedTime: "1m",
		EstimatedSize: "300MB",
		CheckCommand:  "test -d /usr/local/Postman",
	})
}
