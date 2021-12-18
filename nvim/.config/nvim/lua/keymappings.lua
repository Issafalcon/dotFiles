require('vimp')

-- Exit insert mode without using Esc
vim.api.nvim_set_keymap('i', 'jk', '<Esc>', {noremap = true, silent = true})

-- Mapping for paste (Paste over visually selected text with whatever is in unnamed register)
vimp.vnoremap('<Leader>p', '"_dP')

-- Copying to shared clipboard slots in register
vimp.vnoremap('<Leader>y', '"+y')
vimp.nnoremap('<Leader>y', '"+y')
vimp.nnoremap('<Leader>Y', 'gg"+yG')

-- Delete without adding to register
vimp.nnoremap('<A-d>', '"_d')
vimp.vnoremap('<A-d>', '"_d')

-- Quick move lines up and down in buffer
vimp.nnoremap('<A-j>', ':m .+1<CR>==')
vimp.nnoremap('<A-k>', ':m .-2<CR>==')
vimp.vnoremap('<A-k>', ':m \'>+1<CR>gv=gv')
vimp.vnoremap('<A-j>', ':m \'<-2<CR>gv=gv')
vimp.inoremap('<A-j>', '<Esc>:m .+1<CR>==gi')
vimp.inoremap('<A-k>', '<Esc>:m .-2<CR>==gi')

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
