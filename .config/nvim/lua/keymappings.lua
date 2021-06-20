local utils = require('utils')
require('vimp')

-- Mapping for paste (Paste over visually selected text with whatever is in unnamed register)
utils.map('v', '<Leader>p', '"_dP', { noremap = true })

vimp.inoremap('<C-d>', '<Esc>')

-- Copying to shared clipboard slots in register
utils.map('v', '<Leader>y', '"+y', { noremap = true })
utils.map('n', '<Leader>y', '"+y', { noremap = true })
utils.map('n', '<Leader>Y', 'gg"+yG', { noremap = true })

-- Delete without adding to register
utils.map('n', '<Leader>d', '"_d', { noremap = true })
utils.map('v', '<Leader>d', '"_d', { noremap = true })

-- Quick move lines up and down in buffer
utils.map('n', '<A-j>', ':m .+1<CR>==', { noremap = true })
utils.map('n', '<A-k>', ':m .-2<CR>==', { noremap = true })
utils.map('v', '<A-k>', ':m \'>+1<CR>gv=gv', { noremap = true })
utils.map('v', '<A-j>', ':m \'<-2<CR>gv=gv', { noremap = true })
utils.map('i', '<A-j>', '<Esc>:m .+1<CR>==gi', { noremap = true })
utils.map('i', '<A-k>', '<Esc>:m .-2<CR>==gi', { noremap = true })

-- Tab switch buffer
vim.api.nvim_set_keymap('n', '<TAB>', ':bnext<CR>', {noremap = true, silent = true})
vim.api.nvim_set_keymap('n', '<S-TAB>', ':bprevious<CR>', {noremap = true, silent = true})

-- Toggle line number / relative line number
vimp.nnoremap('<leader>n', function()
  vim.wo.number = not vim.wo.number
end)

vimp.nnoremap('<leader>rn', function()
  vim.wo.relativenumber = not vim.wo.relativenumber
end)

-- Undotreee
vimp.nnoremap('<A-t>', ':UndotreeToggle<cr>')

-- QF List toggle and navigation
vim.g.issafalcon_qf_g = 0
vim.g.issafalcon_qf_l = 0

function _G.ToggleQFList(global)
  if global == 1 then
    if vim.g.issafalcon_qf_g == 1 then
      vim.g.issafalcon_qf_g = 0
      vim.cmd('cclose')
    else
      vim.g.issafalcon_qf_g = 1
      vim.cmd('copen')
    end
  else
    if vim.g.issafalcon_qf_l == 1 then
      vim.g.issafalcon_qf_l = 0
      vim.cmd('lclose')
    else
      vim.g.issafalcon_qf_l = 1
      vim.cmd('lopen')
    end
  end
end

vimp.nnoremap('<C-k>', ':cnext<CR>zz')
vimp.nnoremap('<C-j>', ':cprev<CR>zz')
vimp.nnoremap('<leader>k', ':lnext<CR>zz')
vimp.nnoremap('<leader>j', ':lprev<CR>zz')
vimp.nnoremap('<C-q>', ':lua ToggleQFList(1)<CR>')
vimp.nnoremap('<leader>q', ':lua ToggleQFList(0)<CR>')
