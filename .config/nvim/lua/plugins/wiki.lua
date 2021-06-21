-- VimWiki settings
local wikiList = {}

local defaultWiki = {
  path = '~/Repos/wiki/',
  syntax = 'markdown',
  ext = '.wiki'
}

wikiList[1] = defaultWiki
vim.g.vimwiki_list = wikiList

local ext2SyntaxSettings = {}
local wikiExt = '.wiki'
ext2SyntaxSettings[wikiExt] = 'markdown'

vim.g.vimwiki_ext2syntax = ext2SyntaxSettings

-- VimTex settings
vim.g.vimtex_view_method = 'zathura'
vim.g.vimtex_view_general_viewer = 'zathura'
