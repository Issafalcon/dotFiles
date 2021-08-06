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
  vim.cmd("!" .. cmd)
  vim.cmd("let $VSTEST_HOST_DEBUG=0")
end

-- Custom Strategies For Testing
local strategies = {}
strategies["debugJest"] = DebugJestFile
strategies["debugDotNet"] = DebugDotnetFile

vim.g["test#custom#strategies"] = strategies
