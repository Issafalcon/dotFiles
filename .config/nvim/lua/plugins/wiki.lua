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
