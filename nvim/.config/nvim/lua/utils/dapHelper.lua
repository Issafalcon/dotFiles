-- At time of writing, this file isn't currently in use
-- I intend to come back to it when nvim-dap is more mature / I know what I'm doing a bit better and can configure it
-- to work properly
local dap = require("dap")
local dapUtils = require("dap.utils")

local function debugInChrome()
  local launchUrl = vim.fn.input("Launch URL - Full path or relative to http://localhost: ")

  if string.sub(launchUrl, 1, 1) == "/" then
    launchUrl = "http://localhost" .. launchUrl
  end

  local sourcemapKeys = { "/./*", "/src/*", "webpack:///./*", }
  local webrootMapPath = "${webRoot}/*"
  local sourceMapOverrides = {}
  sourceMapOverrides[sourcemapKeys[1]] = webrootMapPath
  sourceMapOverrides[sourcemapKeys[2]] = webrootMapPath
  sourceMapOverrides[sourcemapKeys[3]] = webrootMapPath
       -- "/*": "*",
       -- "/./~/*": "${webRoot}/node_modules/*"
  dap.run(
    {
      type = "chrome",
      request = "launch",
      stopOnEntry = true,
      url = launchUrl,
      webRoot = "${workspaceFolder}",
      runtimeExecutable = "/usr/bin/google-chrome",
      sourceMapPathOverrides = sourceMapOverrides
    }
  )
end

local function debugDotNet()
  dap.set_log_level('TRACE')

  dap.run(
    {
      type = "netcoredbg",
      name = "attach - netcoredbg",
      request = "attach",
      processId = function()
        return dapUtils.pick_process()
      end,
      -- program = function()
      --   return vim.fn.input('Path to dll', vim.fn.getcwd() .. '/bin/Debug/', 'file')
      -- end,
    }
  )
end

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

return {
  debugInChrome = debugInChrome,
  startDebugging = startDebugging
}
