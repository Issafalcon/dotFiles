local map = vim.api.nvim_set_keymap
local o = vim.o
local g = vim.g
local wo = vim.wo
local bo = vim.bo
local cmd = vim.cmd
local fn = vim.fn

g.mapleader                  = ' '             -- set the leader
g.vimspector_enable_mappings = "visual_studio" -- set vimspector mappings to match vscode debugger style
terminal                     = fn.expand('$terminal')

cmd('filetype plugin indent on')
cmd 'set noswapfile'
cmd 'set iskeyword+=-'                         -- treat dash separated words as a word text object"
cmd 'set iskeyword+=@'                        -- Allow coc-css to recognize scss keywords
cmd 'set shortmess+=c'                         -- don't pass messages to |ins-completion-menu|.
cmd 'set inccommand=split'                     -- make substitution work in realtime
cmd 'syntax on'
cmd 'let &titleold="\'..terminal..\'"'
cmd 'set colorcolumn=99999'                    -- fix indentline for now
cmd('set ts=2')                                -- insert 2 spaces for a tab
cmd('set sw=2')                                -- change the number of space characters inserted for indentation

o.hidden          = true                       -- required to keep multiple buffers open
o.title           = true
o.mouse           = "a"                        -- enable your mouse
o.titlestring     = "%<%f%=%l/%l - nvim"
o.guifont         = "firacode nerd font:h17"   -- only used for gui vim clients
o.pumheight       = 10                         -- makes popup menu smaller
o.fileencoding    = "utf-8"                    -- the encoding written to file
o.cmdheight       = 2                          -- more space for displaying messages
o.splitbelow      = true                       -- horizontal splits will automatically be below
o.splitright      = true                       -- vertical splits will automatically be to the right
o.termguicolors   = true                       -- set term giu colors most terminals support this
o.conceallevel    = 0                          -- so that i can see `` in markdown files
o.showtabline     = 2                          -- always show tabs
o.backup          = false                      -- this is recommended by coc
o.writebackup     = false                      -- this is recommended by coc
o.showmode        = true
o.updatetime      = 300                        -- faster completion
o.timeoutlen      = 500                        -- by default timeoutlen is 1000 ms
o.clipboard       = "unnamedplus"              -- copy paste between vim and everything else
o.expandtab       = true                       -- converts tabs to spaces
o.tabstop         = 2
o.shiftwidth      = 2
o.scrolloff       = 8
o.ignorecase      = true
o.smartcase       = true
o.hlsearch        = false
wo.scrolloff      = 8
g.scrolloff       = 8
wo.wrap           = false                      -- Display long lines as just one line
wo.number         = true                       -- set numbered lines
wo.signcolumn     = "yes:2"                      -- Always show the signcolumn, otherwise it would shift the text each time

bo.expandtab      = true 							-- Converts tabs to spaces
bo.smartindent    = true 							-- Makes indenting smart

-- LSP
-- menuone: popup even when there's only one match
-- noinsert: Do not insert text until a selection is made
-- noselect: Do not select, force user to select one from the menu
o.completeopt = "menuone,noinsert,noselect"

-- Visualize diagnostics
g.diagnostic_enable_virtual_text = 1
g.diagnostic_trimmed_virtual_text = "40"

-- Don't show diagnostics while in insert mode
g.diagnostic_insert_delay = 1

-- FOLDS
o.foldlevel = 99

-- Session
o.sessionoptions="blank,buffers,curdir,folds,help,options,tabpages,winsize,resize,winpos,terminal"

-- Backup, undo, swap options
if fn.has('persistent_undo') then
    cmd('set undodir=~/.undodir')
    cmd('set undofile')
end

if fn.has('win32unix') and fn.has('wsl') then
    g.clipboard = {
        name = "win32yank-wsl",
        copy = {
            ["+"] = "win32yank.exe -i --crlf",
            ["*"] = "win32yank.exe -i --crlf"
        },
        paste = {
            ["+"] = "win32yank.exe -o --lf",
            ["*"] = "win32yank.exe -o --lf"
        }
    }
end

-- URL handling. Opens urls in browser with gx
if vim.fn.has("mac") == 1 then
  map('', 'gx', '<Cmd>call jobstart(["open", expand("<cfile>")], {"detach": v:true})<CR>', {})
elseif vim.fn.has("unix") == 1 then
  map('', 'gx', '<Cmd>call jobstart(["xdg-open", expand("<cfile>")], {"detach": v:true})<CR>', {})
else
  map('', 'gx', '<Cmd>lua print("Error: gx is not supported on this OS!")', {})
end
