vim.cmd('let g:nvcode_termcolors=256')

vim.g.material_style = 'darker'
vim.g.material_italic_comments = true
vim.g.material_italic_keywords = false
vim.g.material_italic_functions = false
vim.g.material_italic_variables = false
vim.g.material_contrast = true
vim.g.material_borders = true
vim.g.material_hide_eob = true
vim.g.material_disable_background = false
vim.g.material_custom_colors = {
  blue = "#FF9CAC"
}
vim.cmd('colorscheme material')

vim.api.nvim_set_keymap('n', '<leader>tm', [[<Cmd>lua require('material.functions').toggle_style()<CR>]], { noremap = true, silent = true })
