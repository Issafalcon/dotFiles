-- Automatically generated packer.nvim plugin loader code

if vim.api.nvim_call_function('has', {'nvim-0.5'}) ~= 1 then
  vim.api.nvim_command('echohl WarningMsg | echom "Invalid Neovim version for packer.nvim! | echohl None"')
  return
end

vim.api.nvim_command('packadd packer.nvim')

local no_errors, error_msg = pcall(function()

  local time
  local profile_info
  local should_profile = false
  if should_profile then
    local hrtime = vim.loop.hrtime
    profile_info = {}
    time = function(chunk, start)
      if start then
        profile_info[chunk] = hrtime()
      else
        profile_info[chunk] = (hrtime() - profile_info[chunk]) / 1e6
      end
    end
  else
    time = function(chunk, start) end
  end
  
local function save_profiles(threshold)
  local sorted_times = {}
  for chunk_name, time_taken in pairs(profile_info) do
    sorted_times[#sorted_times + 1] = {chunk_name, time_taken}
  end
  table.sort(sorted_times, function(a, b) return a[2] > b[2] end)
  local results = {}
  for i, elem in ipairs(sorted_times) do
    if not threshold or threshold and elem[2] > threshold then
      results[i] = elem[1] .. ' took ' .. elem[2] .. 'ms'
    end
  end

  _G._packer = _G._packer or {}
  _G._packer.profile_output = results
end

time([[Luarocks path setup]], true)
local package_path_str = "/home/issafalcon/.cache/nvim/packer_hererocks/2.1.0-beta3/share/lua/5.1/?.lua;/home/issafalcon/.cache/nvim/packer_hererocks/2.1.0-beta3/share/lua/5.1/?/init.lua;/home/issafalcon/.cache/nvim/packer_hererocks/2.1.0-beta3/lib/luarocks/rocks-5.1/?.lua;/home/issafalcon/.cache/nvim/packer_hererocks/2.1.0-beta3/lib/luarocks/rocks-5.1/?/init.lua"
local install_cpath_pattern = "/home/issafalcon/.cache/nvim/packer_hererocks/2.1.0-beta3/lib/lua/5.1/?.so"
if not string.find(package.path, package_path_str, 1, true) then
  package.path = package.path .. ';' .. package_path_str
end

if not string.find(package.cpath, install_cpath_pattern, 1, true) then
  package.cpath = package.cpath .. ';' .. install_cpath_pattern
end

time([[Luarocks path setup]], false)
time([[try_loadstring definition]], true)
local function try_loadstring(s, component, name)
  local success, result = pcall(loadstring(s))
  if not success then
    vim.schedule(function()
      vim.api.nvim_notify('packer.nvim: Error running ' .. component .. ' for ' .. name .. ': ' .. result, vim.log.levels.ERROR, {})
    end)
  end
  return result
end

time([[try_loadstring definition]], false)
time([[Defining packer_plugins]], true)
_G.packer_plugins = {
  ["BetterLua.vim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/BetterLua.vim"
  },
  ["editorconfig-vim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/editorconfig-vim"
  },
  ["friendly-snippets"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/friendly-snippets"
  },
  fzf = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/fzf"
  },
  ["galaxyline.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/galaxyline.nvim"
  },
  ["git-messenger.vim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/git-messenger.vim"
  },
  ["gitsigns.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/gitsigns.nvim"
  },
  ["gv.vim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/gv.vim"
  },
  ["lsp_signature.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/lsp_signature.nvim"
  },
  ["lspsaga.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/lspsaga.nvim"
  },
  ["markdown-preview.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/markdown-preview.nvim"
  },
  ["nlua.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nlua.nvim"
  },
  ["nvcode-color-schemes.vim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvcode-color-schemes.vim"
  },
  ["nvim-autopairs"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-autopairs"
  },
  ["nvim-bqf"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-bqf"
  },
  ["nvim-colorizer.lua"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-colorizer.lua"
  },
  ["nvim-comment"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-comment"
  },
  ["nvim-compe"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-compe"
  },
  ["nvim-lsp-installer"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-lsp-installer"
  },
  ["nvim-lspconfig"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-lspconfig"
  },
  ["nvim-treesitter"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-treesitter"
  },
  ["nvim-treesitter-angular"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-treesitter-angular"
  },
  ["nvim-ts-autotag"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-ts-autotag"
  },
  ["nvim-ts-rainbow"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-ts-rainbow"
  },
  ["nvim-web-devicons"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim-web-devicons"
  },
  nvim_utils = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/nvim_utils"
  },
  ["octo.nvim"] = {
    config = { "\27LJ\2\n2\0\0\3\0\3\0\0066\0\0\0'\2\1\0B\0\2\0029\0\2\0B\0\1\1K\0\1\0\nsetup\tocto\frequire\0" },
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/octo.nvim"
  },
  ["omnisharp-vim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/omnisharp-vim"
  },
  ["packer.nvim"] = {
    loaded = false,
    needs_bufread = false,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/opt/packer.nvim"
  },
  ["plenary.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/plenary.nvim"
  },
  ["popup.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/popup.nvim"
  },
  rnvimr = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/rnvimr"
  },
  ["suckless.vim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/suckless.vim"
  },
  tagbar = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/tagbar"
  },
  ["telescope-fzy-native.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/telescope-fzy-native.nvim"
  },
  ["telescope.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/telescope.nvim"
  },
  undotree = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/undotree"
  },
  ["vim-cucumber"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-cucumber"
  },
  ["vim-devicons"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-devicons"
  },
  ["vim-dispatch"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-dispatch"
  },
  ["vim-easy-align"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-easy-align"
  },
  ["vim-easymotion"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-easymotion"
  },
  ["vim-fugitive"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-fugitive"
  },
  ["vim-gitgutter"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-gitgutter"
  },
  ["vim-maximizer"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-maximizer"
  },
  ["vim-polyglot"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-polyglot"
  },
  ["vim-rhubarb"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-rhubarb"
  },
  ["vim-surround"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-surround"
  },
  ["vim-terraform"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-terraform"
  },
  ["vim-terraform-completion"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-terraform-completion"
  },
  ["vim-test"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-test"
  },
  ["vim-twiggy"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-twiggy"
  },
  ["vim-ultest"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-ultest"
  },
  ["vim-unimpaired"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-unimpaired"
  },
  ["vim-vsnip"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-vsnip"
  },
  ["vim-vsnip-integ"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vim-vsnip-integ"
  },
  vimpeccable = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vimpeccable"
  },
  vimspector = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vimspector"
  },
  vimtex = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vimtex"
  },
  vimwiki = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vimwiki"
  },
  ["vscode-csharp-snippets"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vscode-csharp-snippets"
  },
  ["vscode-markdown-snippets"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/vscode-markdown-snippets"
  },
  ["which-key.nvim"] = {
    loaded = true,
    path = "/home/issafalcon/.local/share/nvim/site/pack/packer/start/which-key.nvim"
  }
}

time([[Defining packer_plugins]], false)
-- Config for: octo.nvim
time([[Config for octo.nvim]], true)
try_loadstring("\27LJ\2\n2\0\0\3\0\3\0\0066\0\0\0'\2\1\0B\0\2\0029\0\2\0B\0\1\1K\0\1\0\nsetup\tocto\frequire\0", "config", "octo.nvim")
time([[Config for octo.nvim]], false)
if should_profile then save_profiles() end

end)

if not no_errors then
  vim.api.nvim_command('echohl ErrorMsg | echom "Error in packer_compiled: '..error_msg..'" | echom "Please check your config for correctness" | echohl None')
end
