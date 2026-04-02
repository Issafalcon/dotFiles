package modules

// devops.go registers DevOps-related modules: docker, kubernetes, terraform,
// helm, packer, vagrant, localstack.
//
// These modules provide container orchestration, infrastructure-as-code,
// and local cloud emulation tools used in DevOps workflows.
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- docker ---
	// Docker Engine installed from Docker's official APT repository.
	// Includes Docker CE, CLI, containerd, buildx plugin, and compose plugin.
	// The install script first removes any conflicting packages.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "docker",
		Icon:        "",
		Description: "Docker container runtime and tools",
		Category:    "DevOps",
		Website:     "https://www.docker.com/",
		Repo:        "https://github.com/docker/cli",
		InstallCommands: []string{
			// Remove conflicting packages
			"for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove -y $pkg; done",
			// Add Docker's official GPG key and repo
			"sudo apt-get update && sudo apt-get install -y ca-certificates curl",
			"sudo install -m 0755 -d /etc/apt/keyrings",
			"sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc",
			"sudo chmod a+r /etc/apt/keyrings/docker.asc",
			`echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list >/dev/null`,
			// Install Docker packages
			"sudo apt-get update",
			"sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin",
		},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  "docker --version",
	})

	// --- kubernetes ---
	// kubectl CLI and k9s terminal UI for Kubernetes cluster management.
	// Adds the official Kubernetes APT repository for kubectl, then installs
	// k9s via webinstall.dev. Also generates zsh completions.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "kubernetes",
		Icon:        "󱃾",
		Description: "kubectl CLI and k9s terminal UI",
		Category:    "DevOps",
		Website:     "https://kubernetes.io/",
		Repo:        "https://github.com/kubernetes/kubernetes",
		InstallCommands: []string{
			"sudo apt-get update",
			"sudo apt-get install -y apt-transport-https ca-certificates curl gnupg",
			"curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.31/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg",
			"sudo chmod 644 /etc/apt/keyrings/kubernetes-apt-keyring.gpg",
			`echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.31/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list`,
			"sudo chmod 644 /etc/apt/sources.list.d/kubernetes.list",
			"sudo apt-get update",
			"sudo apt-get install -y kubectl",
			// k9s terminal UI for Kubernetes
			"curl -sS https://webinstall.dev/k9s | bash",
			// Generate zsh completions
			`mkdir -p "$HOME/zsh_local/functions" && kubectl completion zsh > "$HOME/zsh_local/functions/_kubectl"`,
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y kubectl",
		},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "100MB",
		CheckCommand:  "kubectl version --client",
	})

	// --- terraform ---
	// HashiCorp Terraform for infrastructure-as-code, plus Terragrunt
	// (wrapper for Terraform) and Mise (for building terragrunt-ls).
	// Uses the official HashiCorp APT repository.
	module.DefaultRegistry.Register(&module.Module{
		Name:         "terraform",
		Icon:         "󱁢",
		Description:  "Terraform and Terragrunt for IaC",
		Category:     "DevOps",
		Website:      "https://www.terraform.io/",
		Repo:         "https://github.com/hashicorp/terraform",
		Dependencies: []string{"homebrew"},
		InstallCommands: []string{
			"sudo apt update && sudo apt install -y gnupg software-properties-common curl",
			"wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg >/dev/null",
			`gpg --no-default-keyring --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg --fingerprint`,
			`echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(grep -oP '(?<=UBUNTU_CODENAME=).*' /etc/os-release || lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list`,
			"sudo apt update",
			"sudo apt-get install -y terraform",
			// Terragrunt via Homebrew
			"brew install terragrunt",
			"terragrunt --install-autocomplete",
			// Mise install (for building terragrunt-ls)
			"sudo apt update -y && sudo apt install -y gpg sudo wget curl",
			"sudo install -dm 755 /etc/apt/keyrings",
			`wget -qO - https://mise.jdx.dev/gpg-key.pub | gpg --dearmor | sudo tee /etc/apt/keyrings/mise-archive-keyring.gpg 1>/dev/null`,
			`echo "deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.gpg arch=amd64] https://mise.jdx.dev/deb stable main" | sudo tee /etc/apt/sources.list.d/mise.list`,
			"sudo apt update",
			"sudo apt install -y mise",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y terraform",
			"brew uninstall terragrunt",
		},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "200MB",
		CheckCommand:  "terraform --version",
	})

	// --- helm ---
	// Helm — the Kubernetes package manager.
	// Installed from the official Helm APT repository.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "helm",
		Icon:        "󱃾",
		Description: "Kubernetes package manager",
		Category:    "DevOps",
		Website:     "https://helm.sh/",
		Repo:        "https://github.com/helm/helm",
		InstallCommands: []string{
			"sudo apt-get install -y curl gpg apt-transport-https",
			`curl -fsSL https://packages.buildkite.com/helm-linux/helm-debian/gpgkey | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg >/dev/null`,
			`echo "deb [signed-by=/usr/share/keyrings/helm.gpg] https://packages.buildkite.com/helm-linux/helm-debian/any/ any main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list`,
			"sudo apt-get update",
			"sudo apt-get install -y helm",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y helm",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "50MB",
		CheckCommand:  "helm version",
	})

	// --- packer ---
	// HashiCorp Packer for building machine images.
	// Uses the official HashiCorp APT repository.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "packer",
		Icon:        "󱁢",
		Description: "HashiCorp Packer for machine images",
		Category:    "DevOps",
		Website:     "https://www.packer.io/",
		Repo:        "https://github.com/hashicorp/packer",
		InstallCommands: []string{
			"sudo apt update && sudo apt install -y gnupg software-properties-common curl",
			"curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -",
			`echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(grep -oP '(?<=UBUNTU_CODENAME=).*' /etc/os-release || lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list`,
			"sudo apt update && sudo apt install -y packer",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y packer",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "packer --version",
	})

	// --- vagrant ---
	// HashiCorp Vagrant for building development environments.
	// Uses the official HashiCorp APT repository.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "vagrant",
		Icon:        "󰜐",
		Description: "HashiCorp Vagrant for dev environments",
		Category:    "DevOps",
		Website:     "https://www.vagrantup.com/",
		Repo:        "https://github.com/hashicorp/vagrant",
		InstallCommands: []string{
			"curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -",
			`sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"`,
			"sudo apt-get update && sudo apt-get install -y vagrant",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y vagrant",
		},
		StowEnabled:   false,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "vagrant --version",
	})

	// --- localstack ---
	// LocalStack — a fully functional local AWS cloud stack for testing.
	// Includes the CLI and awscli-local (awslocal) in a Python venv.
	module.DefaultRegistry.Register(&module.Module{
		Name:         "localstack",
		Icon:         "󰸏",
		Description:  "Local AWS cloud stack for testing",
		Category:     "DevOps",
		Website:      "https://localstack.cloud/",
		Repo:         "https://github.com/localstack/localstack",
		Dependencies: []string{"python"},
		InstallCommands: []string{
			`curl --output localstack-cli-4.3.0-linux-amd64-onefile.tar.gz --location https://github.com/localstack/localstack-cli/releases/download/v4.3.0/localstack-cli-4.3.0-linux-amd64-onefile.tar.gz`,
			"sudo tar xvzf localstack-cli-4.3.0-linux-*-onefile.tar.gz -C /usr/local/bin",
			"rm localstack-cli-4.3.0-linux-*-onefile.tar.gz",
			// Setup awslocal Python venv
			`if [ ! -d "$HOME/python3/envs/awslocal" ]; then mkdir -p "$HOME/python3/envs" && cd "$HOME/python3/envs" && python3 -m venv awslocal && source "$HOME/python3/envs/awslocal/bin/activate" && python3 -m pip install awscli-local && deactivate; fi`,
		},
		UninstallCommands: []string{
			"sudo rm -f /usr/local/bin/localstack",
			`rm -rf "$HOME/python3/envs/awslocal"`,
		},
		StowEnabled:   false,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "localstack --version",
	})
}
