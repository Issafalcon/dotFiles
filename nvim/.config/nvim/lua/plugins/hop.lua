require("hop").setup()
require("vimp")

-- place this in one of your configuration file(s)
vimp.nnoremap('<leader><leader>w', "<cmd>lua require'hop'.hint_words()<cr>")
vimp.nnoremap('<leader><leader>o', "<cmd>lua require'hop'.hint_char1()<cr>")
vimp.nnoremap('<leader><leader>t', "<cmd>lua require'hop'.hint_char2()<cr>")
vimp.nnoremap('<leader><leader>/', "<cmd>lua require'hop'.hint_patterns()<cr>")
vimp.nnoremap('<leader><leader>l', "<cmd>lua require'hop'.hint_lines()<cr>")
