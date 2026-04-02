package modules

// utility.go registers utility modules: fzf, tmux, ranger, yazi, lazygit,
// git, gh, homebrew, editorconfig, clipboard, wsl-open, libsecret,
// chromedriver, pandoc, quarto, latex, plantuml, buku, zathura, wezterm,
// qmk, snippets.
//
// These are the general-purpose CLI tools, terminal multiplexers, file managers,
// and other utilities that support the development workflow.
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- homebrew ---
	// Homebrew (Linuxbrew) — the package manager used as a dependency by
	// many other modules (lazygit, yazi, terraform, nvim, etc.).
	// This is one of the first modules that should be installed.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "homebrew",
		Icon:        "🍺",
		Description: "Homebrew package manager for Linux",
		Category:    "Utility",
		Website:     "https://brew.sh/",
		Repo:        "https://github.com/Homebrew/brew",
		InstallCommands: []string{
			`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
			"sudo apt-get install -y build-essential",
			"/home/linuxbrew/.linuxbrew/bin/brew install gcc",
		},
		UninstallCommands: []string{
			`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/uninstall.sh)"`,
		},
		StowEnabled:   false,
		EstimatedTime: "3m",
		EstimatedSize: "500MB",
		CheckCommand:  "brew --version",
		RequiresInput: true,
	})

	// --- git ---
	// Git with delta (a syntax-highlighting pager for diffs).
	// Configures delta as the default pager with line numbers and Dracula theme.
	module.DefaultRegistry.Register(&module.Module{
		Name:         "git",
		Icon:         "",
		Description:  "Git with delta diff viewer",
		Category:     "Utility",
		Website:      "https://git-scm.com/",
		Repo:         "https://github.com/git/git",
		Dependencies: []string{"homebrew"},
		InstallCommands: []string{
			"brew install git-delta",
			`git config --global core.pager "delta --dark --paging=never"`,
			`git config --global include.path "~/themes.gitconfig"`,
			`git config --global interactive.diffFilter "delta --color-only"`,
			`git config --global delta.navigate "true"`,
			`git config --global delta.line-numbers "true"`,
			`git config --global delta.side-by-side "false"`,
			`git config --global delta.syntax-theme "Dracula"`,
			`git config --global delta.features "decorations line-numbers zebra-dark"`,
			`git config --global merge.conflictstyle "diff3"`,
			`git config --global credential.helper store`,
		},
		UninstallCommands: []string{
			"brew uninstall git-delta",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "git --version",
	})

	// --- gh ---
	// GitHub CLI — work with GitHub from the command line.
	// Installed from GitHub's official APT repository.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "gh",
		Icon:        "",
		Description: "GitHub CLI",
		Category:    "Utility",
		Website:     "https://cli.github.com/",
		Repo:        "https://github.com/cli/cli",
		InstallCommands: []string{
			"(type -p wget >/dev/null || (sudo apt update && sudo apt install wget -y))",
			"sudo mkdir -p -m 755 /etc/apt/keyrings",
			"wget -nv -O /tmp/githubcli-archive-keyring.gpg https://cli.github.com/packages/githubcli-archive-keyring.gpg",
			"sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg < /tmp/githubcli-archive-keyring.gpg >/dev/null",
			"sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg",
			`echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list >/dev/null`,
			"sudo apt update && sudo apt install -y gh",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y gh",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "50MB",
		CheckCommand:  "gh --version",
	})

	// --- fzf ---
	// fzf — a general-purpose command-line fuzzy finder.
	// Used by zsh for history search, nvim for telescope, and more.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "fzf",
		Icon:        "",
		Description: "Command-line fuzzy finder",
		Category:    "Utility",
		Website:     "https://junegunn.github.io/fzf/",
		Repo:        "https://github.com/junegunn/fzf",
		InstallCommands: []string{
			"sudo apt update && sudo apt install -y fzf",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y fzf",
		},
		StowEnabled:   true,
		EstimatedTime: "15s",
		EstimatedSize: "5MB",
		CheckCommand:  "fzf --version",
	})

	// --- tmux ---
	// tmux — a terminal multiplexer. Allows splitting terminals, detaching
	// sessions, and managing multiple terminal windows from one connection.
	// Also installs TPM (Tmux Plugin Manager).
	module.DefaultRegistry.Register(&module.Module{
		Name:        "tmux",
		Icon:        "",
		Description: "Terminal multiplexer",
		Category:    "Utility",
		Website:     "https://github.com/tmux/tmux/wiki",
		Repo:        "https://github.com/tmux/tmux",
		InstallCommands: []string{
			"sudo apt-get update && sudo apt-get install -y tmux",
			"git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y tmux",
			"rm -rf ~/.tmux/plugins/tpm",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "tmux -V",
	})

	// --- ranger ---
	// ranger — a console file manager with VI key bindings.
	// Includes the ranger_devicons plugin for file icons.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "ranger",
		Icon:        "",
		Description: "Console file manager with VI bindings",
		Category:    "Utility",
		Website:     "https://ranger.github.io/",
		Repo:        "https://github.com/ranger/ranger",
		InstallCommands: []string{
			"sudo apt-get update && sudo apt-get install -y ranger xsel",
			"git clone https://github.com/alexanderjeurissen/ranger_devicons ~/.config/ranger/plugins/ranger_devicons",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y ranger xsel",
			"rm -rf ~/.config/ranger/plugins/ranger_devicons",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "ranger --version",
	})

	// --- yazi ---
	// yazi — a blazing fast terminal file manager written in Rust.
	// Installed via Homebrew along with supporting preview tools
	// (ImageMagick, ffmpeg, fd).
	module.DefaultRegistry.Register(&module.Module{
		Name:         "yazi",
		Icon:         "󰇥",
		Description:  "Blazing fast terminal file manager",
		Category:     "Utility",
		Website:      "https://yazi-rs.github.io/",
		Repo:         "https://github.com/sxyazi/yazi",
		Dependencies: []string{"homebrew"},
		InstallCommands: []string{
			"brew install yazi ImageMagick ffmpeg fd",
		},
		UninstallCommands: []string{
			"brew uninstall yazi ImageMagick ffmpeg fd",
		},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "100MB",
		CheckCommand:  "yazi --version",
	})

	// --- lazygit ---
	// lazygit — a simple terminal UI for Git commands.
	// Installed via Homebrew from Jesse Duffield's tap.
	module.DefaultRegistry.Register(&module.Module{
		Name:         "lazygit",
		Icon:         "",
		Description:  "Simple terminal UI for Git",
		Category:     "Utility",
		Website:      "https://github.com/jesseduffield/lazygit",
		Repo:         "https://github.com/jesseduffield/lazygit",
		Dependencies: []string{"homebrew"},
		InstallCommands: []string{
			"brew install jesseduffield/lazygit/lazygit",
		},
		UninstallCommands: []string{
			"brew uninstall lazygit",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "30MB",
		CheckCommand:  "lazygit --version",
	})

	// --- editorconfig ---
	// EditorConfig — maintains consistent coding styles between editors.
	// Config-only module with no install script; uses stow to link the
	// .editorconfig file.
	module.DefaultRegistry.Register(&module.Module{
		Name:          "editorconfig",
		Icon:          "",
		Description:   "Editor configuration for consistent coding styles",
		Category:      "Utility",
		Website:       "https://editorconfig.org/",
		Repo:          "https://github.com/editorconfig/editorconfig",
		StowEnabled:   true,
		EstimatedTime: "5s",
		EstimatedSize: "1KB",
		CheckCommand:  "test -f $HOME/.editorconfig",
	})

	// --- clipboard ---
	// Clipboard integration — installs xclip on non-WSL systems for
	// clipboard sharing between the terminal and GUI applications.
	// On WSL, win32yank is bundled with Neovim.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "clipboard",
		Icon:        "󰅍",
		Description: "Clipboard integration (xclip)",
		Category:    "Utility",
		Website:     "https://github.com/astrand/xclip",
		Repo:        "https://github.com/astrand/xclip",
		InstallCommands: []string{
			`if ! grep -qEi "(Microsoft|WSL)" /proc/version 2>/dev/null; then sudo apt-get install -y xclip; fi`,
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y xclip",
		},
		StowEnabled:   false,
		EstimatedTime: "15s",
		EstimatedSize: "5MB",
		CheckCommand:  "xclip -version",
	})

	// --- wsl-open ---
	// wsl-open — opens files/URLs from WSL in the default Windows application.
	// Installed as a global npm package.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "wsl-open",
		Icon:        "󰖳",
		Description: "Open files in Windows apps from WSL",
		Category:    "Utility",
		Website:     "https://github.com/4U6U57/wsl-open",
		Repo:        "https://github.com/4U6U57/wsl-open",
		InstallCommands: []string{
			"npm install -g wsl-open",
		},
		UninstallCommands: []string{
			"npm uninstall -g wsl-open",
		},
		StowEnabled:   false,
		EstimatedTime: "15s",
		EstimatedSize: "5MB",
		CheckCommand:  "wsl-open --version",
	})

	// --- libsecret ---
	// libsecret — provides credential storage for Git.
	// On non-WSL: installs gnome-keyring and builds git-credential-libsecret.
	// On WSL: configures git-credential-manager from Windows Git.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "libsecret",
		Icon:        "",
		Description: "Git credential storage via libsecret",
		Category:    "Utility",
		Website:     "https://wiki.gnome.org/Projects/Libsecret",
		Repo:        "https://gitlab.gnome.org/GNOME/libsecret",
		InstallCommands: []string{
			`if grep -qEi "(Microsoft|WSL)" /proc/version 2>/dev/null; then git config --global credential.helper "/mnt/c/Program\ Files/Git/mingw64/bin/git-credential-manager.exe"; else sudo apt install -y gnome-keyring libsecret-1-0 libsecret-1-dev && cd /usr/share/doc/git/contrib/credential/libsecret && sudo make && git config --global credential.helper /usr/share/doc/git/contrib/credential/libsecret/git-credential-libsecret; fi`,
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y gnome-keyring libsecret-1-0 libsecret-1-dev",
		},
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "dpkg -l | grep libsecret-1-0",
	})

	// --- chromedriver ---
	// ChromeDriver — WebDriver for Chrome browser automation and testing.
	// Downloaded from Google's Chrome for Testing CDN.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "chromedriver",
		Icon:        "",
		Description: "ChromeDriver for browser automation",
		Category:    "Utility",
		Website:     "https://chromedriver.chromium.org/",
		Repo:        "https://github.com/nicedoc/chromium",
		InstallCommands: []string{
			`curl "https://storage.googleapis.com/chrome-for-testing-public/130.0.6723.69/linux64/chromedriver-linux64.zip" -o /tmp/chromedriver-linux64.zip`,
			"unzip /tmp/chromedriver-linux64.zip -d /tmp/",
			"sudo mv /tmp/chromedriver-linux64/chromedriver /usr/bin/chromedriver",
			"rm -rf /tmp/chromedriver-linux64 /tmp/chromedriver-linux64.zip",
		},
		UninstallCommands: []string{
			"sudo rm -f /usr/bin/chromedriver",
		},
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "chromedriver --version",
	})

	// --- pandoc ---
	// Pandoc — a universal document converter. Converts between markup
	// formats (Markdown, LaTeX, HTML, DOCX, etc.).
	module.DefaultRegistry.Register(&module.Module{
		Name:        "pandoc",
		Icon:        "",
		Description: "Universal document converter",
		Category:    "Utility",
		Website:     "https://pandoc.org/",
		Repo:        "https://github.com/jgm/pandoc",
		InstallCommands: []string{
			"sudo apt update -y && sudo apt install -y pandoc",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y pandoc",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "50MB",
		CheckCommand:  "pandoc --version",
	})

	// --- quarto ---
	// Quarto — an open-source scientific and technical publishing system.
	// Built on Pandoc, supports Jupyter notebooks, R Markdown, and more.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "quarto",
		Icon:        "󰐗",
		Description: "Scientific and technical publishing",
		Category:    "Utility",
		Website:     "https://quarto.org/",
		Repo:        "https://github.com/quarto-dev/quarto-cli",
		InstallCommands: []string{
			"curl -LO https://github.com/quarto-dev/quarto-cli/releases/download/v1.8.27/quarto-1.8.27-linux-amd64.deb",
			"sudo apt install -y ./quarto-1.8.27-linux-amd64.deb",
			"rm quarto-1.8.27-linux-amd64.deb",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y quarto",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "quarto --version",
	})

	// --- latex ---
	// Full TeX Live installation with latexmk, xdotool, xindy, and
	// Pygments (for minted code highlighting). This is a very large install.
	module.DefaultRegistry.Register(&module.Module{
		Name:         "latex",
		Icon:        "",
		Description:  "Full TeX Live distribution",
		Category:     "Utility",
		Website:      "https://www.latex-project.org/",
		Repo:         "https://github.com/latex3/latex3",
		Dependencies: []string{"python"},
		InstallCommands: []string{
			"sudo apt-get update -y && sudo apt-get install -y texlive-full latexmk xdotool xindy",
			"pip install pygments",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y texlive-full latexmk xdotool xindy",
		},
		StowEnabled:   true,
		EstimatedTime: "10m",
		EstimatedSize: "5GB",
		CheckCommand:  "latex --version",
	})

	// --- plantuml ---
	// PlantUML — generates UML diagrams from text descriptions.
	// Also installs Graphviz for diagram rendering.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "plantuml",
		Icon:        "󰈏",
		Description: "UML diagrams from text descriptions",
		Category:    "Utility",
		Website:     "https://plantuml.com/",
		Repo:        "https://github.com/plantuml/plantuml",
		InstallCommands: []string{
			"sudo apt update -y && sudo apt install -y graphviz plantuml",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y graphviz plantuml",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "plantuml -version",
	})

	// --- buku ---
	// buku — a powerful bookmark manager for the command line.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "buku",
		Icon:        "",
		Description: "Command-line bookmark manager",
		Category:    "Utility",
		Website:     "https://github.com/jarun/buku",
		Repo:        "https://github.com/jarun/buku",
		InstallCommands: []string{
			"sudo apt update && sudo apt install -y buku",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y buku",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "10MB",
		CheckCommand:  "buku --version",
	})

	// --- zathura ---
	// zathura — a highly customizable document viewer with vim-like keybindings.
	// Supports PDF, PostScript, and DjVu.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "zathura",
		Icon:        "",
		Description: "Document viewer with vim-like keybindings",
		Category:    "Utility",
		Website:     "https://pwmt.org/projects/zathura/",
		Repo:        "https://git.pwmt.org/pwmt/zathura",
		InstallCommands: []string{
			"sudo apt-get update -y && sudo apt-get install -y zathura",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y zathura",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "20MB",
		CheckCommand:  "zathura --version",
	})

	// --- wezterm ---
	// WezTerm — a GPU-accelerated terminal emulator and multiplexer.
	// Installed from the official .deb nightly release.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "wezterm",
		Icon:        "",
		Description: "GPU-accelerated terminal emulator",
		Category:    "Utility",
		Website:     "https://wezfurlong.org/wezterm/",
		Repo:        "https://github.com/wez/wezterm",
		InstallCommands: []string{
			"curl -LO https://github.com/wez/wezterm/releases/download/nightly/wezterm-nightly.Ubuntu22.04.deb",
			"sudo apt install -y ./wezterm-nightly.Ubuntu22.04.deb",
			"rm wezterm-nightly.Ubuntu22.04.deb",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y wezterm-nightly",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "wezterm --version",
	})

	// --- qmk ---
	// QMK Firmware — open-source keyboard firmware. Installs QMK CLI in a
	// Python virtualenv and sets up firmware source for keyboard customization.
	module.DefaultRegistry.Register(&module.Module{
		Name:         "qmk",
		Icon:        "󰌌",
		Description:  "QMK keyboard firmware tools",
		Category:     "Utility",
		Website:      "https://qmk.fm/",
		Repo:         "https://github.com/qmk/qmk_firmware",
		Dependencies: []string{"python"},
		InstallCommands: []string{
			`mkdir -p "$HOME/python3/envs"`,
			`if [ ! -d "$HOME/python3/envs/qmk" ]; then cd "$HOME/python3/envs" && python3 -m venv qmk && source "$HOME/python3/envs/qmk/bin/activate" && python3 -m pip install qmk && deactivate; fi`,
		},
		UninstallCommands: []string{
			`rm -rf "$HOME/python3/envs/qmk"`,
		},
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

	// --- snippets ---
	// Snippets — custom code snippet definitions. Config-only module with
	// no install script; uses stow to link snippet files.
	module.DefaultRegistry.Register(&module.Module{
		Name:          "snippets",
		Icon:          "",
		Description:   "Custom code snippet definitions",
		Category:      "Utility",
		Website:       "https://github.com/Issafalcon/dotFiles",
		Repo:          "https://github.com/Issafalcon/dotFiles",
		StowEnabled:   true,
		EstimatedTime: "5s",
		EstimatedSize: "10KB",
		CheckCommand:  "test -d $HOME/.config/snippets",
	})
}
