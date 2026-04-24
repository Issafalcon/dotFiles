# Handling Conflicts with Native Neovim Features

## The problem

Neovim 0.10+ and many popular plugins all hook into `textDocument/signatureHelp`
at the same time. When multiple handlers are active, users see duplicate or
stacked floating windows. Similarly, native keymaps may conflict with plugin
keymaps.

---

## Providing an opt-out config option

Always give users control with an explicit option:

```lua
-- In configuration.lua defaults:
{
  override_native_handler = true,  -- default: we own the handler
}
```

In `setup()`:
```lua
function M.setup(opts)
  configuration.set(opts)
  configuration.initialize_if_needed()
end

-- In initialize_if_needed():
function M.initialize_if_needed()
  if _initialized then return end
  _initialized = true
  local cfg = configuration.get()
  if cfg.override_native_handler then
    vim.lsp.handlers["textDocument/signatureHelp"] = require("my-plugin").handler
  end
  -- register LspAttach autocmd, etc.
end
```

---

## noice.nvim

`noice.nvim` by default intercepts `textDocument/signatureHelp`. Disable
noice's signature handler so your plugin can own it:

```lua
-- In user's config (document this):
require("noice").setup({
  lsp = {
    signature = { enabled = false },
  },
})
require("my-plugin").setup()
```

**Health check hint** — detect noice and warn if its signature is enabled:
```lua
local ok, noice = pcall(require, "noice")
if ok then
  local noice_cfg = noice.config and noice.config.options
  if noice_cfg and noice_cfg.lsp and noice_cfg.lsp.signature
      and noice_cfg.lsp.signature.enabled ~= false then
    vim.health.warn(
      "noice.nvim signature help is enabled and may conflict",
      { "Set lsp.signature.enabled = false in noice setup" }
    )
  end
end
```

---

## lsp_signature.nvim

`lsp_signature.nvim` installs its own global `vim.lsp.handlers["textDocument/signatureHelp"]`.
The two plugins cannot coexist — document this clearly in README:

```
Do not load both lsp_signature.nvim and <your-plugin> simultaneously.
Remove or disable lsp_signature.nvim from your plugin manager.
```

Health check:
```lua
if package.loaded["lsp_signature"] then
  vim.health.warn(
    "lsp_signature.nvim is loaded and will conflict",
    { "Remove lsp_signature.nvim from your plugin list" }
  )
end
```

---

## nvim-cmp

The `cmp-nvim-lsp-signature-help` source adds its own signature popup.
Document how to disable it:

```lua
-- In user's config:
require("cmp").setup({
  sources = require("cmp").config.sources({
    { name = "nvim_lsp" },
    -- Remove: { name = "nvim_lsp_signature_help" }
  }),
})
```

---

## blink.cmp

blink.cmp has a built-in signature popup. Disable it:

```lua
require("blink.cmp").setup({
  signature = { enabled = false },
})
```

---

## Neovim 0.11+ `vim.lsp.config` (new server config API)

With the new `vim.lsp.config` API, per-server `on_attach` is less common.
Your `LspAttach` autocmd approach handles this correctly without any user
wiring, since the autocmd fires regardless of how the LSP client was started.

```lua
-- This works whether the user uses lspconfig, vim.lsp.config, or manual start:
vim.api.nvim_create_autocmd("LspAttach", {
  group = vim.api.nvim_create_augroup("MyPlugin", { clear = true }),
  callback = function(ev)
    local client = vim.lsp.get_client_by_id(ev.data.client_id)
    if client and client.server_capabilities.signatureHelpProvider then
      require("my-plugin").on_attach(client, ev.buf)
    end
  end,
})
```

---

## Native `vim.lsp.buf.signature_help()`

When your plugin owns the `textDocument/signatureHelp` handler, calls to
`vim.lsp.buf.signature_help()` from other plugins or user mappings will
route through your handler — this is expected and correct behaviour.

If a user wants to access the native popup (e.g. to enter the float in
normal mode for scrolling), they can set `ui.focusable = true` and
`ui.focus = false` in your config, then call `vim.lsp.buf.signature_help()`
while the popup is visible — it will re-focus the existing float.

---

## Documenting conflicts in README

Add a section like:

```markdown
## Preventing Conflicting Signature Popups

| Plugin         | Fix |
|----------------|-----|
| noice.nvim     | Set `lsp.signature.enabled = false` in noice setup |
| lsp_signature.nvim | Remove from plugin list — cannot coexist |
| nvim-cmp       | Remove `nvim_lsp_signature_help` source |
| blink.cmp      | Set `signature.enabled = false` |
| Native Neovim  | Handled automatically when `override_native_handler = true` (default) |
```
