local custom_attach = require("lsp.custom_attach")

require "lspconfig/configs".emmet_ls = {
  default_config = {
    cmd = {"emmet-ls", "--stdio"},
    filetypes = {"html", "css"},
    root_dir = require "lspconfig".util.root_pattern(".git", vim.fn.getcwd())
  }
}

require "lspconfig".emmet_ls.setup {
  on_attach = custom_attach
}
