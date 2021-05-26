require('nvim_utils')
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

-- If another buffer tries to replace NERDTree, put it in the other window, and bring back NERDTree.
local autocmds = {
  nerdtree = {
    {'BufEnter', '*', [[if bufname('#') =~ 'NERD_tree_\d\+' && bufname('%') !~ 'NERD_tree_\d\+' && winnr('$') > 1 | let buf=bufnr() | buffer# | execute "normal! \<C-W>w" | execute 'buffer'.buf | endif]]}
  }
}

nvim_create_augroups(autocmds)
