package rust

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "rust",
		Icon:        "",
		Description: "Rust programming language via rustup",
		Category:    "Language",
		Website:     "https://www.rust-lang.org/",
		Repo:        "https://github.com/rust-lang/rust",
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  "rustc --version",
		RequiresInput: true,
	})
}
