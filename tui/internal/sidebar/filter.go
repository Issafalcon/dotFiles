// Package sidebar implements the module list sidebar component for the TUI.
//
// This file contains the filtering and fuzzy-matching logic used to search
// through the list of available modules. The filter is intentionally simple
// — a case-insensitive substring match — which is fast and intuitive for
// small lists of dotfile modules.
//
// # Key Go Concepts Used
//
//   - Structs: Custom composite types that group related fields together.
//     See: https://go.dev/tour/moretypes/2
//   - Slices: Dynamic-length sequences backed by arrays. Slices are one of Go's
//     most-used data structures and are passed by reference (cheap to pass around).
//     See: https://go.dev/tour/moretypes/7
//   - The strings package: Standard library for string manipulation.
//     See: https://pkg.go.dev/strings
//
// See also: https://go.dev/doc/effective_go
package sidebar

import (
	// The strings package provides functions for manipulating UTF-8 encoded strings.
	// We use it for case-insensitive matching via strings.ToLower() and strings.Contains().
	// See: https://pkg.go.dev/strings
	"strings"
)

// ModuleItem represents a single stow module that can be installed.
//
// In Go, structs are the primary way to define custom types with named fields.
// Each field has a name and a type. Exported fields (starting with uppercase)
// are visible outside the package, while unexported fields (lowercase) are private.
//
// See: https://go.dev/ref/spec#Struct_types
// See: https://go.dev/tour/moretypes/2
type ModuleItem struct {
	// Name is the directory name of the stow module (e.g., "nvim", "zsh").
	// This is the primary identifier used for installation and display.
	Name string

	// Icon is the Nerd Font glyph displayed next to the module name.
	// Populated by theme.GetModuleIcon() during initialization.
	// See: https://www.nerdfonts.com/
	Icon string

	// Description is a brief human-readable summary of what this module configures.
	Description string

	// Category groups related modules together (e.g., "editor", "shell", "devops").
	// Used for organizational display and potential future category filtering.
	Category string

	// Installed indicates whether this module is currently stowed (symlinked).
	// In Go, bool fields default to false (the zero value for booleans).
	// See: https://go.dev/ref/spec#The_zero_value
	Installed bool
}

// FuzzyMatch performs a case-insensitive substring match of query against target.
//
// This is a simple "fuzzy" approach — it checks whether the query string appears
// anywhere within the target string, ignoring case. For a small module list
// (typically < 50 items), this is fast and provides an intuitive search experience.
//
// Parameters:
//   - query: The search string entered by the user.
//   - target: The string to search within (e.g., a module name or description).
//
// Returns true if query is a substring of target (case-insensitive).
//
// In Go, functions are first-class citizens. They can be assigned to variables,
// passed as arguments, and returned from other functions.
// See: https://go.dev/tour/moretypes/24
// See: https://pkg.go.dev/strings#Contains
func FuzzyMatch(query, target string) bool {
	// strings.ToLower converts a string to lowercase for case-insensitive comparison.
	// Go strings are immutable UTF-8 byte sequences — ToLower returns a new string.
	// See: https://pkg.go.dev/strings#ToLower
	lowerQuery := strings.ToLower(query)
	lowerTarget := strings.ToLower(target)

	// strings.Contains reports whether the second argument is a substring of the first.
	// See: https://pkg.go.dev/strings#Contains
	return strings.Contains(lowerTarget, lowerQuery)
}

// FilterModules returns the subset of modules whose Name, Description, or Category
// match the given query string. If the query is empty, all modules are returned.
//
// This function creates and returns a new slice rather than modifying the input.
// In Go, slices are reference types — they point to an underlying array. However,
// append() may allocate a new array if the capacity is exceeded, so we always
// capture its return value.
//
// # How append() Works
//
// append(slice, element) adds an element to a slice and returns the updated slice.
// If the underlying array has capacity, it grows in place. Otherwise, Go allocates
// a new, larger array and copies the data. This is why you must always write:
//
//	slice = append(slice, item)
//
// See: https://go.dev/tour/moretypes/15
// See: https://pkg.go.dev/builtin#append
func FilterModules(modules []ModuleItem, query string) []ModuleItem {
	// If the query is empty, return all modules. strings.TrimSpace removes
	// leading and trailing whitespace so that "  " is treated as empty.
	// See: https://pkg.go.dev/strings#TrimSpace
	if strings.TrimSpace(query) == "" {
		return modules
	}

	// Create a new empty slice to hold matching results.
	// var declares a variable; a nil slice is perfectly valid and works with append.
	// The difference between nil slice and empty slice is subtle:
	//   - var s []T        → nil slice (len=0, cap=0, == nil)
	//   - s := []T{}       → empty slice (len=0, cap=0, != nil)
	// Both work identically with append, range, and len.
	// See: https://go.dev/doc/effective_go#allocation_new
	var result []ModuleItem

	// range iterates over a slice, yielding (index, value) pairs.
	// The underscore (_) discards the index since we don't need it.
	// See: https://go.dev/tour/moretypes/16
	for _, module := range modules {
		// Check if the query matches any of the module's searchable fields.
		// The || operator short-circuits: if the first condition is true, the
		// rest are not evaluated. This is standard boolean short-circuit evaluation.
		// See: https://go.dev/ref/spec#Logical_operators
		if FuzzyMatch(query, module.Name) ||
			FuzzyMatch(query, module.Description) ||
			FuzzyMatch(query, module.Category) {
			result = append(result, module)
		}
	}

	return result
}
