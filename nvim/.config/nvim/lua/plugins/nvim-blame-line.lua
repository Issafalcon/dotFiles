require('vimp')

vim.g.blameLineVirtualTextHighlight = 'Question'

vimp.nnoremap({silent = true }, '<leader>gb', ':ToggleBlameLine<CR>')
