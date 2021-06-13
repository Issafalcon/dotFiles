local wikiLists = {};

local markdownWiki =   {
  path = '~/vimwiki/',
  syntax = 'markdown',
  ext = '.md'
};

local default =   {
  path = '~/vimwiki/',
  syntax = 'default',
  ext = '.wiki'
};

wikiLists[1] = default;
wikiLists[2] = markdownWiki;

vim.g.vimwiki_list = wikiLists;
