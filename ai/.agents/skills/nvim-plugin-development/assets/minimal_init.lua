-- spec/minimal_init.lua
-- Headless test bootstrap for busted via plenary runner.
-- Usage: nvim --headless --clean --noplugin -u spec/minimal_init.lua \
--          -c "lua vim.cmd([[PlenaryBustedDirectory spec/my-plugin { minimal_init = 'spec/minimal_init.lua' }]])"

-- ── 1. Find plenary ──────────────────────────────────────────────────────────
-- Adjust this path to match your system's plenary installation.
-- Common locations:
--   ~/.local/share/nvim/site/pack/*/opt/plenary.nvim  (manual / packer opt)
--   ~/.local/share/nvim/lazy/plenary.nvim             (lazy.nvim)
--   ~/.local/share/nvim/site/pack/*/start/plenary.nvim (start plugins)
local function find_plenary()
  local candidates = {
    vim.fn.expand("~/.local/share/nvim/lazy/plenary.nvim"),
    vim.fn.expand("~/.local/share/nvim/site/pack/core/opt/plenary.nvim"),
    vim.fn.expand("~/.local/share/nvim/site/pack/packer/opt/plenary.nvim"),
  }
  -- Also try a glob search as fallback
  local glob = vim.fn.glob(vim.fn.expand("~/.local/share/nvim") .. "/**/plenary.nvim", false, true)
  for _, p in ipairs(glob) do
    table.insert(candidates, p)
  end
  for _, p in ipairs(candidates) do
    if vim.fn.isdirectory(p) == 1 then
      return p
    end
  end
  return nil
end

local plenary_path = find_plenary()
if not plenary_path then
  error(
    "plenary.nvim not found. Install it or set PLENARY_PATH and update this file.\n"
    .. "  lazy.nvim:  { 'nvim-lua/plenary.nvim' }\n"
    .. "  packer:     use 'nvim-lua/plenary.nvim'\n"
  )
end

-- ── 2. Configure runtimepath ─────────────────────────────────────────────────
vim.opt.rtp:prepend(plenary_path)
vim.opt.rtp:prepend(vim.fn.getcwd()) -- plugin root

-- ── 3. Source plugin entry point (for user-command registration) ─────────────
-- Uncomment if your specs test commands registered in plugin/*.lua:
-- local ok, err = pcall(vim.cmd, "source plugin/my_plugin.lua")
-- if not ok then
--   vim.notify("minimal_init: could not source plugin entry: " .. err, vim.log.levels.WARN)
-- end

-- ── 4. (Optional) Silence noisy startup messages ────────────────────────────
vim.o.shortmess = vim.o.shortmess .. "I"
