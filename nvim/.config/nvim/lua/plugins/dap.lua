require('nvim_utils')
local dap = require('dap')

-- Adapters
dap.adapters.node2 = {
  type = 'executable',
  command = 'node',
  args = {os.getenv('HOME') .. '/debug-adapters/vscode-node-debug2/out/src/nodeDebug.js'}
}

dap.adapters.chrome = {
  type = 'executable',
  command = 'node',
  args = {os.getenv('HOME') .. '/debug-adapters/vscode-chrome-debug/out/src/chromeDebug.js'}
}

-- Configuration
dap.configurations.typescript = {
  {
    type = 'chrome',
    request = 'launch',
    program = '${file}',
    cwd = vim.fn.getcwd(),
    sourceMaps = true,
    breakOnLoad = true,
    protocol = 'inspector',
    port = 9222,
    webRoot = '${workspaceFolder}/src',
  }
}

-- Mappings
nvim.fn.sign_define('DapBreakPoint', {text='⭕', texthl='', linehl='', numhl=''})
nvim.fn.sign_define('DapStopped', {text='🟢', texthl='', linehl='', numhl=''})

vimp.nmap('<leader>db', ':lua require"dap".toggle_breakpoint()<CR>')
vimp.nmap('<leader>dk', ':lua require"dap".step_out()<CR>')
vimp.nmap('<leader>dl', ':lua require"dap".step_into()<CR>')
vimp.nmap('<leader>dj', ':lua require"dap".step_over()<CR>')
vimp.nmap('<leader>dp', ':lua require"dap".up()<CR>')
vimp.nmap('<leader>dn', ':lua require"dap".down()<CR>')
vimp.nmap('<leader>d_', ':lua require"dap".disconnect();require"dap".stop();require"dap".run_last()<CR>')
vimp.nmap('<leader>dr', ':lua require"dap".repl.open({}, "vsplit")<CR><C-w>l')
vimp.nmap('<leader>di', ':lua require"dap.ui.variables".hover()<CR>')
vimp.vmap('<leader>di', ':lua require"dap.ui.variables".visual_hover()<CR>')
vimp.nmap('<leader>d?', ':lua require"dap.ui.variables".scopes()<CR>')
vimp.nmap('<leader>de', ':lua require"dap".set_exception_breakpoints({"all"})<CR>')

vimp.nmap('<leader>dd', ':lua require"utils.debugHelper".startDebugging()<CR>')

