require('vimp')

function _G.DebugJestFile(cmd)
  print(cmd)
  local cwd = vim.fn.getcwd()
  vim.cmd('call vimspector#LaunchWithSettings( #{ configuration: "Jest watch current file", VimCwd: "' .. cwd .. '"})')
end

function _G.DebugDotnetFile(cmd)
  print(cmd)
  --Fix this
  vim.cmd("let $VSTEST_HOST_DEBUG=1")
  vim.cmd("echo $VSTEST_HOST_DEBUG")
  vim.cmd("terminal " .. cmd)
  vim.cmd("let $VSTEST_HOST_DEBUG=0")
end

-- Custom Strategies For Testing
local strategies = {}
strategies["debugJest"] = DebugJestFile
strategies["debugDotNet"] = DebugDotnetFile

vim.g["test#custom_strategies"] = strategies

-- Ultest Keybindings
vimp.nmap('<leader>us', '<Plug>(ultest-summary-toggle)')
vimp.nmap('<leader>uf', '<Plug>(ultest-run-file)')
vimp.nmap('<leader>un', '<Plug>(ultest-run-nearest)')
vimp.nmap('<leader>uc', ':UltestClear<CR>')
vimp.nmap('<leader>uo', '<Plug>(ultest-output-jump)')
vimp.nmap('<leader>[u', '<Plug>(ultest-prev-fail)')
vimp.nmap('<leader>]u', '<Plug>(ultest-next-fail)')
