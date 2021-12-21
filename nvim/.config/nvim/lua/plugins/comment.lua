require("nvim_comment").setup(
  {
    hook = function()
      require("ts_context_commentstring.internal").update_commentstring()
    end
  }
)

require "nvim-treesitter.configs".setup {
  context_commentstring = {
    enable = true,
    enable_autocmd = false
  }
}

local wk = require("which-key")

wk.register(
  {
    ["<leader>"] = {
      ["/"] = {":CommentToggle<CR>", "Comment Toggle", mode = "n"}
    }
  }
)

wk.register(
  {
    ["<leader>"] = {
      ["/"] = {":CommentToggle<CR>", "Comment Toggle", mode = "v"}
    }
  }
)
