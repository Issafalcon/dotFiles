local dap = require('dap')

local function debugInChrome()
  local launchUrl = vim.fn.input("Launch URL - Full path or relative to http://localhost: ")

  if string.sub(launchUrl, 1, 1) == '/' then
    launchUrl = 'http://localhost' .. launchUrl
  end

  dap.run({
    type = 'chrome',
    request = 'launch',
    stopOnEntry = true,
    cwd = vim.fn.getcwd(),
    url = launchUrl,
    webRoot = '${workspaceFolder}',
    runtimeExecutable = '/usr/bin/google-chrome'
  })
end

local function startDebugging()
 if dap.session() then
   dap.disconnect()
   dap.close()
  end

  if vim.bo.filetype == 'typescript' then
   debugInChrome()
  end
end

return {
  debugInChrome = debugInChrome,
  startDebugging = startDebugging
}
