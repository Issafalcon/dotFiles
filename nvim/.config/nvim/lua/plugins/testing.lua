require('vimp')

function _G.DebugJestFile(cmd)
  print(cmd)
  local cwd = vim.fn.getcwd()
  vim.cmd('call vimspector#LaunchWithSettings( #{ configuration: "Jest watch current file", VimCwd: "' .. cwd .. '"})')
end

-- Custom Strategies For Testing
local strategies = {}
strategies["debugJest"] = DebugJestFile

vim.g["test#custom_strategies"] = strategies

-- Ultest Keybindings
vimp.nmap('<leader>us', '<Plug>(ultest-summary-toggle)')
vimp.nmap('<leader>uf', '<Plug>(ultest-run-file)')
vimp.nmap('<leader>un', '<Plug>(ultest-run-nearest)')
vimp.nmap('<leader>uc', ':UltestClear<CR>')
vimp.nmap('<leader>uo', '<Plug>(ultest-output-jump)')
vimp.nmap('<leader>[u', '<Plug>(ultest-prev-fail)')
vimp.nmap('<leader>]u', '<Plug>(ultest-next-fail)')
