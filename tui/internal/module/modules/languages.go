package modules

// languages.go registers programming language modules: node, python, go, rust,
// lua, cpp, dotnet.
//
// Each language module installs the compiler/runtime and associated tooling.
// These are foundational modules — many other tools depend on them (e.g.,
// nvim needs python, node, and go for its plugin ecosystem).
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- node ---
	// Node.js via NVM (Node Version Manager). NVM allows installing and
	// switching between multiple Node.js versions.
	//
	// Dependencies: none — node is a foundational dependency for many tools.
	// The install script clones nvm, installs the latest LTS version, and
	// sets up npm with progress disabled (helps behind VPNs on WSL2).
	module.DefaultRegistry.Register(&module.Module{
		Name:        "node",
		Icon:        "󰎙",
		Description: "Node.js via NVM (Node Version Manager)",
		Category:    "Language",
		Website:     "https://nodejs.org/",
		Repo:        "https://github.com/nvm-sh/nvm",
		InstallCommands: []string{
			`rm -rf "$NVM_DIR" ~/.npm ~/.bower`,
			`cd "$HOME" && git clone https://github.com/nvm-sh/nvm.git .nvm && cd "$HOME/.nvm" && git checkout v0.38.0`,
			`export NVM_DIR="$HOME/.nvm" && [ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh" && nvm install --lts && nvm use --lts && nvm install-latest-npm`,
			"npm set progress=false",
		},
		UninstallCommands: []string{
			`rm -rf "$NVM_DIR" ~/.npm ~/.bower`,
		},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "200MB",
		CheckCommand:  "node --version",
	})

	// --- python ---
	// Python 3 with pip, dev headers, and venv support. Many other modules
	// depend on python for their tooling (nvim, latex, localstack, qmk, etc.).
	module.DefaultRegistry.Register(&module.Module{
		Name:        "python",
		Icon:        "",
		Description: "Python 3 with pip and venv support",
		Category:    "Language",
		Website:     "https://www.python.org/",
		Repo:        "https://github.com/python/cpython",
		InstallCommands: []string{
			"sudo apt update",
			"sudo apt install -y python3 python3-pip python3-dev python3-venv",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y python3-pip python3-dev python3-venv",
		},
		StowEnabled:   false,
		EstimatedTime: "1m",
		EstimatedSize: "150MB",
		CheckCommand:  "python3 --version",
	})

	// --- go ---
	// Go programming language — installed from the official tarball.
	// Go is needed by nvim (for gopls and other Go tools) and various
	// CLI utilities.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "go",
		Icon:        "",
		Description: "Go programming language",
		Category:    "Language",
		Website:     "https://go.dev/",
		Repo:        "https://github.com/golang/go",
		InstallCommands: []string{
			"sudo rm -rf /usr/local/go",
			"wget https://dl.google.com/go/go1.23.3.linux-amd64.tar.gz",
			"sudo tar -C /usr/local/ -xzf go1.23.3.linux-amd64.tar.gz",
			"rm -f go1.23.3.linux-amd64.tar.gz",
		},
		UninstallCommands: []string{
			"sudo rm -rf /usr/local/go",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "500MB",
		CheckCommand:  "go version",
	})

	// --- rust ---
	// Rust via rustup — the official Rust toolchain installer.
	// Rust is used for fast CLI tools (ripgrep, fd, bat, etc.).
	module.DefaultRegistry.Register(&module.Module{
		Name:        "rust",
		Icon:        "",
		Description: "Rust programming language via rustup",
		Category:    "Language",
		Website:     "https://www.rust-lang.org/",
		Repo:        "https://github.com/rust-lang/rust",
		InstallCommands: []string{
			"curl https://sh.rustup.rs -sSf | sh",
			"rustup update stable",
		},
		UninstallCommands: []string{
			"rustup self uninstall",
		},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  "rustc --version",
		RequiresInput: true,
	})

	// --- lua ---
	// Lua 5.4 scripting language with development headers.
	// Used by Neovim for its configuration language and plugins.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "lua",
		Icon:        "",
		Description: "Lua 5.4 scripting language",
		Category:    "Language",
		Website:     "https://www.lua.org/",
		Repo:        "https://github.com/lua/lua",
		InstallCommands: []string{
			"sudo apt update",
			"sudo apt install -y lua5.4 liblua5.4-dev",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y lua5.4 liblua5.4-dev",
		},
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "10MB",
		CheckCommand:  "lua5.4 -v",
	})

	// --- cpp ---
	// C/C++ toolchain — clang, gcc, g++, make, cmake, and OpenSSL headers.
	// Required for compiling native extensions and building C/C++ projects.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "cpp",
		Icon:        "",
		Description: "C/C++ toolchain (clang, gcc, cmake)",
		Category:    "Language",
		Website:     "https://clang.llvm.org/",
		Repo:        "https://github.com/llvm/llvm-project",
		InstallCommands: []string{
			"sudo apt update && sudo apt install -y clang gcc g++ make cmake libssl-dev",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y clang gcc g++ make cmake libssl-dev",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "300MB",
		CheckCommand:  "clang --version",
	})

	// --- dotnet ---
	// .NET SDK (multiple versions) with global tools.
	// Installs .NET 8.0, 9.0, and 10.0 SDKs plus PlantUML class diagram
	// generator, lazydotnet, and dotnet-outdated-tool.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "dotnet",
		Icon:        "󰪮",
		Description: ".NET SDK with global tools",
		Category:    "Language",
		Website:     "https://dotnet.microsoft.com/",
		Repo:        "https://github.com/dotnet/sdk",
		InstallCommands: []string{
			"sudo apt update -y",
			"sudo add-apt-repository ppa:dotnet/backports",
			"sudo apt install -y dotnet-sdk-8.0",
			"sudo apt install -y dotnet-sdk-9.0",
			"sudo apt install -y dotnet-sdk-10.0",
			"dotnet tool install --global PlantUmlClassDiagramGenerator",
			"dotnet tool install --global lazydotnet",
			"dotnet tool install --global dotnet-outdated-tool",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y dotnet-sdk-8.0 dotnet-sdk-9.0 dotnet-sdk-10.0",
		},
		StowEnabled:   true,
		EstimatedTime: "3m",
		EstimatedSize: "2GB",
		CheckCommand:  "dotnet --version",
		RequiresInput: true,
		ConfigOptions: []module.ConfigOption{
			{
				Name:        "dotnet_versions",
				Description: "Which .NET SDK versions to install",
				Default:     "8.0,9.0,10.0",
				Choices:     []string{"8.0", "9.0", "10.0"},
			},
		},
	})
}
