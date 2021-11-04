vim.g.vsnip_filetypes = {
  typescriptreact = {"typescript"},
  javascriptreact = {"javascript"},
  vimwiki = {"markdown"}
}

vim.api.nvim_set_keymap('i', '<Tab>',
  table.concat{
    'vsnip#jumpable(1)',
    '? "<Plug>(vsnip-jump-next)"',
    ': "<Tab>"'
  },
  { silent = true, expr = true}
)

vim.api.nvim_set_keymap('s', '<Tab>',
  table.concat{
    'vsnip#jumpable(1)',
    '? "<Plug>(vsnip-jump-next)"',
    ': "<Tab>"'
  },
  { silent = true, expr = true}
)

vim.api.nvim_set_keymap('i', '<S-Tab>',
  table.concat{
    'vsnip#jumpable(-1)',
    '? "<Plug>(vsnip-jump-prev)"',
    ': "<Tab>"'
  },
  { silent = true, expr = true}
)

vim.api.nvim_set_keymap('s', '<S-Tab>',
  table.concat{
    'vsnip#jumpable(-1)',
    '? "<Plug>(vsnip-jump-prev)"',
    ': "<Tab>"'
  },
  { silent = true, expr = true}
)
