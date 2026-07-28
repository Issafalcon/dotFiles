// Package docs embeds user-facing documentation into the binary.
package docs

import _ "embed"

// AddingModules is the guide for creating modules (YAML + scripts).
//
//go:embed ADDING_MODULES.md
var AddingModules string

// Troubleshooting covers common AppImage / modules-path issues.
//
//go:embed TROUBLESHOOTING.md
var Troubleshooting string
