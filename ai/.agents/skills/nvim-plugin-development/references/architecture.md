# Architecture & Directory Layout

## Canonical structure

```
my-plugin/
├── lua/
│   └── my-plugin/
│       ├── init.lua          ← public API surface
│       ├── types.lua         ← LuaCATS @class definitions only
│       ├── health.lua        ← :checkhealth my-plugin
│       ├── _core/            ← infrastructure (no feature logic)
│       │   ├── configuration.lua
│       │   ├── log.lua
│       │   └── cmdparse.lua
│       └── _<domain>/        ← feature modules (e.g. _lsp/, _ui/)
│           ├── content.lua
│           ├── handler.lua
│           └── mappings.lua
├── plugin/
│   └── my_plugin.lua         ← lazy-loaded entry; double-source guard
├── spec/
│   ├── minimal_init.lua      ← headless test bootstrap
│   └── my-plugin/
│       └── *_spec.lua
├── .busted                   ← busted config
├── .luacheckrc               ← luacheck config
├── .luarc.json               ← lua-language-server config
├── Makefile
└── README.md
```

## Module responsibilities

### `init.lua` — public API

- `M.setup(opts)` — single opts table, calls `configuration.set(opts)`,
  registers autocmds and handler overrides.
- `M.on_attach(client, bufnr)` — manual per-buffer setup for users who
  don't want the auto-attach `LspAttach` approach.
- Exported fields (e.g. `M.handler`) for advanced users.
- **Do not** put implementation logic here; delegate to domain modules.

### `plugin/<name>.lua` — lazy entry point

```lua
if vim.g.loaded_my_plugin_plugin then return end
vim.g.loaded_my_plugin_plugin = true

-- Register user commands; do NOT call setup() here.
-- Users call setup() in their own config.
require("my-plugin._core.cmdparse").create("MyPlugin", {
  signature = function(_) ... end,
  toggle    = function(_) ... end,
})
```

### `_core/configuration.lua` — settings

- Single source of truth for all defaults.
- `M.set(opts)` — `vim.tbl_deep_extend("force", current, opts)`.
- `M.get()` — returns current merged config.
- `M.reset()` — restores defaults; **must** nil `vim.g.loaded_<plugin>` guard.
  Call in test `before_each` for isolation.
- `M.initialize_if_needed()` — once-per-session side effects (e.g. setting
  the global LSP handler).

### `_core/log.lua` — logger

- Wraps `vim.notify` with log-level filtering.
- No external dependencies.
- Level set from `config.log_level`; defaults to `"warn"`.

### `_core/cmdparse.lua` — command dispatcher

- `M.create(name, subcommands, opts)` — registers a `nvim_create_user_command`
  with tab-completion over the `subcommands` keys.
- Uses `error()` (not `vim.notify`) for unknown subcommands so callers
  can catch them with `pcall`.

### `types.lua` — type declarations

```lua
---@class MyPluginConfig
---@field ui MyPluginUiConfig
---@field keymaps MyPluginKeymaps
---@field display_automatically boolean

---@class MyPluginUiConfig
---@field border string
---@field silent boolean
---@field zindex integer
```
Only `@class`/`@field` here — no runtime code.

### `health.lua` — checkhealth

Neovim auto-discovers `lua/<plugin>/health.lua` when `:checkhealth <plugin>`
is run; no explicit registration needed.

```lua
local M = {}
function M.check()
  vim.health.start("my-plugin")
  -- version
  local v = vim.version()
  if v.minor < 11 then
    vim.health.error("Neovim 0.11+ required, got " .. tostring(v))
  else
    vim.health.ok("Neovim " .. tostring(v))
  end
  -- handler override
  if vim.lsp.handlers["textDocument/signatureHelp"] == require("my-plugin").handler then
    vim.health.ok("Native signatureHelp handler overridden")
  else
    vim.health.warn("Native handler not overridden (override_native_handler = false?)")
  end
end
return M
```

## Naming conventions

| Item | Convention | Example |
|------|-----------|---------|
| Plugin module path | kebab-case | `lua/my-plugin/` |
| Plugin entry var | snake_case with underscores | `vim.g.loaded_my_plugin` |
| Lua module requires | match directory | `require("my-plugin._core.log")` |
| Public functions | `M.verb_noun` | `M.open_signature()` |
| Private functions | `_verb_noun` (local) | `local function _build_lines()` |
| Test files | `*_spec.lua` | `handler_spec.lua` |

## `Makefile` targets

```makefile
PLENARY_PATH ?= $(shell find ~/.local/share/nvim -name plenary.nvim -maxdepth 6 -type d 2>/dev/null | head -1)

.PHONY: test
test:
	nvim --headless --clean --noplugin -u spec/minimal_init.lua \
	  -c "lua vim.cmd([[PlenaryBustedDirectory spec/my-plugin { minimal_init = 'spec/minimal_init.lua' }]])"

.PHONY: lint
lint:
	luacheck lua/ spec/ --config .luacheckrc
```

## `.luarc.json` (lua-language-server)

```json
{
  "$schema": "https://raw.githubusercontent.com/LuaLS/vscode-lua/master/setting/schema.json",
  "runtime": { "version": "LuaJIT" },
  "workspace": {
    "library": [
      "$VIMRUNTIME",
      "${3rd}/luv/library"
    ],
    "checkThirdParty": false
  },
  "diagnostics": { "globals": ["vim", "describe", "it", "before_each", "after_each", "assert"] }
}
```
