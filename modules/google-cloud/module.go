package googlecloud

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "google-cloud",
		Icon:        "󱇶",
		Description: "Google Cloud SDK (gcloud CLI)",
		Category:    "Cloud",
		Website:     "https://cloud.google.com/sdk",
		Repo:        "https://github.com/GoogleCloudPlatform/cloud-sdk-docker",
		StowEnabled:   true,
		EstimatedTime: "2m",
		EstimatedSize: "500MB",
		CheckCommand:  "gcloud --version",
		RequiresInput: true,
	})
}
