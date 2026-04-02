package modules

// cloud.go registers cloud provider CLI modules: aws, azure, google-cloud, oracle.
//
// These modules install the CLI tools for major cloud providers, enabling
// infrastructure management directly from the terminal.
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- aws ---
	// AWS CLI v2 — installed from the official Amazon zip distribution.
	// The install script downloads the installer, runs it, and cleans up.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "aws",
		Icon:        "",
		Description: "AWS CLI v2",
		Category:    "Cloud",
		Website:     "https://aws.amazon.com/cli/",
		Repo:        "https://github.com/aws/aws-cli",
		InstallCommands: []string{
			`curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "/tmp/awscliv2.zip"`,
			`unzip /tmp/awscliv2.zip -d /tmp/aws-install`,
			"sudo /tmp/aws-install/aws/install",
			"rm -rf /tmp/awscliv2.zip /tmp/aws-install",
		},
		UninstallCommands: []string{
			"sudo rm -rf /usr/local/aws-cli",
			"sudo rm -f /usr/local/bin/aws /usr/local/bin/aws_completer",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "200MB",
		CheckCommand:  "aws --version",
	})

	// --- azure ---
	// Azure CLI with Azure Functions Core Tools.
	// Installed from Microsoft's official APT repository.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "azure",
		Icon:        "󰠅",
		Description: "Azure CLI with Functions Core Tools",
		Category:    "Cloud",
		Website:     "https://docs.microsoft.com/en-us/cli/azure/",
		Repo:        "https://github.com/Azure/azure-cli",
		InstallCommands: []string{
			"curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash",
			// Azure Functions Core Tools
			"curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > microsoft.gpg",
			"sudo mv microsoft.gpg /etc/apt/trusted.gpg.d/microsoft.gpg",
			`sudo sh -c 'echo "deb [arch=amd64] https://packages.microsoft.com/repos/microsoft-ubuntu-$(lsb_release -cs)-prod $(lsb_release -cs) main" > /etc/apt/sources.list.d/dotnetdev.list'`,
			"sudo apt-get update",
			"sudo apt-get install -y azure-functions-core-tools-4",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y azure-cli azure-functions-core-tools-4",
		},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  "az --version",
	})

	// --- google-cloud ---
	// Google Cloud SDK (gcloud, gsutil, bq).
	// Installed by extracting the tarball to $HOME and running the installer.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "google-cloud",
		Icon:        "󱇶",
		Description: "Google Cloud SDK (gcloud CLI)",
		Category:    "Cloud",
		Website:     "https://cloud.google.com/sdk",
		Repo:        "https://github.com/GoogleCloudPlatform/cloud-sdk-docker",
		InstallCommands: []string{
			`tar -C "$HOME" -xzf google-cloud-cli-441.0.0-linux-x86_64.tar.gz`,
			"~/google-cloud-sdk/install.sh",
			"sudo rm -f google-cloud-sdk-441.0.0-linux-x86_64.tar.gz",
		},
		UninstallCommands: []string{
			`rm -rf "$HOME/google-cloud-sdk"`,
		},
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  "gcloud --version",
		RequiresInput: true,
	})

	// --- oracle ---
	// Oracle Cloud Infrastructure (OCI) CLI, installed in a Python venv.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "oracle",
		Icon:        "󰮆",
		Description: "Oracle Cloud Infrastructure CLI",
		Category:    "Cloud",
		Website:     "https://docs.oracle.com/en-us/iaas/tools/oci-cli/latest/",
		Repo:        "https://github.com/oracle/oci-cli",
		InstallCommands: []string{
			"python3 -m venv ~/oci-cli-env",
			"source ~/oci-cli-env/bin/activate && pip install oci-cli && deactivate",
		},
		UninstallCommands: []string{
			"rm -rf ~/oci-cli-env",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "~/oci-cli-env/bin/oci --version",
	})
}
