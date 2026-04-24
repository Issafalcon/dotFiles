---
name: nvim-plugin-development
description: >
  Write, refactor, or review Neovim plugins following modern best practices.
  Use this skill when creating a new Neovim plugin, refactoring an existing one,
  writing busted tests for a plugin, setting up LSP-related plugin behaviour,
  handling conflicts with native Neovim features or other plugins, or ensuring
  compatibility with Neovim 0.11+. Covers directory structure, public API
  design, configuration patterns, keymap lifecycle, health checks, LSP handler
  overrides, and the busted+plenary test setup.
compatibility: Neovim >= 0.11. Requires lua, luacheck (optional). Tests use busted via plenary runner.
license: MIT
---

# Neovim Plugin Development

## Quick reference

- [Architecture & directory layout](references/architecture.md)
- [Neovim API conventions and version notes](references/neovim-api.md)
- [Testing with busted + plenary](references/testing.md)
- [Handling native feature conflicts](references/conflict-resolution.md)
- [Test bootstrap template](assets/minimal_init.lua)

---

## Core rules (always apply)

1. **Target Neovim ≥ 0.11.** Use `vim.version().minor >= 11` in health checks.
   Never use APIs removed before 0.11. See [neovim-api.md](references/neovim-api.md).

2. **Directory structure** — keep infrastructure separate from domain logic:

   ```
   lua/<plugin-name>/
   ├── init.lua          ← public API: setup(), and exported fields
   ├── types.lua         ← LuaCATS @class / @field declarations
   ├── health.lua        ← :checkhealth <plugin-name>
   ├── _core/
   │   ├── configuration.lua   ← settings, deep-merge, reset()
   │   ├── log.lua             ← structured logger, no external deps
   │   └── cmdparse.lua        ← user-command subcommand dispatcher
   └── _lsp/ (or _<domain>/)
       └── *.lua               ← feature modules
   plugin/
   └── <plugin_name>.lua ← lazy-loaded entry point, double-source guard
   spec/
   ├── minimal_init.lua
   └── <plugin-name>/
       └── *_spec.lua
   ```
   See [architecture.md](references/architecture.md) for full details.

3. **Public API in `init.lua`**:
   - `M.setup(opts)` — single options table, deep-merged with defaults.
     Never use positional non-opts arguments like `setup(client, config)`.
   - Auto-attach via `LspAttach` autocmd (for LSP plugins) rather than
     requiring users to call `on_attach` per server.
   - Export `M.on_attach(client, bufnr)` for users who want manual control.

4. **Double-source guard** in `plugin/<plugin_name>.lua`:
   ```lua
   if vim.g.loaded_<plugin_name>_plugin then return end
   vim.g.loaded_<plugin_name>_plugin = true
   ```

5. **LuaCATS annotations** on every public function and `@class`/`@field`
   on every table type. This enables lua-language-server diagnostics.

6. **Health check** via `lua/<plugin-name>/health.lua`:
   ```lua
   local M = {}
   function M.check()
     vim.health.start("<plugin-name>")
     if vim.version().minor < 11 then
       vim.health.error("Neovim 0.11+ required")
     else
       vim.health.ok("Neovim version OK")
     end
   end
   return M
   ```
   Register it in `init.lua`: no additional wiring needed; Neovim discovers
   `health.lua` automatically when the module is `lua/<plugin>/health.lua`.

---

## Configuration module pattern

```lua
-- lua/<plugin>/_core/configuration.lua
local DEFAULTS = {
  ui = { border = "single", silent = true, zindex = 50 },
  keymaps = { close = "<A-s>" },
  display_automatically = true,
  log_level = "warn",
}

local M = {}
local _config = vim.deepcopy(DEFAULTS)
local _initialized = false

--- Deep-merge user opts over defaults.
function M.set(opts)
  _config = vim.tbl_deep_extend("force", _config, opts or {})
end

--- Restore defaults (use in test before_each).
function M.reset()
  _config = vim.deepcopy(DEFAULTS)
  _initialized = false
  vim.g.loaded_<plugin_name> = nil
end

--- Run once per session; call from setup().
function M.initialize_if_needed()
  if _initialized then return end
  _initialized = true
  -- one-time side effects here
end

function M.get() return _config end
return M
```

---

## Keymap lifecycle (buffer-local, popup-scoped)

When installing temporary buffer-local keymaps (e.g. while a popup is open),
save and restore originals reliably:

```lua
-- Save: use nvim_buf_call to find buffer-local originals
local original = vim.api.nvim_buf_call(bufnr, function()
  return vim.fn.maparg(lhs, mode, false, true)
end)

-- Install
vim.keymap.set(mode, lhs, rhs, { buffer = bufnr, silent = true })

-- Restore on cleanup — use vim.keymap.set (NOT vim.fn.mapset)
if original and original.lhs ~= "" then
  vim.keymap.set(original.mode, lhs, original.callback or original.rhs,
    { buffer = bufnr, expr = original.expr == 1,
      silent = original.silent == 1, noremap = original.noremap == 1 })
else
  pcall(vim.keymap.del, mode, lhs, { buffer = bufnr })
end
```

**Never** use `vim.fn.mapset()` for restoration — it crashes in Neovim 0.10+
when `callback` is a Lua function.  
**Never** add `expr = true` to keymaps unless the rhs is actually an
expression string; it breaks `<Plug>` and `<Cmd>` mappings silently.

---

## Logging (no external deps)

```lua
-- lua/<plugin>/_core/log.lua
local M = {}
local LEVELS = { debug=1, info=2, warn=3, error=4 }
local _level = "warn"

function M.set_level(l) _level = l end

local function _log(level, msg, ...)
  if LEVELS[level] < LEVELS[_level] then return end
  local text = select("#", ...) > 0 and string.format(msg, ...) or msg
  vim.notify(("[plugin] " .. text), vim.log.levels[level:upper()])
end

for _, l in ipairs({"debug","info","warn","error"}) do
  M[l] = function(msg, ...) _log(l, msg, ...) end
end
return M
```

---

## Command dispatcher (no external deps)

```lua
-- lua/<plugin>/_core/cmdparse.lua
local M = {}

--- Register a :PluginName subcommand dispatcher.
---@param name string  command name
---@param subcommands table<string, fun(args: table)>
function M.create(name, subcommands, opts)
  opts = opts or {}
  vim.api.nvim_create_user_command(name, function(info)
    local sub = info.fargs[1]
    local fn = subcommands[sub]
    if not fn then
      error(("[%s] unknown subcommand: %s"):format(name, sub))
    end
    fn({ args = vim.list_slice(info.fargs, 2), bang = info.bang })
  end, {
    nargs = "+",
    complete = function(lead)
      return vim.tbl_filter(
        function(k) return k:sub(1, #lead) == lead end,
        vim.tbl_keys(subcommands))
    end,
    desc = opts.desc or name,
  })
end

return M
```

---

## Tests

See [testing.md](references/testing.md) for the full busted+plenary setup.
Key points:
- Use `before_each` with `configuration.reset()` for test isolation.
- Use `nvim_buf_call(bufnr, fn)` to inspect buffer-local keymaps reliably.
- Mock `vim.notify` by replacing it in `before_each` and restoring in `after_each`.
- Use `describe` / `it` blocks; group by module and function name.

---

## Handling conflicts with native Neovim features

See [conflict-resolution.md](references/conflict-resolution.md) for how to:
- Override `vim.lsp.handlers["textDocument/signatureHelp"]` safely.
- Document incompatibilities with noice.nvim, lsp_signature.nvim, nvim-cmp, blink.cmp.
- Provide an `override_native_handler = true/false` option to let users opt out.
