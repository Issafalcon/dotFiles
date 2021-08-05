require('vimp')

vim.g.vimspector_base_dir = vim.fn.expand('$HOME/.config/vimspector-config')

function _G.GotoWindow(id)
  vim.call('win_gotoid', id)
end

function _G.AddToWatch()
  local word = vim.fn.expand('<cexpr>')
  vim.call('vimspector#AddWatch', word)
end

-- Debugger remaps
vimp.nnoremap('<leader>m', ':MaximizerToggle!<CR>')
vimp.nnoremap('<leader>dc', ':lua GotoWindow(vim.g.vimspector_session_windows.code)<CR>')
vimp.nnoremap('<leader>dt', ':lua GotoWindow(vim.g.vimspector_session_windows.tagpage)<CR>')
vimp.nnoremap('<leader>dv', ':lua GotoWindow(vim.g.vimspector_session_windows.variables)<CR>')
vimp.nnoremap('<leader>dw', ':lua GotoWindow(vim.g.vimspector_session_windows.watches)<CR>')
vimp.nnoremap('<leader>ds', ':lua GotoWindow(vim.g.vimspector_session_windows.stack_trace)<CR>')
vimp.nnoremap('<leader>do', ':lua GotoWindow(vim.g.vimspector_session_windows.output)<CR>')
vimp.nnoremap('<leader>de', ':call vimspector#Reset()<CR>')
vimp.nnoremap('<leader>d?', ':lua AddToWatch()<CR>')

vimp.nmap('<leader>dl', '<Plug>VimspectorStepInto')
vimp.nmap('<leader>dj', '<Plug>VimspectorStepOver')
vimp.nmap('<leader>dk', '<Plug>VimspectorStepOut')
vimp.nmap('<leader>d_', '<Plug>VimspectorRestart')

vimp.nnoremap('<leader>d<space>', ':call vimspector#Continue()<CR>')

vimp.nmap('<leader>drc', '<Plug>VimspectorRunToCursor')
vimp.nmap('<leader>db', '<Plug>VimspectorToggleBreakpoint')
vimp.nnoremap('<leader>dB', '<Plug>VimspectorToggleConditionalBreakpoint')
vimp.nmap('<leader>dX', ':call vimspector#ClearBreakpoints()<CR>')

-- for normal mode - the word under the cursor
vimp.nmap('<Leader>di', '<Plug>VimspectorBalloonEval')
-- for visual mode, the visually selected text
vimp.xmap('<Leader>di', '<Plug>VimspectorBalloonEval')

-- Mapping to begin debugging for specific 'Modes'
--  dA = Debug Attach (Attach to a running process)
--  dL = Launch a process in debug
--  dT = Run a test in debug mode
--  dd = Launch file specific vimspector.json and select an option
vimp.nnoremap('<leader>dJ', ':lua require"utils.vimspectorHelper".startDebugTest()<CR>')
vimp.nnoremap('<leader>dd', ':call vimspector#Launch()<CR>')

