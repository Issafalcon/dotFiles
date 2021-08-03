local dap = require("dap")

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

local function startDebugging()
  if dap.session() then
    dap.disconnect()
    dap.close()
  end

  if vim.bo.filetype == "typescript" then
    debugInChrome()
  end
end

return {
  debugInChrome = debugInChrome,
  startDebugging = startDebugging
}
