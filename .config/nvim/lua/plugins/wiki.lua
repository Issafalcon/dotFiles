require('vimp')

-- Call the global function AutoCommit for .wiki. files
vimp.nnoremap('<leader>cg', ':lua AutoCommit()<CR>')

-- vim.cmd "autocmd BufWritePost */wiki/* call AutoCommit"

-- Enable markdown snippets in vimwiki files
vim.cmd('let g:vsnip_filetypes = {}')
vim.cmd('let g:vsnip_filetypes.vimwiki = ["markdown"]')

-- VimWiki settings
local wikiList = {}

local defaultWiki = {
  path = '~/repos/wiki/',
  syntax = 'markdown',
  ext = '.md'
}

wikiList[1] = defaultWiki
vim.g.vimwiki_list = wikiList

vim.g.vimwiki_table_mappings = 0
vim.g.vimwiki_global_ext = 0 -- Prevents vimwiki from treating every .md file as a wiki
vim.g.vimwiki_hl_headers = 1

-- VimTex settings
vim.g.vimtex_view_method = 'zathura'
vim.g.vimtex_view_general_viewer = 'zathura'
