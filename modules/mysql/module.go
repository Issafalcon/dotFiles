package mysql

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	module.DefaultRegistry.Register(&module.Module{
		Name:        "mysql",
		Icon:        "",
		Description: "MySQL client tools",
		Category:    "Database",
		Website:     "https://www.mysql.com/",
		Repo:        "https://github.com/mysql/mysql-server",
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "50MB",
		CheckCommand:  "mysql --version",
	})
}
