require('vimp')

-- Call the global function AutoCommit for .wiki. files
vimp.nnoremap('<leader>cg', ':lua AutoCommit()<CR>')

-- vim.cmd "autocmd BufWritePost */wiki/* call AutoCommit"

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

local latexmkOptions = {
  '-verbose', '-file-line-error', '-synctex=1', '-interaction=nonstopmode', '-shell-escape'
}

vim.g.vimtex_compiler_latexmk = {
  build_dir = '',
  callback = 1,
  continuous = 1,
  executable = 'latexmk',
  hooks = {},
  options = latexmkOptions
}

vimp.nnoremap({"silent"}, "\\lg", [[:Start latexmk-custom-launch.sh %:p<CR>]])

-- vim.cmd([[
-- let g:vimtex_compiler_latexmk = {
--               \ 'build_dir' : '',
--               \ 'callback' : 1,
--               \ 'continuous' : 1,
--               \ 'executable' : 'latexmk',
--               \ 'hooks' : [],
--               \ 'options' : [
--               \   '-verbose',
--               \   '-file-line-error',
--               \   '-synctex=1',
--               \   '-interaction=nonstopmode',
--               \   '-shell-escape'
--               \ ],
--               \}
--   ]]
-- )
