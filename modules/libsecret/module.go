package libsecret

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "libsecret",
		Icon:        "",
		Description: "Git credential storage via libsecret",
		Category:    "Utility",
		Website:     "https://wiki.gnome.org/Projects/Libsecret",
		Repo:        "https://gitlab.gnome.org/GNOME/libsecret",
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "dpkg -l | grep libsecret-1-0",
	})
}
