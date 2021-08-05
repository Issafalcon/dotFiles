-- General settings, config and utilities
require("plugins")
require("utils/global-functions")
require("settings")
require("colourscheme")
require("keymappings")

-- LSP
require("lsp")

-- Plugin customizations
require("plugins/comment")
require("plugins/which-key")
require("plugins/compe")
require("plugins/telescope")
require("plugins/treesitter")
require("plugins/gitsigns")
require("plugins/galaxyline")
require("plugins/omnisharp-vim")
require("plugins/colorizer")
require("plugins/git-messenger")
require("plugins/easy-align")
require("plugins/autopairs")
require("plugins/editorconfig")
require("plugins/wiki")
require("plugins/icons")
require("plugins/rnvimr")
require("plugins/testing")
require("plugins/vimspector")
-- dap not quite mature enough yet / I don't know how to configure it to work as well as vimspector
-- require("plugins/dap")
require("plugins/tagbar")
require("plugins/suckless")
