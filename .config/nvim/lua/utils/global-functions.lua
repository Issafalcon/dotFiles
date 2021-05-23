-- Inspect objects
function _G.dump(...)
    local objects = vim.tbl_map(vim.inspect, {...})
    print(unpack(objects))
end

-- Termcodes escapiing function
function _G.t(str)
  return vim.api.nvim_replace_termcodes(str, true, true, true)
end
