local DATA_PATH = vim.fn.stdpath("data")
local custom_attach = require("lsp.custom_attach")

local config = {
  cmd = {"ngserver", "--stdio", "--tsProbeLocations", "", "--ngProbeLocations", ""},
  filetypes = {"typescript", "html", "typescriptreact", "typescript.tsx"},
  root_dir = require("lspconfig/util").root_pattern("angular.json", ".git"),
  on_attach = function(client)
    client.resolved_capabilities.document_range_formatting = true
    client.resolved_capabilities.document_formatting = true

    if client.config.flags then
      client.config.flags.allow_incremental_sync = true
    end
    custom_attach(client)
  end,
  settings = {documentFormatting = true}
}

return config
