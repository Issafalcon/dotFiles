local actions = require("telescope.actions")

-- Custom function to search vim config files
function _G.search_dev_config()
  require("telescope.builtin").find_files(
    {
      prompt_title = "< Config Files >",
      search_dirs = {"$HOME/dotFiles"},
      hidden = true
    }
  )
end

-- Telescope mappings
vimp.nnoremap("<Leader>ss", ':lua require(\'telescope.builtin\').grep_string({ search = vim.fn.input("Grep For > ")})<CR>')
vimp.nnoremap("<C-p>", ":lua require('telescope.builtin').git_files()<CR>")
vimp.nnoremap("<Leader>sf", ":lua require('telescope.builtin').find_files({hidden = true})<CR>")
vimp.nnoremap("<Leader>sw", ':lua require(\'telescope.builtin\').grep_string { search = vim.fn.expand("<cword>") }<CR>')
vimp.nnoremap("<Leader>sb", ":lua require('telescope.builtin').buffers()<CR>")
vimp.nnoremap("<Leader>sh", ":lua require('telescope.builtin').help_tags()<CR>")
vimp.nnoremap("<Leader>sc", ":lua search_dev_config()<CR>")
vimp.nnoremap("<Leader>sgc", ":lua require('telescope.builtin').git_commits()<CR>")
vimp.nnoremap("<Leader>sgf", ":lua require('telescope.builtin').git_bcommits()<CR>")
vimp.nnoremap("<Leader>sgb", ":lua require('telescope.builtin').git_branches()<CR>")
vimp.nnoremap("<Leader>sgs", ":lua require('telescope.builtin').git_status()<CR>")
vimp.nnoremap("<Leader>st", ":lua require('telescope.builtin').colorscheme()<CR>")
vimp.nnoremap('<A-2>', ":lua require('telescope.builtin').registers()<CR>")
vimp.inoremap('<A-2>', ":lua require('telescope.builtin').registers()<CR>")

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
    },
    prompt_prefix = " ",
    selection_caret = " ",
    file_sorter = require "telescope.sorters".get_fzy_sorter,
    file_ignore_patterns = { "node_modules" },
    file_previewer = require "telescope.previewers".vim_buffer_cat.new,
    grep_previewer = require "telescope.previewers".vim_buffer_vimgrep.new,
    qflist_previewer = require "telescope.previewers".vim_buffer_qflist.new,
    layout_strategy = "vertical",
    layout_config = {
      vertical = {
        height = {
          padding = 4
        }
      }
    },
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
  pickers = {
    grep_string = {
      vimgrep_arguments = {
        "rg",
        "--hidden",
        "--color=never",
        "--no-heading",
        "--with-filename",
        "--line-number",
        "--column",
        "--smart-case",
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
