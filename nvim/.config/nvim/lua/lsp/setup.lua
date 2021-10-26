local custom_attach = require("lsp.custom_attach")
local lsp_installer = require("nvim-lsp-installer")
local efmConfig = require("lsp.lsp-servers.efm-ls")
local stylelintConfig = require("lsp.lsp-servers.stylelint-ls")
local texlabConfig = require("lsp.lsp-servers.texlab-ls")
local terraformConfig = require("lsp.lsp-servers.terraform-ls")

-- symbols for autocomplete
vim.lsp.protocol.CompletionItemKind = {
  "   (Text) ",
  "   (Method)",
  "   (Function)",
  "   (Constructor)",
  "   (Field)",
  "[] (Variable)",
  "   (Class)",
  "   (Interface)",
  "   (Module)",
  " 襁 (Property)",
  "   (Unit)",
  "   (Value)",
  " 練 (Enum)",
  "   (Keyword)",
  "   (Snippet)",
  "   (Color)",
  "   (File)",
  "   (Reference)",
  "   (Folder)",
  "   (EnumMember)",
  "   (Constant)",
  "   (Struct)",
  "   (Event)",
  "   (Operator)",
  "   (TypeParameter)"
}

local function addAutoFormatOnSave(client)
  -- Format on save if the lsp has document formattin enabled
  if client.resolved_capabilities.document_formatting then
    vim.api.nvim_command [[augroup Format]]
    vim.api.nvim_command [[autocmd! * <buffer>]]
    vim.api.nvim_command [[autocmd BufWritePost <buffer> lua vim.lsp.buf.formatting()]]
    vim.api.nvim_command [[augroup END]]
  end
end

-- vim.lsp.set_log_level("debug")

-- Set the capabilities for all servers. Snippet support can be enabled for all
local clientCapabilities = vim.lsp.protocol.make_client_capabilities()
clientCapabilities.textDocument.completion.completionItem.snippetSupport = true

-- Get the servers installed via nvim-lsp-installer plugin
-- see https://github.com/williamboman/nvim-lsp-installer/tree/main/lua/nvim-lsp-installer/servers
local installed_servers = lsp_installer.get_installed_servers()

-- Attach and customize servers installed from nvim-lsp-intaller plugin
for _, server in pairs(installed_servers) do
  local opts = {}

  -- Add custom options here.
  if server.name == "eslintls" then
    opts.on_attach = function(client)
      client.resolved_capabilities.document_formatting = true
      custom_attach(client)
      -- addAutoFormatOnSave(client) --Add autoformat on save: This is annoying when multiple LSPs are attached
    end
    opts.settings = {
      format = {
        enable = true
      }
    }
  end

  if server.name == "sumneko_lua" then
    opts.settings = {Lua = {diagnostics = {globals = {"vim", "vimp", "nvim"}}}}
  end

  if server.name == "bashls" then
   opts.filetypes = {"sh", "zsh"}
  end

  if server.name == "tsserver" then
    opts.on_attach = function(client)
      custom_attach(client)
      client.resolved_capabilities.formatting = true
      client.resolved_capabilities.document_formatting = true
    end
  end

  if server.name == "jsonls" then
    opts.on_attach = function(client)
      custom_attach(client)
    end
    opts.filetypes = {"json", "jsonc"}
    opts.settings = {
      json = {
        -- Schemas https://www.schemastore.org
        schemas = {
          {
            fileMatch = {"package.json"},
            url = "https://json.schemastore.org/package.json"
          },
          {
            fileMatch = {"tsconfig*.json"},
            url = "https://json.schemastore.org/tsconfig.json"
          },
          {
            fileMatch = {
              ".prettierrc",
              ".prettierrc.json",
              "prettier.config.json"
            },
            url = "https://json.schemastore.org/prettierrc.json"
          },
          {
            fileMatch = {".eslintrc", ".eslintrc.json"},
            url = "https://json.schemastore.org/eslintrc.json"
          },
          {
            fileMatch = {".babelrc", ".babelrc.json", "babel.config.json"},
            url = "https://json.schemastore.org/babelrc.json"
          },
          {
            fileMatch = {"lerna.json"},
            url = "https://json.schemastore.org/lerna.json"
          },
          {
            fileMatch = {"now.json", "vercel.json"},
            url = "https://json.schemastore.org/now.json"
          },
          {
            fileMatch = {
              ".stylelintrc",
              ".stylelintrc.json",
              "stylelint.config.json"
            },
            url = "http://json.schemastore.org/stylelintrc.json"
          }
        }
      }
    }
  end

  if opts.on_attach == nil then
    opts.on_attach = custom_attach
  end

  opts.capabilities = clientCapabilities

  server:setup(opts)
end

-- Add lsp servers from local setup (i.e. Not installed as per nvim-lsp-installer
local servers = {}

-- Supplement omnisharp-vim with diagnostic info
servers.omnisharp = require("lsp.lsp-servers.omnisharp-ls").getConfig()

servers.efm = efmConfig

servers.stylelint_lsp = stylelintConfig

servers.texlab = texlabConfig

servers.terraformls = terraformConfig

for server, config in pairs(servers) do
  require("lspconfig")[server].setup(vim.tbl_deep_extend("force", {capabilities = clientCapabilities}, config))
end
