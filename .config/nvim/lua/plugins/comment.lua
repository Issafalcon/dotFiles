require("nvim_comment").setup()
local wk = require("which-key")

wk.register({
        ["<leader>"] = {
            ["/"] = {":CommentToggle<CR>", "Comment Toggle", mode = "n"},
        }
    })

wk.register({
        ["<leader>"] = {
            ["/"] = {":CommentToggle<CR>", "Comment Toggle", mode = "v"},
        }
    })
