package cpp

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "cpp",
		Icon:        "",
		Description: "C/C++ toolchain (clang, gcc, cmake)",
		Category:    "Language",
		Website:     "https://clang.llvm.org/",
		Repo:        "https://github.com/llvm/llvm-project",
		ExternalDeps: []module.ExternalDep{
			{
				Name:           "build-essential",
				CheckCommand:   "dpkg -s build-essential",
				InstallCommand: "sudo apt-get install -y build-essential",
				InstallMethod:  "apt",
			},
			{
				Name:           "libssl-dev",
				CheckCommand:   "dpkg -s libssl-dev",
				InstallCommand: "sudo apt-get install -y libssl-dev",
				InstallMethod:  "apt",
			},
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "300MB",
		CheckCommand:  "clang --version",
	})
}
