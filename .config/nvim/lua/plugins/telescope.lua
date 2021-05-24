local actions = require("telescope.actions")
local utils = require("utils")

-- Custom function to search vim config files
function _G.search_dev_config()
  require("telescope.builtin").find_files(
    {
      prompt_title = "< Nvim Config Files >",
      cwd = "$HOME/.config"
    }
  )
end

-- Telescope mappings
utils.map(
  "n",
  "<Leader>ss",
  ':lua require(\'telescope.builtin\').grep_string({ search = vim.fn.input("Grep For > ")})<CR>',
  {noremap = true}
)
utils.map("n", "<C-p>", ":lua require('telescope.builtin').git_files()<CR>", {noremap = true})
utils.map("n", "<Leader>sf", ":lua require('telescope.builtin').find_files()<CR>", {noremap = true})
utils.map("n", "<Leader>sw", ':lua require(\'telescope.builtin\').grep_string { search = vim.fn.expand("<cword>") }<CR>', {noremap = true})
utils.map("n", "<Leader>sb", ":lua require('telescope.builtin').buffers()<CR>", {noremap = true})
utils.map("n", "<Leader>sh", ":lua require('telescope.builtin').help_tags()<CR>", {noremap = true})
utils.map("n", "<Leader>sc", ":lua search_dev_config()<CR>", {noremap = true})
utils.map("n", "<Leader>sgc", ":lua require('telescope.builtin').git_commits()<CR>", {noremap = true})
utils.map("n", "<Leader>sgb", ":lua require('telescope.builtin').git_branches()<CR>", {noremap = true})
utils.map("n", "<Leader>sgs", ":lua require('telescope.builtin').git_status()<CR>", {noremap = true})
utils.map("n", "<Leader>scs", ":lua require('telescope.builtin').colorscheme()<CR>", {noremap = true})

require("telescope").setup {
  defaults = {
    vimgrep_arguments = {
      "rg",
      "--color=never",
      "--no-heading",
      "--with-filename",
      "--line-number",
      "--column",
      "--smart-case",
      "-u",
      "-u"
    },
    prompt_position = "top",
    prompt_prefix = " ",
    selection_caret = " ",
    file_sorter = require "telescope.sorters".get_fzy_sorter,
    file_ignore_patterns = {},
    file_previewer = require "telescope.previewers".vim_buffer_cat.new,
    grep_previewer = require "telescope.previewers".vim_buffer_vimgrep.new,
    qflist_previewer = require "telescope.previewers".vim_buffer_qflist.new,
    mappings = {
      i = {
        ["<C-j>"] = actions.move_selection_next,
        ["<C-k>"] = actions.move_selection_previous,
        ["<C-q>"] = actions.smart_send_to_qflist + actions.open_qflist,
        -- To disable a keymap, put [map] = false
        -- So, to not map "<C-n>", just put
        -- ["<c-x>"] = false,
        ["<esc>"] = actions.close,
        -- Add up multiple actions
        ["<CR>"] = actions.select_default + actions.center
      },
      n = {
        ["<C-j>"] = actions.move_selection_next,
        ["<C-k>"] = actions.move_selection_previous,
        ["<C-q>"] = actions.smart_send_to_qflist + actions.open_qflist
      }
    }
  },
  extensions = {
    fzy_native = {
      override_generic_sorter = false,
      override_file_sorter = true
    }
  }
}
require("telescope").load_extension("fzy_native")
