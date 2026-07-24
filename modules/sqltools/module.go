package sqltools

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "sqltools",
		Icon:        "",
		Description: "MS SQL Server command-line tools",
		Category:    "Database",
		Website:     "https://docs.microsoft.com/en-us/sql/tools/sqlcmd-utility",
		Repo:        "https://github.com/microsoft/mssql-tools",
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "sqlcmd -?",
		RequiresInput: true,
	})
}
