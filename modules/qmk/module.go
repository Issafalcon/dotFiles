package qmk

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:         "qmk",
		Icon:        "󰌌",
		Description:  "QMK keyboard firmware tools",
		Category:     "Utility",
		Website:      "https://qmk.fm/",
		Repo:         "https://github.com/qmk/qmk_firmware",
		Dependencies: []string{"python"},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  `test -d "$HOME/python3/envs/qmk"`,
		RequiresInput: true,
		ConfigOptions: []module.ConfigOption{
			{
				Name:        "qmk_variant",
				Description: "QMK firmware variant to set up",
				Default:     "default",
				Choices:     []string{"default", "vial"},
			},
		},
	})
}
