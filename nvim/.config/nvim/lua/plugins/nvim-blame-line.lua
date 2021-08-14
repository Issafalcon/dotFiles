require('vimp')
require('nvim_utils')

vim.g.blameLineVirtualTextHighlight = 'Question'

vimp.nnoremap({silent = true }, '<leader>gb', ':ToggleBlameLine<CR>')

local autoCmds = {
  blameLine = {
    { "BufEnter", "*", "EnableBlameLine" };
  };
}

nvim_create_augroups(autoCmds)
