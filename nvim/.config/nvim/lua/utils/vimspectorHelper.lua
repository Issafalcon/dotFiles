require("vimp")

local function cwd()
  return vim.fn.getcwd()
end

local function pickPID()
  return require("utils.processPicker").pick_process()
end

local function getUrlInput()
  local launchUrl = vim.fn.input("Launch URL - Full path or relative to http://localhost ")
  local initChar = string.sub(launchUrl, 1, 1)

  if initChar == "/" or initChar == ":" then
    launchUrl = "http://localhost" .. launchUrl
  end

  return launchUrl
end

local function getLaunchFile()
  return vim.fn.input("Path to dll", vim.fn.getcwd() .. "/bin/Debug/", "file")
end

-- Generic debug function. Pick from list of 'launch.json' style options
local function selectDebug()
  -- Provide common parameters (can't use workspaceRoot for example, as we are using a global vimspector.json)
  vim.cmd("call vimspector#LaunchWithSettings( #{ VimCwd: '" .. cwd() .. "'})")
end

local function startDebugAttach()
  if vim.bo.filetype == "typescript" or vim.bo.filetype == "typescriptreact" then
    vim.cmd("call vimspector#LaunchWithSettings( #{ configuration: 'Chrome - Attach', VimCwd: '" .. cwd() .. "'})")
  elseif vim.bo.filetype == "cs" then
    vim.cmd(
      "call vimspector#LaunchWithSettings( #{ configuration: 'netcoredbg attach', VimCwd: '" ..
        cwd() .. "', processId: " .. pickPID() .. "})"
    )
  end
end

local function startDebugLaunch()
  if vim.bo.filetype == "typescript" or vim.bo.filetype == "typescriptreact" then
    vim.cmd(
      "call vimspector#LaunchWithSettings( #{ configuration: 'Chrome - Launch', VimCwd: '" ..
        cwd() .. "', launchUrl: '" .. getUrlInput() .. "'})"
    )
  elseif vim.bo.filetype == "cs" then
    vim.cmd(
      "call vimspector#LaunchWithSettings( #{ configuration: 'netcoredbg', VimCwd: '" ..
        cwd() .. "', dllPath: '" .. getLaunchFile() .. "'})"
    )
  end
end

local function startDebugTest()
  if vim.bo.filetype == "typescript" then
    -- This makes the assumption that the test runner is jest
    vim.cmd("TestFile -strategy=debugJest")
  elseif vim.bo.filetype == "cs" then
    vim.cmd("TestNearest -strategy=debugDotNet")
  end
end

return {
  startDebugAttach = startDebugAttach,
  startDebugLaunch = startDebugLaunch,
  startDebugTest = startDebugTest,
  selectDebug = selectDebug
}
