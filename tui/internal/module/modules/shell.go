package modules

// shell.go registers shell-related modules: zsh, powershell.
//
// These modules provide the core shell environments used across the system.
// zsh is the primary interactive shell; powershell is available for
// cross-platform scripting and Azure/Windows interoperability.
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- zsh ---
	// The Zsh shell with zinit plugin manager and Powerlevel10k theme.
	// Installs zinit from source, sets up italics terminfo, and installs
	// powerline fonts for the prompt theme.
	//
	// Dependencies: none — zsh is a foundational module that other tools
	// (like nvim's terminal integration) may assume is present.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "zsh",
		Icon:        "",
		Description: "Z Shell with zinit plugin manager",
		Category:    "Shell",
		Website:     "https://www.zsh.org/",
		Repo:        "https://github.com/zsh-users/zsh",
		InstallCommands: []string{
			`bash -c "$(curl --fail --show-error --silent --location https://raw.githubusercontent.com/zdharma-continuum/zinit/HEAD/scripts/install.sh)"`,
			`sudo apt-get update && sudo apt-get install -y fonts-powerline powerline`,
			// After stow links .zshrc, switch default shell to zsh (mirrors bootstrap.sh logic).
			`ZSH=$(which zsh) && [ -n "$ZSH" ] && command -v chsh >/dev/null 2>&1 && chsh -s "$ZSH"`,
		},
		UninstallCommands: []string{
			// Switch back to bash before removing zsh (mirrors bootstrap.sh uninstall logic).
			`BASH=$(which bash) && command -v chsh >/dev/null 2>&1 && chsh -s "$BASH"`,
			"sudo apt-get remove -y fonts-powerline powerline",
			"rm -rf $HOME/.local/share/zinit",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "50MB",
		CheckCommand:  "zsh --version",
	})

	// --- powershell ---
	// Microsoft PowerShell for Linux — useful for Azure management and
	// cross-platform scripting. Installed from Microsoft's official APT repo.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "powershell",
		Icon:        "󰨊",
		Description: "Microsoft PowerShell for Linux",
		Category:    "Shell",
		Website:     "https://docs.microsoft.com/en-us/powershell/",
		Repo:        "https://github.com/PowerShell/PowerShell",
		InstallCommands: []string{
			"sudo apt-get update",
			"sudo apt-get install -y wget apt-transport-https software-properties-common",
			`wget -q "https://packages.microsoft.com/config/ubuntu/$(lsb_release -rs)/packages-microsoft-prod.deb"`,
			"sudo dpkg -i packages-microsoft-prod.deb",
			"rm packages-microsoft-prod.deb",
			"sudo apt-get update",
			"sudo apt-get install -y powershell",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y powershell",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "pwsh --version",
		RequiresInput: true,
	})
}
