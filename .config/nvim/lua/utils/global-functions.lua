-- Inspect objects
function _G.dump(...)
    local objects = vim.tbl_map(vim.inspect, {...})
    print(unpack(objects))
end

-- Termcodes escapiing function
function _G.t(str)
  return vim.api.nvim_replace_termcodes(str, true, true, true)
end

-- Autocommits the current file
function _G.AutoCommit()
  local output = vim.api.nvim_exec(':Git rev-parse --git-dir', true)
  print(output)
  if output ~= ".git" then
    print("Error: Not a git repository")
    return
  end
  local message = 'Updated ' .. vim.fn.expand('%:.')
  print(message)
  vim.cmd('Git add ' .. vim.fn.expand('%:p'))
  vim.cmd('Git commit -m ' .. vim.fn.shellescape(message, 1))
end
