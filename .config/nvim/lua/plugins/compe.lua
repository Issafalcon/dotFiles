local utils = require('utils')
local wk = require("which-key")

require'compe'.setup {
    enabled = true,
    autocomplete = true,
    debug = false,
    min_length = 1,
    preselect = 'enable',
    throttle_time = 80,
    source_timeout = 200,
    incomplete_delay = 400,
    max_abbr_width = 100,
    max_kind_width = 100,
    max_menu_width = 100,
    documentation = true,

    source = {
        path = {kind = "   (Path)"},
        buffer = {kind = "   (Buffer)"},
        calc = {kind = "   (Calc)"},
        vsnip = {kind = "   (Snippet)"},
        nvim_lsp = {kind = "   (LSP)"},
        nvim_lua = {kind = "  "},
        spell = {kind = "   (Spell)", filetypes={"markdown", "text"}},
        tags = false,
        vim_dadbod_completion = true,
        snippets_nvim = {kind = "  "},
        treesitter = {kind = "  "},
        emoji = {kind = " ﲃ  (Emoji)", filetypes={"markdown", "text"}}
        -- for emoji press : (idk if that in compe tho)
    }
}

-- Below is taken directly from https://github.com/hrsh7th/nvim-compe config
local check_back_space = function()
    local col = vim.fn.col('.') - 1
    if col == 0 or vim.fn.getline('.'):sub(col, col):match('%s') then
        return true
    else
        return false
    end
end

-- Use (s-)tab to:
-- move to prev/next item in completion menuone
-- jump to prev/next snippet's placeholder
_G.tab_complete = function()
    if vim.fn.pumvisible() == 1 then
        return t"<C-n>"
    elseif vim.fn.call("vsnip#available", {1}) == 1 then
        return t"<Plug>(vsnip-expand-or-jump)"
    elseif check_back_space() then
        return t"<Tab>"
    else
        return vim.fn['compe#complete']()
    end
end

_G.s_tab_complete = function()
    if vim.fn.pumvisible() == 1 then
        return t"<C-p>"
    elseif vim.fn.call("vsnip#jumpable", {-1}) == 1 then
        return t"<Plug>(vsnip-jump-prev)"
    else
        return t"<S-Tab>"
    end
end

-- Key Mappings for Which-Key
vim.api.nvim_set_keymap("i", "<CR>", "compe#confirm('<CR>')", {expr = true, silent = true })
vim.api.nvim_set_keymap("i", "<C-space>", "compe#complete()", {expr = true, silent = true })
vim.api.nvim_set_keymap("i", "<C-e>", "compe#close('<C-e>')", {expr = true})
vim.api.nvim_set_keymap("i", "<C-f>", "compe#scroll({ 'delta': +4 })", {expr = true, noremap = true})
vim.api.nvim_set_keymap("i", "<C-b>", "compe#scroll({ 'delta': -4 })", {expr = true, noremap = true})

vim.api.nvim_set_keymap("i", "<Tab>", "v:lua.tab_complete()", {expr = true, noremap = true})
vim.api.nvim_set_keymap("i", "<Tab>", "v:lua.tab_complete()", {expr = true, noremap = true})
vim.api.nvim_set_keymap("s", "<Tab>", "v:lua.tab_complete()", {expr = true})
vim.api.nvim_set_keymap("i", "<S-Tab>", "v:lua.s_tab_complete()", {expr = true})
vim.api.nvim_set_keymap("s", "<S-Tab>", "v:lua.s_tab_complete()", {expr = true})
