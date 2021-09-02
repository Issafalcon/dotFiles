vim.cmd('let g:nvcode_termcolors=256')

vim.g.material_style = 'darker'

require('material').setup({

	contrast = true, -- Enable contrast for sidebars, floating windows and popup menus like Nvim-Tree
	borders = true, -- Enable borders between verticaly split windows

	italics = {
		comments = true, -- Enable italic comments
		keywords = false, -- Enable italic keywords
		functions = false, -- Enable italic functions
		strings = false, -- Enable italic strings
		variables = false -- Enable italic variables
	},

	contrast_windows = { -- Specify which windows get the contrasted (darker) background
		"terminal", -- Darker terminal background
		"packer", -- Darker packer background
		"qf" -- Darker qf list background
	},

	text_contrast = {
		lighter = false, -- Enable higher contrast text for lighter style
		darker = true -- Enable higher contrast text for darker style
	},

	disable = {
		background = false, -- Prevent the theme from setting the background (NeoVim then uses your teminal background)
		term_colors = false, -- Prevent the theme from setting terminal colors
		eob_lines = true -- Hide the end-of-buffer lines
	},

	custom_colors = {
		blue = '#FF9CAC',
	} -- Overwrite highlights with your own
})
vim.cmd('colorscheme material')

vim.api.nvim_set_keymap('n', '<leader>tm', [[<Cmd>lua require('material.functions').toggle_style()<CR>]], { noremap = true, silent = true })
