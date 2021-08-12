local excludeDirsFunc = function(opts)
  local excludeDirs = {}
  local dir = ""

  repeat
    dir = vim.fn.input("Exclude Dir (be sure to add trailing '/') > ", "", "file")
    if dir ~= nil and dir ~= "" then
      table.insert(excludeDirs, "--glob=!" .. dir)
    end
  until dir == ""

  for index, direc in ipairs(excludeDirs) do
    print(direc)
  end
  return excludeDirs
end

local function custom_live_grep()
  local includeDirs = {}
  local dir = ""

  repeat
    dir = vim.fn.input("Include Dir > ", "", "file")
    if dir ~= nil and dir ~= "" then
      table.insert(includeDirs, dir)
    end
  until dir == ""

  require("telescope.builtin").live_grep({search_dirs = includeDirs, additional_args = excludeDirsFunc})
end

return {
  custom_live_grep = custom_live_grep
}
