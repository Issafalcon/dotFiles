require("nvim_utils")
local dap = require("dap")

-- Adapters
dap.adapters.node2 = {
  type = "executable",
  command = "node",
  args = {os.getenv("HOME") .. "/debug-adapters/vscode-node-debug2/out/src/nodeDebug.js"}
}

dap.adapters.chrome = {
  type = "executable",
  command = "node",
  args = {os.getenv("HOME") .. "/debug-adapters/vscode-chrome-debug/out/src/chromeDebug.js"}
}

dap.adapters.netcoredbg = {
  type = "executable",
  command = "/usr/local/netcoredbg",
  args = {"--interpreter=vscode"}
}

-- Configuration
dap.configurations.typescript = {
  {
    type = "chrome",
    request = "launch",
    program = "${file}",
    cwd = vim.fn.getcwd(),
    sourceMaps = true,
    breakOnLoad = true,
    protocol = "inspector",
    port = 9222,
    webRoot = "${workspaceFolder}/src"
  }
}

dap.configurations.cs = {
  {
    type = "netcoredbg",
    name = "launch - netcoredbg",
    request = "launch",
    stopOnEntry = true,
    -- processId = function()
    --   return dapUtils.pick_process()
    -- end,
    program = function()
      return vim.fn.input("Path to dll", vim.fn.getcwd() .. "/bin/Debug/", "file")
    end
  }
}

-- Mappings
vimp.nmap("<leader>dd", ':lua require"dap".continue()<CR>')
vimp.nmap("<leader>db", ':lua require"dap".toggle_breakpoint()<CR>')
vimp.nmap("<leader>dk", ':lua require"dap".step_out()<CR>')
vimp.nmap("<leader>dl", ':lua require"dap".step_into()<CR>')
vimp.nmap("<leader>dj", ':lua require"dap".step_over()<CR>')
vimp.nmap("<leader>dp", ':lua require"dap".up()<CR>')
vimp.nmap("<leader>dn", ':lua require"dap".down()<CR>')
vimp.nmap("<leader>d_", ':lua require"dap".disconnect();require"dap".stop();require"dap".run_last()<CR>')
vimp.nmap("<leader>di", ':lua require"dap.ui.variables".hover()<CR>')
vimp.vmap("<leader>di", ':lua require"dap.ui.variables".visual_hover()<CR>')
vimp.nmap("<leader>d?", ':lua require"dap.ui.variables".scopes()<CR>')
vimp.nmap("<leader>de", ':lua require"dap".set_exception_breakpoints({"all"})<CR>')

-- Dap windows
vimp.nmap("<leader>dr", ':lua require"dap".repl.open({}, "vsplit")<CR><C-w>l')

vimp.nmap("<F5>", ':lua require"utils.dapHelper".startDebugging()<CR>')
