package oracle

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "oracle",
		Icon:        "󰮆",
		Description: "Oracle Cloud Infrastructure CLI",
		Category:    "Cloud",
		Website:     "https://docs.oracle.com/en-us/iaas/tools/oci-cli/latest/",
		Repo:        "https://github.com/oracle/oci-cli",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "~/oci-cli-env/bin/oci --version",
	})
}
