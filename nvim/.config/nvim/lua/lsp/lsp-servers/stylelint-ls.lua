local custom_attach = require("lsp.custom_attach")

local config = {
  on_attach = function(client)
    custom_attach(client)
  end,
  settings = {
    stylelintplus = {
      autoFixOnSave = true
    }
  }
}

return config
