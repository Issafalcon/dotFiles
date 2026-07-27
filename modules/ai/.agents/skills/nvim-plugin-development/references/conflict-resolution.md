# Handling Conflicts with Other Plugins and Native Features

## The problem

Neovim plugins frequently compete for the same resources: autocmd events,
buffer-local keymaps, floating windows, and LSP handlers. Identifying and
handling these conflicts gracefully — while giving users escape hatches —
distinguishes a polished plugin from a fragile one.

---

## 1. Autocmd group conflicts

Always create autocmd groups with a namespaced name and `clear = true` to
prevent duplicate registrations on `setup()` re-calls:

```lua
local group = vim.api.nvim_create_augroup("MyPlugin", { clear = true })
vim.api.nvim_create_autocmd("BufWritePre", {
  group = group,
  callback = function(ev) ... end,
})
```

If users call `setup()` twice (common during config reload), `clear = true`
ensures the old autocmds are removed before re-registering.

---

## 2. Keymap conflicts

When installing keymaps, save and restore originals so you don't permanently
destroy the user's existing bindings:

```lua
-- Save original before installing
local original = vim.api.nvim_buf_call(bufnr, function()
  return vim.fn.maparg(lhs, mode, false, true)
end)

vim.keymap.set(mode, lhs, rhs, { buffer = bufnr, silent = true })

-- Restore on cleanup
if original and original.lhs ~= "" then
  vim.keymap.set(original.mode, lhs, original.callback or original.rhs,
    { buffer = bufnr, expr = original.expr == 1,
      silent = original.silent == 1, noremap = original.noremap == 1 })
else
  pcall(vim.keymap.del, mode, lhs, { buffer = bufnr })
end
```

Give users the option to disable your default keymaps entirely:

```lua
-- In defaults:
{ keymaps = { enable = true, close = "<Esc>" } }

-- In setup():
if config.keymaps.enable then
  install_keymaps(bufnr)
end
```

---

## 3. Floating window stacking (zindex)

When multiple plugins open floating windows simultaneously (completion,
diagnostics, your plugin), `zindex` controls which float appears on top.

```lua
local bufnr, winnr = vim.lsp.util.open_floating_preview(lines, "markdown", {
  border  = config.ui.border,
  zindex  = config.ui.zindex,   -- default 50; let users tune this
})
```

Document the default and mention that users may need to adjust if floats from
other plugins appear above or below yours unexpectedly.

Health check hint:
```lua
vim.health.info("Float zindex: " .. tostring(cfg.ui.zindex) ..
  " (increase if other floats appear on top)")
```

---

## 4. LSP handler conflicts

If your plugin overrides a global LSP handler (e.g.
`vim.lsp.handlers["textDocument/signatureHelp"]`), multiple plugins doing
the same will silently overwrite each other. Always:

1. Give users an explicit opt-out:

```lua
-- In defaults:
{ override_native_handler = true }

-- In initialize_if_needed():
if cfg.override_native_handler then
  vim.lsp.handlers["the/method"] = require("my-plugin").handler
end
```

2. Expose `M.handler` as a public field so power users can install it
   themselves rather than relying on the global override.

3. Document any plugins known to set the same handler in your README.

---

## 5. Native feature conflicts

When your plugin duplicates a native Neovim feature, users may end up with
duplicate UI. Provide a config option and guide users in your README:

```lua
-- Example: plugin provides its own diagnostic float
{ replace_native_feature = false }
```

Health check:
```lua
if cfg.replace_native_feature then
  -- check if the native feature is still active and warn
  vim.health.warn(
    "Native feature X is still enabled alongside my-plugin",
    { "Consider disabling it to avoid duplicate UI" }
  )
end
```

---

## 6. Documenting known conflicts in README

Add a conflict table so users can resolve issues without opening bug reports:

```markdown
## Known conflicts

| Plugin / Feature   | Symptom               | Fix                                        |
|--------------------|-----------------------|--------------------------------------------|
| some-plugin        | Duplicate popups      | Disable the overlapping feature in that plugin |
| another-plugin     | Keymap overwritten    | Set `keymaps.enable = false` and map manually  |
| Native feature X   | Duplicate UI          | See `replace_native_feature` config option |
```

---

## Neovim 0.11+ `vim.lsp.config` (new server config API)

With the new `vim.lsp.config` API, per-server `on_attach` is less common.
An `LspAttach` autocmd approach handles this correctly without any user
wiring, since the autocmd fires regardless of how the LSP client was started:

```lua
vim.api.nvim_create_autocmd("LspAttach", {
  group = vim.api.nvim_create_augroup("MyPlugin", { clear = true }),
  callback = function(ev)
    local client = vim.lsp.get_client_by_id(ev.data.client_id)
    -- filter by server capability if your plugin is LSP-specific:
    -- if client and client.server_capabilities.someCapability then
    require("my-plugin").on_attach(client, ev.buf)
  end,
})
```
