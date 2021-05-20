local wk = require("which-key")
local M = {}

wk.setup {
    plugins = {
        marks = true,                                                           -- shows a list of your marks on ' and `
        registers = true,                                                       -- shows your registers on " in NORMAL or <C-r> in INSERT mode
        -- the presets plugin, adds help for a bunch of default keybindings in Neovim
        -- No actual key bindings are created
        presets = {
            operators = true,                                                   -- adds help for operators like d, y, ...
            motions = true,                                                     -- adds help for motions
            text_objects = true,                                                -- help for text objects triggered after entering an operator
            windows = true,                                                     -- default bindings on <c-w>
            nav = true,                                                         -- misc bindings to work with windows
            z = true,                                                           -- bindings for folds, spelling and others prefixed with z
            g = true                                                            -- bindings for prefixed with g
        }
    },
    icons = {
        breadcrumb = "»",                                                       -- symbol used in the command line area that shows your active key combo
        separator = "➜",                                                        -- symbol used between a key and it's label
        group = "+"                                                             -- symbol prepended to a group
    },
    window = {
        border = "single",                                                      -- none, single, double, shadow
        position = "bottom",                                                    -- bottom, top
        margin = {1, 0, 1, 0},                                                  -- extra window margin [top, right, bottom, left]
        padding = {2, 2, 2, 2}                                                  -- extra window padding [top, right, bottom, left]
    },
    layout = {
        height = {min = 4, max = 25},                                           -- min and max height of the columns
        width = {min = 20, max = 50},                                           -- min and max width of the columns
        spacing = 3                                                             -- spacing between columns
    },
    hidden = {"<silent>", "<cmd>", "<Cmd>", "<CR>", "call", "lua", "^:", "^ "}, -- hide mapping boilerplate
    show_help = true                                                            -- show help message on the command line when the popup is visible
}

M.defaultOpts = {
    mode = "n",     -- NORMAL mode
    prefix = "",
    buffer = nil,   -- Global mappings. Specify a buffer number for buffer local mappings
    silent = true,  -- use `silent` when creating keymaps
    noremap = true, -- use `noremap` when creating keymaps
    nowait = false  -- use `nowait` when creating keymaps
}

wk.register({
        ["<C-r>"] = "which_key_ignore",     -- Interferes with <C-r><C-l> to toggle relative numbers
        ["<leader>"] = {
            d = {
                name         = "+Debug",
                d            = {"Launch debugger"},
                e            = {"Close / Reset debugger"},
                ["d_"]       = {"Restart debugger"},
                l            = {"Step into"},
                j            = {"Step over"},
                k            = {"Step out"},
                ["d<space>"] = {"Continue"},
                rc           = {"Run to cursor"},
                b            = {"Toggle breakpoint"},
                cb           = {"Toggle conditional breakpoint"},
                X            = {"Clear all breakpoints"},
                i            = {"Inspect"},
                c            = {"Code window"},
                t            = {"Tag window"},
                v            = {"Variables window"},
                w            = {"Watches window"},
                s            = {"Stack trace window"},
                o            = {"Output window"}
            },
          s = {
              name           = "+Search",
              s              = "String (prompt)",
              w              = "String (under cursor)",
              f              = "File (all files)",
              b              = "Buffer list",
              h              = "Help tags",
              c              = "Config files",
              gc             = "Git commits",
              gb             = "Git branches",
              gs             = "Git status",
              cs             = "Colour schemes"
          }
        }
    }, defaultOpts)
return M

