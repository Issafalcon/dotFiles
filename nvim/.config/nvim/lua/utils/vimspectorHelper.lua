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
    -- Needs reworking to run tests
    vim.cmd('TestFile -strategy=debugJest')
  elseif vim.bo.filetype == "cs" then
    vim.cmd('TestNearest -strategy=debugDotNet')
    vim.cmd('buffer #')
    vim.cmd('call vimspector#LaunchWithSettings( #{ configuration: "netcoredbg attach"})')
    vim.cmd('let $VSTEST_HOST_DEBUG=0')
  end  
end

return {
  startDebugAttach = startDebugAttach,
  startDebugLaunch = startDebugLaunch,
  startDebugTest = startDebugTest,
}
