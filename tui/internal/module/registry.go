// Package module — registry.go provides the module registry and dependency resolution.
//
// The Registry is a central store for all module definitions. It provides
// lookup by name, listing, grouping by category, and topological sorting
// for install-order resolution.
//
// # Maps in Go
//
// Go's map type is a hash table — an unordered collection of key-value pairs.
// Maps are reference types (like slices), so passing a map to a function
// doesn't copy the data. The zero value of a map is nil, so you must
// initialize it with make() or a literal before using it.
//
// See: https://go.dev/doc/effective_go#maps
// See: https://go.dev/tour/moretypes/19
//
// # Package-Level Variables
//
// DefaultRegistry is a package-level variable — it's accessible from any file
// in the module package and from external packages that import this one.
// Go initializes package-level variables before main() runs, and init()
// functions run after that, so modules can safely register themselves in init().
//
// See: https://go.dev/doc/effective_go#init
// See: https://go.dev/ref/spec#Package_initialization
package module

import (
	"fmt"
	"sort"
)

// DefaultRegistry is the global registry where all module definitions are
// registered. Module definition files (in the modules/ sub-package) use
// their init() functions to register with this registry.
//
// In Go, package-level variables are initialized before any init() functions
// run in that package. Since NewRegistry() just creates a map, this is safe.
//
// See: https://go.dev/ref/spec#Package_initialization
var DefaultRegistry = NewRegistry()

// Registry stores all known module definitions, indexed by module name.
//
// The struct has a single field: a map from module name (string) to *Module.
// Using a pointer (*Module) means the registry stores references to modules,
// not copies. This allows the registry and callers to share the same Module
// instances.
//
// See: https://go.dev/ref/spec#Map_types
type Registry struct {
	// modules maps module name → module definition pointer.
	// The lowercase first letter makes this field unexported (private).
	// In Go, identifiers starting with uppercase are exported (public),
	// and lowercase are unexported (package-private).
	// See: https://go.dev/doc/effective_go#names
	modules map[string]*Module
}

// NewRegistry creates and returns a new empty Registry.
//
// This is a constructor function — Go's convention for creating initialized
// structs. Go doesn't have constructors like Java/C#, so you write a
// regular function named NewXxx that returns an initialized *Xxx.
//
// The make() built-in function creates and initializes maps, slices, and
// channels. You must use make() for maps before inserting keys.
//
// See: https://go.dev/doc/effective_go#allocation_make
// See: https://go.dev/doc/effective_go#composite_literals
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]*Module),
	}
}

// Register adds a module definition to the registry.
//
// If a module with the same Name already exists, it will be overwritten.
// This method is called from init() functions in the modules/ sub-package.
//
// The (r *Registry) receiver means this is a method on *Registry.
// We use a pointer receiver because Register modifies the registry's internal map.
//
// See: https://go.dev/tour/methods/4 (pointer receivers)
func (r *Registry) Register(m *Module) {
	r.modules[m.Name] = m
}

// Get retrieves a module by name, or nil if not found.
//
// In Go, accessing a map key that doesn't exist returns the zero value for
// the value type. For *Module, the zero value is nil. We use the two-value
// map access pattern (value, ok) internally, but for simplicity this method
// just returns the value (which is nil if not found).
//
// See: https://go.dev/doc/effective_go#maps
func (r *Registry) Get(name string) *Module {
	// The comma-ok idiom: ok is true if the key exists, false otherwise.
	// We discard ok here and just return m (nil if not found).
	m, _ := r.modules[name]
	return m
}

// All returns every registered module, sorted alphabetically by name.
//
// This method demonstrates several Go patterns:
//  1. Iterating over a map with range (which yields key, value pairs)
//  2. Appending to a slice with the built-in append() function
//  3. Sorting with sort.Slice() and a custom less function
//
// # Why Sort?
//
// Go maps are unordered — iterating over them yields keys in a random order
// (intentionally randomized since Go 1). We sort the result so the UI
// displays modules in a consistent, predictable order.
//
// See: https://go.dev/doc/effective_go#for
// See: https://pkg.go.dev/sort#Slice
// See: https://go.dev/blog/maps
func (r *Registry) All() []*Module {
	// make() with a length of 0 and capacity hint for efficiency.
	// The third argument to make() is the capacity — Go will pre-allocate
	// that much space, avoiding reallocations as we append.
	// See: https://go.dev/doc/effective_go#allocation_make
	result := make([]*Module, 0, len(r.modules))

	// range over a map yields (key, value) pairs.
	// The underscore _ discards the key since we only need the value.
	for _, m := range r.modules {
		result = append(result, m)
	}

	// sort.Slice sorts in-place using a "less" function.
	// The less function returns true if result[i] should come before result[j].
	// Here we sort alphabetically by module Name.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// ByCategory groups all modules by their Category field.
//
// Returns a map where keys are category names (e.g., "Shell", "Editor")
// and values are slices of modules in that category, sorted alphabetically.
//
// See: https://go.dev/doc/effective_go#maps
func (r *Registry) ByCategory() map[string][]*Module {
	// Create the result map. We don't know the exact size, so we use make()
	// without a size hint.
	categories := make(map[string][]*Module)

	for _, m := range r.modules {
		// append() adds an element to a slice. If the key doesn't exist
		// in the map yet, categories[m.Category] returns a nil slice,
		// and append() on a nil slice works fine (it allocates a new one).
		categories[m.Category] = append(categories[m.Category], m)
	}

	// Sort each category's module list alphabetically for consistent display.
	for _, mods := range categories {
		sort.Slice(mods, func(i, j int) bool {
			return mods[i].Name < mods[j].Name
		})
	}

	return categories
}

// GetInstallOrder computes a topological installation order for the given
// module names, respecting each module's Dependencies field.
//
// # Topological Sort
//
// A topological sort orders items so that each item comes after all its
// dependencies. For example, if nvim depends on python, python appears
// before nvim in the result. This is essential for installing modules
// in the right order.
//
// # Algorithm: Kahn's Algorithm (BFS-based)
//
// We use Kahn's algorithm, which works like this:
//  1. Build a graph of dependencies and count incoming edges (in-degree)
//  2. Start with nodes that have no incoming edges (no unmet dependencies)
//  3. Process each node: add it to the result, then "remove" its outgoing
//     edges by decrementing the in-degree of its dependents
//  4. Repeat until all nodes are processed
//  5. If some nodes remain unprocessed, there's a circular dependency
//
// # Error Handling in Go
//
// Go uses explicit error returns instead of exceptions. Functions that can
// fail return (result, error). The caller checks if error is nil.
// fmt.Errorf creates a formatted error message.
//
// See: https://go.dev/doc/effective_go#errors
// See: https://go.dev/blog/error-handling-and-go
func (r *Registry) GetInstallOrder(moduleNames []string) ([]string, error) {
	// Step 1: Collect all modules we need to install, including transitive
	// dependencies (dependencies of dependencies).
	//
	// We use a set (map[string]bool) to track which modules we've already
	// visited, preventing infinite loops and duplicate work.
	needed := make(map[string]bool)

	// collectDeps is a recursive closure (anonymous function) that walks the
	// dependency tree. In Go, closures can capture variables from their
	// enclosing scope — here it captures `needed` and `r`.
	//
	// See: https://go.dev/tour/moretypes/25
	// See: https://go.dev/doc/effective_go#functions
	var collectDeps func(name string) error
	collectDeps = func(name string) error {
		// If we've already visited this module, skip it.
		if needed[name] {
			return nil
		}

		m := r.Get(name)
		if m == nil {
			return fmt.Errorf("unknown module: %q", name)
		}

		needed[name] = true

		// Recursively collect dependencies of this module.
		for _, dep := range m.Dependencies {
			if err := collectDeps(dep); err != nil {
				return err
			}
		}
		return nil
	}

	// Collect all needed modules from the requested list.
	for _, name := range moduleNames {
		if err := collectDeps(name); err != nil {
			return nil, err
		}
	}

	// Step 2: Build the dependency graph for Kahn's algorithm.
	//
	// inDegree counts how many unresolved dependencies each module has.
	// dependents maps each module to the list of modules that depend on it.
	inDegree := make(map[string]int)
	dependents := make(map[string][]string)

	for name := range needed {
		// Initialize in-degree. If the key doesn't exist yet, the zero
		// value (0) is used, which is exactly what we want.
		if _, exists := inDegree[name]; !exists {
			inDegree[name] = 0
		}

		m := r.Get(name)
		for _, dep := range m.Dependencies {
			// Only count dependencies that are in our needed set.
			if needed[dep] {
				inDegree[name]++
				dependents[dep] = append(dependents[dep], name)
			}
		}
	}

	// Step 3: Start with modules that have no unresolved dependencies
	// (in-degree == 0). These can be installed first.
	//
	// We use a slice as a queue (FIFO). Go doesn't have a built-in queue
	// type, but slices work fine for this purpose.
	var queue []string
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	// Sort the initial queue alphabetically for deterministic output.
	// Without this, the order would depend on map iteration order (random).
	sort.Strings(queue)

	// Step 4: Process the queue (BFS).
	var result []string
	for len(queue) > 0 {
		// Pop the first element from the queue.
		// queue[0] gets the first element; queue[1:] creates a new slice
		// starting from index 1 (everything after the first element).
		// See: https://go.dev/blog/slices-intro
		current := queue[0]
		queue = queue[1:]

		result = append(result, current)

		// "Remove" outgoing edges: for each module that depends on `current`,
		// decrement its in-degree. If it reaches 0, it's ready to install.
		//
		// We collect newly ready items, sort them, then add to the queue
		// to ensure deterministic ordering.
		var ready []string
		for _, dependent := range dependents[current] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				ready = append(ready, dependent)
			}
		}
		sort.Strings(ready)
		queue = append(queue, ready...)
	}

	// Step 5: Check for circular dependencies.
	// If we couldn't process all modules, some have unresolvable dependencies
	// (they depend on each other in a cycle).
	if len(result) != len(needed) {
		// Find the modules that are stuck in the cycle for a helpful error message.
		var stuck []string
		for name := range needed {
			if inDegree[name] > 0 {
				stuck = append(stuck, name)
			}
		}
		sort.Strings(stuck)
		return nil, fmt.Errorf("circular dependency detected among modules: %v", stuck)
	}

	return result, nil
}
