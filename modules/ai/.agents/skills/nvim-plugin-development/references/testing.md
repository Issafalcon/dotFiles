# Testing with Busted + Plenary

## Setup

Neovim plugins are tested headlessly using [busted](https://lunarmodules.github.io/busted/)
run through [plenary.nvim](https://github.com/nvim-lua/plenary.nvim)'s test runner.
This avoids needing `nlua` or a standalone `busted` binary — plenary is nearly
always available in a Neovim user's environment.

### `spec/minimal_init.lua`

```lua
-- Find plenary on the system (adjust path as needed)
local plenary_path = vim.fn.expand("~/.local/share/nvim/site/pack/core/opt/plenary.nvim")
  -- or: vim.fn.stdpath("data") .. "/lazy/plenary.nvim"

if vim.fn.isdirectory(plenary_path) == 0 then
  error("plenary.nvim not found at: " .. plenary_path)
end

vim.opt.rtp:prepend(plenary_path)
vim.opt.rtp:prepend(vim.fn.getcwd())  -- add plugin root

-- Source plugin entry for command registration (optional)
-- vim.cmd("source plugin/my_plugin.lua")
```

See the [asset template](../assets/minimal_init.lua) for the full version.

### `Makefile`

```makefile
# Portably discover plenary; CI can override with PLENARY_PATH=...
PLENARY_PATH ?= $(shell find ~/.local/share/nvim -name plenary.nvim \
                  -maxdepth 6 -type d 2>/dev/null | head -1)

.PHONY: test
test:
	@nvim --version | head -n 1
	nvim --headless --clean --noplugin -u spec/minimal_init.lua \
	  -c "lua vim.cmd([[PlenaryBustedDirectory spec/my-plugin \
	      { minimal_init = 'spec/minimal_init.lua' }]])"
```

### `.busted`

```lua
return {
  _all = {
    lua = "nlua",   -- used when nlua/busted installed directly
  },
}
```

---

## Test file structure

```lua
-- spec/my-plugin/configuration_spec.lua
local configuration = require("my-plugin._core.configuration")

describe("configuration", function()
  before_each(function()
    configuration.reset()   -- ← always reset before each test
  end)

  describe("defaults", function()
    it("has the expected border", function()
      assert.equals("single", configuration.get().ui.border)
    end)
  end)

  describe("set()", function()
    it("deep-merges nested tables", function()
      configuration.set({ ui = { border = "rounded" } })
      assert.equals("rounded", configuration.get().ui.border)
      -- other ui fields still present:
      assert.is_not_nil(configuration.get().ui.silent)
    end)
  end)

  describe("reset()", function()
    it("restores defaults after set()", function()
      configuration.set({ ui = { border = "double" } })
      configuration.reset()
      assert.equals("single", configuration.get().ui.border)
    end)

    it("clears the loaded guard flag", function()
      configuration.reset()
      assert.is_nil(vim.g.loaded_my_plugin)
    end)
  end)
end)
```

---

## Patterns

### Isolate state with `before_each` + `reset()`

Every spec file that touches plugin state should call `configuration.reset()`
(and any other module-level state reset) in `before_each`:

```lua
before_each(function()
  configuration.reset()
  -- also reset any module-level caches your spec touches
end)
```

### Testing buffer-local keymaps

`nvim_buf_get_keymap` lhs values may differ in case (e.g. `<c-j>` vs `<C-j>`).
Use `nvim_buf_call` + `maparg` for reliable lookup:

```lua
local bufnr = vim.api.nvim_create_buf(false, true)
mappings.setup(bufnr, config.keymaps, win)

local found = vim.api.nvim_buf_call(bufnr, function()
  local info = vim.fn.maparg("<C-j>", "i", false, true)
  return info and info.lhs ~= ""
end)
assert.is_true(found)
```

Or check the module's internal tracking table if it exposes one:
```lua
local state = require("my-plugin._core.mappings")
assert.is_not_nil(state._buf_mappings[bufnr])
```

### Mocking `vim.notify`

```lua
local notifications = {}
local original_notify

before_each(function()
  original_notify = vim.notify
  vim.notify = function(msg, level)
    table.insert(notifications, { msg = msg, level = level })
  end
end)

after_each(function()
  vim.notify = original_notify
  notifications = {}
end)

it("emits a warning", function()
  -- trigger action
  assert.equals(1, #notifications)
  assert.matches("expected text", notifications[1].msg)
end)
```

### Testing pcall-caught errors

For functions that use `error()` (e.g. cmdparse unknown subcommand):

```lua
it("errors on unknown subcommand", function()
  local ok, err = pcall(vim.cmd, "MyPlugin unknown")
  assert.is_false(ok)
  assert.matches("unknown subcommand", err)
end)
```

### Testing user commands

```lua
it("is registered", function()
  local cmds = vim.api.nvim_get_commands({})
  assert.is_not_nil(cmds["MyPlugin"])
end)

it("provides tab-completion", function()
  local completions = vim.fn.getcompletion("MyPlugin ", "cmdline")
  assert.is_true(vim.tbl_contains(completions, "enable"))
end)
```

### Testing domain modules in isolation

Test feature modules directly by constructing inputs rather than triggering
the full plugin flow. This makes tests fast and avoids needing a real LSP
server, filesystem, or external process:

```lua
local renderer = require("my-plugin._<domain>.renderer")

it("produces expected output for known input", function()
  local input = { key = "value", other = 42 }
  local result = renderer.build(input, {})
  assert.is_not_nil(result)
  assert.is_true(vim.tbl_contains(result, "expected line"))
end)
```

Build the minimal data structures your module needs rather than faking a full
external response — it documents your module's contract and catches regressions
without coupling tests to external APIs.

---

## Running tests

```bash
make test           # run all specs
```

For a single file during development:
```bash
nvim --headless --clean --noplugin -u spec/minimal_init.lua \
  -c "lua vim.cmd([[PlenaryBustedFile spec/my-plugin/handler_spec.lua]])"
```

---

## CI (GitHub Actions)

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: rhysd/action-setup-vim@v1
        with:
          neovim: true
          version: v0.11.0
      - name: Install plenary
        run: |
          git clone --depth 1 https://github.com/nvim-lua/plenary.nvim \
            ~/.local/share/nvim/site/pack/core/opt/plenary.nvim
      - run: make test
```
