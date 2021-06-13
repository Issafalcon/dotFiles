local wikiLists = {};

local markdownWiki =   {
  path = '~/vimwiki/',
  syntax = 'markdown',
  ext = '.md'
};

wikiLists[1] = markdownWiki;

vim.g.vimwiki_list = wikiLists;
