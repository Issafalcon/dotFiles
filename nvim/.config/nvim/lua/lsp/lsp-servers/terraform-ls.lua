local custom_attach = require("lsp.custom_attach")
local util = require("lspconfig/util")

local config = {
  on_attach = function(client)
    custom_attach(client)
  end,
  cmd = { "terraform-ls", "serve" },
  filetypes = { "terraform", "tf" },
  root_dir = util.root_pattern(".terraform", ".git")
}

return config
