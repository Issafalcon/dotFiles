require('vimp')

local function startDebuggingDap()
  if dap.session() then
    dap.disconnect()
    dap.close()
  end

  if vim.bo.filetype == "typescript" then
    debugInChrome()
  elseif vim.bo.filetype == "cs" then
    debugDotNet()
  end
end

local function startDebugTest()
  if vim.bo.filetype == "typescript" then
    -- This makes the assumption that the test runner is jest
    vim.g["test#strategy"] = 'debugJest'
    vim.cmd('UltestNearest -strategy=debugJest')
  elseif vim.bo.filetype == "cs" then

    vim.cmd('UltestNearest -strategy=debugDotNet<CR>')
    
    local pid = require('utils.processPicker').pick_process()
    vim.cmd('call vimspector#LaunchWithSettings( #{ configuration: "netcoredbg attach", processId: "' .. pid .. '"})')
  end  
end

return {
  startDebugAttach = startDebugAttach,
  startDebugLaunch = startDebugLaunch,
  startDebugTest = startDebugTest,
}
