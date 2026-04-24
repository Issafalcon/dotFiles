# Neovim API Conventions (0.11+)

## Minimum version: Neovim 0.11

All plugins in this guide target **Neovim ≥ 0.11**. Never use APIs removed
or deprecated before that version.

---

## API do / don't table

| Purpose | ✅ Use (0.11+) | ❌ Avoid |
|---------|--------------|---------|
| Get LSP clients | `vim.lsp.get_clients({bufnr=0})` | `vim.lsp.get_active_clients()` (removed 0.11) |
| Set buffer keymap | `vim.keymap.set(mode, lhs, rhs, {buffer=n})` | `vim.api.nvim_buf_set_keymap()` |
| Delete keymap | `vim.keymap.del(mode, lhs, {buffer=n})` | `vim.api.nvim_buf_del_keymap()` |
| Restore keymap | `vim.keymap.set()` with saved `maparg` data | `vim.fn.mapset()` (crashes with Lua callbacks) |
| Iterate tables | `vim.iter(t):map():filter():totable()` | Manual `for` loops where iter is cleaner |
| Safe calls | `vim.F.npcall(fn, ...)` | Raw `pcall` when you want nil on error |
| Notify user | `vim.notify(msg, vim.log.levels.WARN)` | `print()` or `error()` for user messages |
| Validate args | `vim.validate({ name = {val, "string"} })` | No validation |
| Schedule deferred | `vim.schedule(fn)` | `vim.defer_fn(fn, 0)` for immediate defer |
| Deep copy | `vim.deepcopy(t)` | Manual table copy |
| Deep merge | `vim.tbl_deep_extend("force", base, override)` | Shallow `vim.tbl_extend` for nested tables |
| Float window | `vim.lsp.util.open_floating_preview(lines, ft, opts)` | `nvim_open_win` directly (loses markdown rendering) |
| Open float above cursor | pass `offset_y` / `anchor` to float opts | `floating_window_above_cur_line` shim |

---

## Key API details

### `vim.fn.maparg` and buffer-local keymaps

`vim.fn.maparg(lhs, mode, false, true)` returns an empty dict `{}` for
buffer-local mappings **unless the target buffer is the current buffer**.

Always wrap with `nvim_buf_call`:
```lua
local info = vim.api.nvim_buf_call(bufnr, function()
  return vim.fn.maparg(lhs, mode, false, true)
end)
```

### `vim.keymap.set` options to avoid

- `expr = true` — only use when `rhs` is a Lua function that returns a string
  to be re-evaluated as a keymap. Setting it on `<Cmd>` or `<Plug>` mappings
  silently breaks them in insert mode.
- `replace_keycodes = false` — only needed when `expr = true` and you're
  returning raw terminal keys.

### `vim.lsp.handlers`

Override signature:
```lua
vim.lsp.handlers["textDocument/signatureHelp"] = function(err, result, ctx, config)
  -- err: nil or string
  -- result: lsp.SignatureHelp or nil
  -- ctx: { bufnr, client_id, method, params }
  -- config: handler config table
end
```
Install once in `setup()` when `override_native_handler = true`. Provide
`M.handler` as a public export so power users can install it themselves.

### `vim.lsp.util.open_floating_preview`

```lua
local bufnr, winnr = vim.lsp.util.open_floating_preview(lines, "markdown", {
  border    = config.ui.border,
  focusable = config.ui.focusable,
  focus     = config.ui.focus,
  zindex    = config.ui.zindex,        -- controls stacking above other floats
  max_width = config.ui.max_width,
  max_height= config.ui.max_height,
  close_events = config.ui.close_events,
  -- offset_x / offset_y shift from cursor position
})
```

`zindex` is critical when other plugins (noice, cmp) also open floats.
Default to 50; document how users can adjust it.

### Health checks

```lua
-- Neovim 0.10+ only (use pcall for older):
vim.health.start("section")
vim.health.ok("message")
vim.health.warn("message", {"hint1", "hint2"})
vim.health.error("message")
vim.health.info("message")
```

### `vim.api.nvim_buf_call`

Executes a function in the context of a specific buffer (sets it current
temporarily). Use for any `vim.fn.*` call that is buffer-context-sensitive:
```lua
vim.api.nvim_buf_call(bufnr, function()
  -- vim.fn.* calls here see bufnr as current buffer
end)
```

### Autocommands

```lua
local group = vim.api.nvim_create_augroup("MyPlugin", { clear = true })

vim.api.nvim_create_autocmd("LspAttach", {
  group = group,
  callback = function(ev)
    local client = vim.lsp.get_client_by_id(ev.data.client_id)
    if client and client.server_capabilities.signatureHelpProvider then
      require("my-plugin").on_attach(client, ev.buf)
    end
  end,
})
```

Always pass `group =` to avoid duplicate handlers on `setup()` re-calls.

### User commands

```lua
vim.api.nvim_create_user_command("MyPlugin", function(info)
  -- info.fargs = split args list
  -- info.args  = raw args string
  -- info.bang  = true if ! was used
end, {
  nargs = "+",         -- "0","1","?","+","*"
  complete = function(arglead, cmdline, pos)
    return vim.tbl_keys(subcommands)  -- list of completion candidates
  end,
  desc = "My plugin command",
})
```

---

## Removed / breaking APIs timeline

| API | Removed in | Replacement |
|-----|-----------|-------------|
| `vim.lsp.get_active_clients()` | 0.11 | `vim.lsp.get_clients()` |
| `vim.lsp.buf_get_clients()` | 0.10 | `vim.lsp.get_clients({bufnr=n})` |
| `vim.api.nvim_exec()` | deprecated 0.9 | `vim.cmd()` |
| `vim.loop` (libuv) | alias still works but | `vim.uv` preferred 0.10+ |
| `unpack()` (global) | Lua 5.1 compat | `table.unpack()` |
| `vim.fn.mapset()` with Lua fn | broken 0.10 | `vim.keymap.set()` for restore |

---

## Version guard pattern

```lua
-- In setup() or health.lua
local v = vim.version()
if v.major == 0 and v.minor < 11 then
  vim.notify("my-plugin requires Neovim 0.11+", vim.log.levels.ERROR)
  return
end
```
