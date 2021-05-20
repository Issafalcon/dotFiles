local pid = vim.fn.getpid()
local util = require("lspconfig/util")
local home = vim.fn.expand("$HOME")
local M = {}
local custom_attach = require("lsp.custom_attach")

function M.getConfig()
  local config = {

  cmd = { home .. "/.cache/omnisharp-vim/omnisharp-roslyn/run", "--languageserver", "--hostPID", tostring(pid)},
  on_attach = function(client)
    custom_attach(client)
  end,
  filetypes = {"cs", "vb"},
  root_dir = util.root_pattern("*.csproj", "*.sln")
  }

  return config
end

return M
