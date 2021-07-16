-- Make Ranger replace netrw and be the file explorer
vim.g.rnvimr_ex_enable = 1
vim.g.rnvimr_draw_border = 1
vim.g.rnvimr_enable_picker = 0
vim.g.rnvimr_bw_enable = 1
-- Change the border's color
vim.g.rnvimr_border_attr = {fg=14, bg=-1}

-- Draw border with both
vim.g.rnvimr_ranger_cmd = 'ranger --cmd="set draw_borders both"'
vim.api.nvim_set_keymap('n', '-', ':RnvimrToggle<CR>', {noremap = true, silent = true})
