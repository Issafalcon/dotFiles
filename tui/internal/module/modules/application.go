package modules

// application.go registers application modules: postman, obsidian, drawio,
// googlechrome, seafile, nx, gojira, jira-cli.
//
// These are standalone GUI and CLI applications — API testing tools, note-taking
// apps, diagramming tools, browsers, and project management utilities.
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- postman ---
	// Postman — API development and testing platform.
	// Installed from the official Postman download as a tarball.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "postman",
		Icon:        "󰛳",
		Description: "API development and testing platform",
		Category:    "Application",
		Website:     "https://www.postman.com/",
		Repo:        "https://github.com/postmanlabs/postman-app-support",
		InstallCommands: []string{
			"wget https://dl.pstmn.io/download/latest/linux_64/ -O postman-linux-x64.tar.gz",
			"sudo tar -C /usr/local/ -xzf postman-linux-x64.tar.gz",
			"rm -f postman-linux-x64.tar.gz",
		},
		UninstallCommands: []string{
			"sudo rm -rf /usr/local/Postman",
		},
		StowEnabled:   false,
		EstimatedTime: "1m",
		EstimatedSize: "300MB",
		CheckCommand:  "test -d /usr/local/Postman",
	})

	// --- obsidian ---
	// Obsidian — a knowledge base and note-taking app built on local Markdown files.
	// Installed as an AppImage from the official GitHub releases.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "obsidian",
		Icon:        "󰈙",
		Description: "Knowledge base on local Markdown files",
		Category:    "Application",
		Website:     "https://obsidian.md/",
		Repo:        "https://github.com/obsidianmd/obsidian-releases",
		InstallCommands: []string{
			"sudo apt update -y && sudo apt install -y libnss3-dev",
			"sudo curl -Lo /usr/bin/obsidian https://github.com/obsidianmd/obsidian-releases/releases/download/v1.12.4/Obsidian-1.12.4.AppImage",
			"sudo chmod 777 /usr/bin/obsidian",
		},
		UninstallCommands: []string{
			"sudo rm -f /usr/bin/obsidian",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "test -f /usr/bin/obsidian",
	})

	// --- drawio ---
	// draw.io Desktop — a diagramming application for flowcharts, UML, etc.
	// Installed as an AppImage from the official GitHub releases.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "drawio",
		Icon:        "󰺷",
		Description: "Diagramming application (draw.io)",
		Category:    "Application",
		Website:     "https://www.drawio.com/",
		Repo:        "https://github.com/jgraph/drawio-desktop",
		InstallCommands: []string{
			"sudo curl -Lo /usr/bin/drawio https://github.com/jgraph/drawio-desktop/releases/download/v29.6.1/drawio-x86_64-29.6.1.AppImage",
			"sudo chmod 777 /usr/bin/drawio",
		},
		UninstallCommands: []string{
			"sudo rm -f /usr/bin/drawio",
		},
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "150MB",
		CheckCommand:  "test -f /usr/bin/drawio",
	})

	// --- googlechrome ---
	// Google Chrome — the web browser. Installed from the official .deb package.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "googlechrome",
		Icon:        "",
		Description: "Google Chrome web browser",
		Category:    "Application",
		Website:     "https://www.google.com/chrome/",
		Repo:        "https://github.com/nicedoc/chromium",
		InstallCommands: []string{
			"sudo apt-get update -y && sudo apt-get upgrade -y",
			"sudo wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb",
			"sudo dpkg -i google-chrome-stable_current_amd64.deb",
			"sudo apt install --fix-broken -y",
			"sudo dpkg -i google-chrome-stable_current_amd64.deb",
			"rm google-chrome-stable_current_amd64.deb",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y google-chrome-stable",
		},
		StowEnabled:   false,
		EstimatedTime: "2m",
		EstimatedSize: "300MB",
		CheckCommand:  "google-chrome --version",
	})

	// --- seafile ---
	// Seafile — a self-hosted file sync and share solution.
	// Installed as an AppImage.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "seafile",
		Icon:        "󰅟",
		Description: "Self-hosted file sync and share",
		Category:    "Application",
		Website:     "https://www.seafile.com/",
		Repo:        "https://github.com/haiwen/seafile",
		InstallCommands: []string{
			"sudo curl -Lo /usr/bin/seafile https://s3.eu-central-1.amazonaws.com/download.seadrive.org/Seafile-x86_64-9.0.8.AppImage",
			"sudo chmod 777 /usr/bin/seafile",
		},
		UninstallCommands: []string{
			"sudo rm -f /usr/bin/seafile",
		},
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "100MB",
		CheckCommand:  "test -f /usr/bin/seafile",
	})

	// --- nx ---
	// Nx — a build system for monorepos. Installs the global nx CLI
	// and jq (used for JSON processing in Nx scripts).
	module.DefaultRegistry.Register(&module.Module{
		Name:        "nx",
		Icon:        "󰝖",
		Description: "Nx monorepo build system CLI",
		Category:    "Application",
		Website:     "https://nx.dev/",
		Repo:        "https://github.com/nrwl/nx",
		InstallCommands: []string{
			"sudo apt update && sudo apt install -y jq",
			"npm add --global nx@latest",
		},
		UninstallCommands: []string{
			"npm remove --global nx",
		},
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "50MB",
		CheckCommand:  "nx --version",
	})

	// --- gojira ---
	// go-jira — a command-line Jira client written in Go.
	// Downloaded as a pre-built binary from GitHub releases.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "gojira",
		Icon:        "󰌃",
		Description: "Command-line Jira client (go-jira)",
		Category:    "Application",
		Website:     "https://github.com/go-jira/jira",
		Repo:        "https://github.com/go-jira/jira",
		InstallCommands: []string{
			"sudo curl -fsSLo /usr/bin/jira https://github.com/go-jira/jira/releases/download/v1.0.27/jira-linux-386",
			"sudo chmod 777 /usr/bin/jira",
		},
		UninstallCommands: []string{
			"sudo rm -f /usr/bin/jira",
		},
		StowEnabled:   true,
		EstimatedTime: "15s",
		EstimatedSize: "20MB",
		CheckCommand:  "jira --version",
	})

	// --- jira-cli ---
	// jira-cli — an interactive Jira CLI by Ankitpokhrel.
	// Downloaded as a pre-built tarball from GitHub releases.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "jira-cli",
		Icon:        "󰌃",
		Description: "Interactive Jira CLI tool",
		Category:    "Application",
		Website:     "https://github.com/ankitpokhrel/jira-cli",
		Repo:        "https://github.com/ankitpokhrel/jira-cli",
		InstallCommands: []string{
			"sudo rm -rf /usr/local/jira_1.1.0_linux_x86_64",
			"wget https://github.com/ankitpokhrel/jira-cli/releases/download/v1.1.0/jira_1.1.0_linux_x86_64.tar.gz",
			"sudo tar -C /usr/local/ -xzf jira_1.1.0_linux_x86_64.tar.gz",
			"rm -f jira_1.1.0_linux_x86_64.tar.gz",
		},
		UninstallCommands: []string{
			"sudo rm -rf /usr/local/jira_1.1.0_linux_x86_64",
		},
		StowEnabled:   false,
		EstimatedTime: "30s",
		EstimatedSize: "30MB",
		CheckCommand:  "test -d /usr/local/jira_1.1.0_linux_x86_64",
	})
}
