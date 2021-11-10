local fn = vim.fn
local execute = vim.api.nvim_command

-- Auto install packer.nvim if it doesn't exist
local install_path = fn.stdpath("data") .. "/site/pack/packer/opt/packer.nvim"

if fn.empty(fn.glob(install_path)) > 0 then
  execute("!git clone https://github.com/wbthomason/packer.nvim.git " .. install_path)
end

vim.cmd [[packadd packer.nvim]]
vim.cmd "autocmd BufWritePost plugins.lua PackerCompile" -- Autocompile changes to plugins.lua

return require("packer").startup(
  function(use)
    -- Packer can manage itself as optional plugin
    use {"wbthomason/packer.nvim", opt = true}

    -- Lua development
    use {"tjdevries/nlua.nvim"}
    use {"euclidianAce/BetterLua.vim"}
    use {"svermeulen/vimpeccable"}
    use {"norcalli/nvim_utils"}

    -- LSP
    use {"neovim/nvim-lspconfig"}
    use {"williamboman/nvim-lsp-installer"} -- Adds missing lspinstall command with some buindles LSP servers
    -- Original lspsaga looks like it isn't being maintained. Incompatible with newer versions of neovim.
    -- Swap to fork for nvim > 0.51
    use {"glepnir/lspsaga.nvim"} -- Adds nice UIs and shortcut keys for LSP interactions
    -- use { 'tami5/lspsaga.nvim' }
    use {"ray-x/lsp_signature.nvim"}
    use {"onsails/lspkind-nvim"}

    -- EditorConfig
    use "editorconfig/editorconfig-vim"

    -- Quickfix / Location list tools
    use {"tpope/vim-unimpaired"}
    use {"kevinhwang91/nvim-bqf"}
    use {"junegunn/fzf"}

    -- Colours and fonts and icons
    use "ryanoasis/vim-devicons"
    use "kyazdani42/nvim-web-devicons"
    use "christianchiarulli/nvcode-color-schemes.vim"
    use "norcalli/nvim-colorizer.lua"
    use {"yonlu/omni.vim"}
    use {"marko-cerovac/material.nvim"}
    use {"bluz71/vim-nightfly-guicolors"}

    -- Vim undo history
    use {"mbbill/undotree"}

    -- Git integration
    use {"airblade/vim-gitgutter"}
    use {"tpope/vim-fugitive"}
    use {"tpope/vim-rhubarb"} -- Browse Github URLs
    use "junegunn/gv.vim" -- Git commit browser
    use "tveskag/nvim-blame-line"
    use "rhysd/git-messenger.vim" -- Show commits under the cursor
    use {"kdheepak/lazygit.nvim"}
    use "sindrets/diffview.nvim"

    -- Github
    use {
      "pwntester/octo.nvim",
      config = function()
        require "octo".setup()
      end
    }

    -- Window resizing
    use {"fabi1cazenave/suckless.vim"}
    use {"szw/vim-maximizer"}

    -- File navigation and searching
    use {"nvim-telescope/telescope-fzy-native.nvim"}
    use {
      "nvim-telescope/telescope.nvim",
      requires = {
        {"nvim-lua/popup.nvim"},
        {"nvim-lua/plenary.nvim"}
      }
    }
    use {"kevinhwang91/rnvimr"}
    use {"preservim/tagbar"}
    use {"phaazon/hop.nvim", as = "hop"}

    -- .NET Development
    use {"OmniSharp/omnisharp-vim"} -- Omnisharpe (.NET LSP interface) - Use to install omnisharpe-rosyln LSP

    -- Testing tools
    use {"rcarriga/vim-ultest", requires = {"vim-test/vim-test"}, run = ":UpdateRemotePlugins"}

    -- Terraform
    use "hashivim/vim-terraform"
    use "juliosueiras/vim-terraform-completion"

    -- Gherkin / Cucumber
    use "tpope/vim-cucumber"

    -- Debugging
    use {"puremourning/vimspector"}
    -- use {'mfussenegger/nvim-dap'}

    -- Autocomplete and Snippets
    use {"hrsh7th/vim-vsnip"}
    use {"hrsh7th/cmp-vsnip"}
    use {"hrsh7th/cmp-nvim-lsp"}
    use {"hrsh7th/cmp-buffer"}
    use {"hrsh7th/cmp-path"}
    use {"hrsh7th/cmp-cmdline"}
    use {"hrsh7th/nvim-cmp"}
    use {"rafamadriz/friendly-snippets"}
    use {"J0rgeSerran0/vscode-csharp-snippets"}
    use {"robole/vscode-markdown-snippets"}
    use {"xabikos/vscode-react"}
    use {"dsznajder/vscode-es7-javascript-react-snippets"}
    use {"mattn/emmet-vim"} -- Install this only for javascript (React) or jsx files support, which emmet-ls doesn't support

    -- Treesitter
    use {"nvim-treesitter/nvim-treesitter", run = ":TSUpdate"}
    use {"nvim-treesitter/nvim-treesitter-textobjects"}
    use {"nvim-treesitter/playground"}
    use {"nvim-treesitter/nvim-treesitter-angular"}
    use {"windwp/nvim-ts-autotag"}
    use {"p00f/nvim-ts-rainbow"}

    -- StatusLine and bufferline
    use {"glepnir/galaxyline.nvim"}

    -- Markdown tools
    use {"iamcco/markdown-preview.nvim", run = "cd app && npm install"}

    -- Wiki and Note taking
    use {"lervag/vimtex"}
    use {"vimwiki/vimwiki"}

    -- Marks management
    use {"chentau/marks.nvim"}

    -- Utility
    use {"terrortylor/nvim-comment"}
    use {"tpope/vim-surround"}
    use {"sheerun/vim-polyglot"} -- Better syntax support
    use {"junegunn/vim-easy-align"}
    use {"folke/which-key.nvim"} -- Key binding support
    use {
      "sudormrfbin/cheatsheet.nvim",
      requires = {
        {"nvim-telescope/telescope.nvim"},
        {"nvim-lua/popup.nvim"},
        {"nvim-lua/plenary.nvim"}
      }
    }
    use {"windwp/nvim-autopairs"}
    use {"tpope/vim-dispatch"}
    use {"dbeniamine/cheat.sh-vim"}
    use {"kkoomen/vim-doge", run = ":call doge#install()"}
    use {
      "rmagatti/session-lens",
      requires = {"rmagatti/auto-session", "nvim-telescope/telescope.nvim"},
      config = function()
        require("session-lens").setup()
      end
    }
  end
)
