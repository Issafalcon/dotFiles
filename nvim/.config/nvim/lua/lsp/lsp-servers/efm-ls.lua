local util = require("lspconfig/util")
local custom_attach = require("lsp.custom_attach")

local lua_fmt = {
  formatCommand = "luafmt --indent-count 2 --line-width 120 --stdin",
  formatStdin = true
}

local prettier = {
  formatCommand = "prettier --stdin --stdin-filepath ${INPUT}",
  formatStdin = true
}

local eslint_d = {
  lintCommand = "eslint_d -f unix --stdin --stdin-filename ${INPUT}",
  lintStdin = true,
  lintFormats = {"%f:%l:%c: %m"},
  lintIgnoreExitCode = true,
  formatCommand = "eslint_d --fix-to-stdout --stdin --stdin-filename=${INPUT}",
  formatStdin = true
}

-- shell

local shellcheck = {
  lintCommand = "shellcheck -f gcc -x -",
  lintStdin = true,
  lintFormats = {
    "%f:%l:%c: %trror: %m",
    "%f:%l:%c: %tarning: %m",
    "%f:%l:%c: %tote: %m"
  }
}

local shfmt = {
  formatCommand = "shfmt -ci -s -bn",
  formatStdin = true
}

-- vlang

local vfmt = {
  formatCommand = "v fmt"
}

-- data

local markdownlint = {
  lintCommand = "markdownlint -s -c ~/.config/efm-langserver/markdown.yaml",
  lintStdin = true,
  lintFormats = {"%f:%l %m", "%f:%l:%c %m", "%f: %l: %m"}
}

local yamllint = {
  lintCommand = "yamllint -c ~/.config/efm-langserver/yamllint.yaml -f parsable -",
  lintStdin = true,
  lintFormats = {"%f:%l %m", "%f:%l:%c %m", "%f: %l: %m"}
}

local efmConfig = {
  cmd = {
    "/home/linuxbrew/.linuxbrew/Cellar/efm-langserver/0.0.36/bin/efm-langserver"
  },
  root_dir = util.root_pattern("package.json", "tsconfig.json", "jsconfig.json", ".git"),
  on_attach = function(client)
    custom_attach(client)
  end,
  filetypes = {
    "javascript",
    "typescript",
    "typescriptreact",
    "javascriptreact",
    -- "css",
    -- "scss",
    "json",
    "html",
    "lua",
    "sh",
    "v",
    "markdown",
    "yaml"
  },
  init_options = {
    documentFormatting = true,
    codeAction = false,
    hover = false,
    completion = false
  },
  settings = {
    languages = {
      typescript = {eslint_d, prettier},
      javascript = {eslint_d, prettier},
      typescriptreact = {eslint_d, prettier},
      javascriptreact = {eslint_d, prettier},
      -- css = {prettier},
      -- scss = {prettier},
      json = {prettier},
      html = {prettier},
      lua = {lua_fmt},
      sh = {shfmt, shellcheck},
      v = {vfmt},
      markdown = {markdownlint},
      yaml = {yamllint, prettier}
    }
  }
}

return efmConfig
