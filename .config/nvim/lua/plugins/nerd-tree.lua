require('vimp')

vimp.nnoremap('<C-n>', ":NERDTreeToggle<CR>")
vim.g.NERDTreeGitStatusUseNerdFonts = 1
vim.g.NERDTreeGitStatusIndicatorMapCustom = {
  Modified  = '✹',
  Staged    = '✚',
  Untracked = '✭',
  Renamed   = '➜',
  Unmerged  = '═',
  Deleted   = '✖',
  Dirty     = '✗',
  Ignored   = '☒',
  Clean     = '✔︎',
  Unknown   = '?',
}

-- Workaround for a bug in nerd-tree devicons where tex files move across the screen
vim.cmd('let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols = {}')
vim.cmd('let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols["tex"] = "ƛ"')

