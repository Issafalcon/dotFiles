package modules

// database.go registers database-related modules: mysql, sqltools, dbeaver.
//
// These modules provide database clients and management tools for working
// with MySQL, SQL Server, and various databases through DBeaver.
//
// See template.go for a detailed explanation of how module registration works.

import "github.com/issafalcon/dotfiles-tui/internal/module"

func init() {
	// --- mysql ---
	// MySQL client only (not the full server). The server is typically
	// run in a Docker container; this just installs the CLI client.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "mysql",
		Icon:        "",
		Description: "MySQL client tools",
		Category:    "Database",
		Website:     "https://www.mysql.com/",
		Repo:        "https://github.com/mysql/mysql-server",
		InstallCommands: []string{
			"sudo apt-get update -y && sudo apt-get install -y mysql-client",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y mysql-client",
		},
		StowEnabled:   true,
		EstimatedTime: "30s",
		EstimatedSize: "50MB",
		CheckCommand:  "mysql --version",
	})

	// --- sqltools ---
	// Microsoft SQL Server command-line tools (sqlcmd, bcp) with ODBC driver.
	// Installed from Microsoft's APT repository.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "sqltools",
		Icon:        "",
		Description: "MS SQL Server command-line tools",
		Category:    "Database",
		Website:     "https://docs.microsoft.com/en-us/sql/tools/sqlcmd-utility",
		Repo:        "https://github.com/microsoft/mssql-tools",
		InstallCommands: []string{
			"sudo curl https://packages.microsoft.com/keys/microsoft.asc | sudo apt-key add -",
			"sudo curl https://packages.microsoft.com/config/ubuntu/20.04/prod.list | sudo tee /etc/apt/sources.list.d/msprod.list",
			"sudo apt-get update",
			"sudo apt-get install -y mssql-tools unixodbc-dev",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y mssql-tools unixodbc-dev",
		},
		StowEnabled:   true,
		EstimatedTime: "1m",
		EstimatedSize: "100MB",
		CheckCommand:  "sqlcmd -?",
		RequiresInput: true,
	})

	// --- dbeaver ---
	// DBeaver Community Edition — a universal database management GUI tool.
	// Installed from the official DBeaver PPA.
	module.DefaultRegistry.Register(&module.Module{
		Name:        "dbeaver",
		Icon:        "",
		Description: "DBeaver universal database manager",
		Category:    "Database",
		Website:     "https://dbeaver.io/",
		Repo:        "https://github.com/dbeaver/dbeaver",
		InstallCommands: []string{
			"sudo add-apt-repository ppa:serge-rider/dbeaver-ce",
			"sudo apt-get update",
			"sudo apt-get install -y dbeaver-ce",
		},
		UninstallCommands: []string{
			"sudo apt-get remove -y dbeaver-ce",
		},
		StowEnabled:   false,
		EstimatedTime: "2m",
		EstimatedSize: "300MB",
		CheckCommand:  "dbeaver --version",
	})
}
