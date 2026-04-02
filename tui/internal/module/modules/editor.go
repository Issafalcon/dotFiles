package modules

// editor.go registers editor-related modules: nvim, vimspector.
//
// Neovim is the primary editor in this dotfiles setup, with extensive plugin
// configuration. Vimspector provides debugging integration.
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- nvim ---
	// Neovim with a full plugin ecosystem, LSP support, DAP debugging,
	// treesitter syntax highlighting, and telescope fuzzy finding.
	//
	// This is one of the most complex modules — it depends on python
	// (for pynvim), node (for tree-sitter-cli and neovim npm package),
	// go (for Go LSP tools), and homebrew (for ast-grep).
	//
	// The install commands are translated from nvim/install.sh and include:
	//   - System dependencies (ripgrep, cmake, ninja, ctags, etc.)
	//   - Python virtualenv with neovim packages
	//   - Node packages (tree-sitter-cli, neovim)
	//   - Neovim AppImage download
	//   - vscode-js-debug for JavaScript debugging
	//   - FUSE for running AppImages
	//   - ast-grep for structural search/replace
	module.DefaultRegistry.Register(&module.Module{
		Name:         "nvim",
		Icon:         "",
		Description:  "Neovim editor with full plugin ecosystem",
		Category:     "Editor",
		Website:      "https://neovim.io/",
		Repo:         "https://github.com/neovim/neovim",
		Dependencies: []string{"python", "node", "go", "homebrew", "yazi"},
		ExternalDeps: []module.ExternalDep{
			{
				Name:           "ripgrep",
				CheckCommand:   "rg --version",
				InstallCommand: "sudo apt-get install -y ripgrep",
				InstallMethod:  "apt",
			},
			{
				Name:           "cmake",
				CheckCommand:   "cmake --version",
				InstallCommand: "sudo apt-get install -y cmake",
				InstallMethod:  "apt",
			},
		},
		InstallCommands: []string{
			"sudo apt-get update",
			"sudo apt-get install -y ripgrep cmake automake ninja-build silversearcher-ag exuberant-ctags clang sqlite3 libsqlite3-dev",
			"sudo apt-get install -y libjpeg8-dev zlib1g-dev libxtst-dev libxext-dev libtool-bin gettext lua5.1 liblua5.1-dev",
			// Python virtual environment for neovim integration
			`mkdir -p "$HOME/python3/envs"`,
			`if [ ! -d "$HOME/python3/envs/neovim" ]; then cd "$HOME/python3/envs" && python3 -m venv neovim && source "$HOME/python3/envs/neovim/bin/activate" && python3 -m pip install pynvim neovim neovim-remote jupyter ipykernel nbclient nbformat jupyter-cache PyYAML matplotlib plotly pandas && deactivate; fi`,
			"sudo apt install -y python3-pynvim",
			// Node packages for treesitter and neovim integration
			"npm install -g tree-sitter-cli",
			"npm install -g neovim",
			// Download Neovim as AppImage
			"sudo curl -Lo /usr/bin/nvim https://github.com/neovim/neovim/releases/download/v0.12.0/nvim-linux-x86_64.appimage",
			"sudo chmod 777 /usr/bin/nvim",
			// Formatters and linters
			"sudo apt-get install -y chktex",
			// vscode-js-debug for JavaScript/TypeScript debugging
			`mkdir -p "$HOME/.local/share/nvim" && cd "$HOME/.local/share/nvim" && git clone https://github.com/microsoft/vscode-js-debug && cd vscode-js-debug && npm install --legacy-peer-deps && npx gulp vsDebugServerBundle && mv dist out`,
			// FUSE for AppImage support
			"sudo apt install -y libfuse2",
			// ast-grep for grug-far search and replace
			"brew install ast-grep",
		},
		UninstallCommands: []string{
			"sudo rm -f /usr/bin/nvim",
			`rm -rf "$HOME/.local/share/nvim/vscode-js-debug"`,
			`rm -rf "$HOME/python3/envs/neovim"`,
		},
		StowEnabled:   true,
		EstimatedTime: "5m",
		EstimatedSize: "500MB",
		CheckCommand:  "nvim --version",
	})

	// --- vimspector ---
	// Vimspector provides multi-language debugging support for Vim/Neovim.
	// It has no install.sh — configuration is managed entirely through stow.
	module.DefaultRegistry.Register(&module.Module{
		Name:          "vimspector",
		Icon:          "",
		Description:   "Multi-language debugging for Vim/Neovim",
		Category:      "Editor",
		Website:       "https://puremourning.github.io/vimspector-web/",
		Repo:          "https://github.com/puremourning/vimspector",
		StowEnabled:   true,
		EstimatedTime: "10s",
		EstimatedSize: "1MB",
		CheckCommand:  "test -d $HOME/.config/vimspector",
	})
}
