local sign_define = vim.fn.sign_define

-- Works with the autoformat on save autocommand that is in custom_attach function
-- This will 1. create an autocommand for every buffer to format on save. And then save again after
-- formatting is done (only if there are no changes to the buffer)
vim.lsp.handlers["textDocument/formatting"] = function(err, _, result, _, bufnr)
  if err ~= nil or result == nil then
    return
  end
  if not vim.api.nvim_buf_get_option(bufnr, "modified") then
    local view = vim.fn.winsaveview()
    vim.lsp.util.apply_text_edits(result, bufnr)
    vim.fn.winrestview(view)
    if bufnr == vim.api.nvim_get_current_buf() then
      vim.api.nvim_command("noautocmd :update")
    end
  end
end

-- Diagnostics handlers
vim.lsp.handlers["textDocument/publishDiagnostics"] =
  vim.lsp.with(
  vim.lsp.diagnostic.on_publish_diagnostics,
  {
    underline = true,
    signs = true,
    update_in_insert = true,
    virtual_text = false
  }
)

-- Additional diagnostic configurations
-- vim.cmd [[autocmd CursorHold,CursorHoldI * lua vim.lsp.diagnostic.show_line_diagnostics({focusable=false})]]

sign_define(
  "LspDiagnosticsSignError",
  {texthl = "LspDiagnosticsSignError", text = "", numhl = "LspDiagnosticsSignError"}
)
sign_define(
  "LspDiagnosticsSignWarning",
  {texthl = "LspDiagnosticsSignWarning", text = "", numhl = "LspDiagnosticsSignWarning"}
)
sign_define("LspDiagnosticsSignHint", {texthl = "LspDiagnosticsSignHint", text = "", numhl = "LspDiagnosticsSignHint"})
sign_define(
  "LspDiagnosticsSignInformation",
  {texthl = "LspDiagnosticsSignInformation", text = "", numhl = "LspDiagnosticsSignInformation"}
)
