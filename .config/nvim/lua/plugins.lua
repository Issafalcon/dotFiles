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
    use "svermeulen/vimpeccable"
    use {"norcalli/nvim_utils"}

    -- LSP
    use "neovim/nvim-lspconfig"
    use "williamboman/nvim-lsp-installer" -- Adds missing lspinstall command with some buindles LSP servers
    use "glepnir/lspsaga.nvim" -- Adds nice UIs and shortcut keys for LSP interactions
    use "kosayoda/nvim-lightbulb" -- Adds the code actions lightbulb with integration with codeaction menu

    -- EditorConfig
    use 'editorconfig/editorconfig-vim'

    -- Quickfix / Location list tools
    use {"tpope/vim-unimpaired"}
    use {"kevinhwang91/nvim-bqf"}
    use {'junegunn/fzf'}

    -- Filesystem navigation
    -- use {"kyazdani42/nvim-tree.lua"} -- Toggle hidden files and refresh tree aren't working currently (25/05/21)
    use 'preservim/nerdtree'
    use 'Xuyuanp/nerdtree-git-plugin'
    use 'tiagofumo/vim-nerdtree-syntax-highlight'
    use 'PhilRunninger/nerdtree-visual-selection'

    -- Colours and fonts and icons
    use "ryanoasis/vim-devicons"
    use "kyazdani42/nvim-web-devicons"
    use "christianchiarulli/nvcode-color-schemes.vim"
    use "norcalli/nvim-colorizer.lua"

    -- Vim undo history
    use {"mbbill/undotree"}

    -- Git integration
    use {"airblade/vim-gitgutter"}
    use {"tpope/vim-fugitive"}
    use "sodapopcan/vim-twiggy" -- Fugitive extension to manage branches
    use {"tpope/vim-rhubarb"} -- Browse Github URLs
    use {"lewis6991/gitsigns.nvim"}
    use "junegunn/gv.vim" -- Git commit browser
    use "rhysd/git-messenger.vim" -- Show commits under the cursor

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

    -- Fuzzy searching
    use {"nvim-telescope/telescope-fzy-native.nvim"}
    use {
      "nvim-telescope/telescope.nvim",
      requires = {
        {"nvim-lua/popup.nvim"},
        {"nvim-lua/plenary.nvim"}
      }
    }

    -- .NET Development
    use {"OmniSharp/omnisharp-vim"} -- Omnisharpe (.NET LSP interface) - Use to install omnisharpe-rosyln LSP

    -- Testing tools
    use {"vim-test/vim-test"}

    -- Terraform
    use "hashivim/vim-terraform"
    use "juliosueiras/vim-terraform-completion"

    -- Debugging
    use {"puremourning/vimspector"}

    -- Autocomplete and Snippets
    use "hrsh7th/nvim-compe"
    use "hrsh7th/vim-vsnip"
    use "hrsh7th/vim-vsnip-integ"
    use "rafamadriz/friendly-snippets"
    use "J0rgeSerran0/vscode-csharp-snippets"

    -- Treesitter
    use {"nvim-treesitter/nvim-treesitter", run = ":TSUpdate"}
    use {"nvim-treesitter/nvim-treesitter-angular"}
    use {"windwp/nvim-ts-autotag"}
    use {"p00f/nvim-ts-rainbow"}

    -- StatusLine and bufferline
    use {"glepnir/galaxyline.nvim"}

    -- Markdown tools
    use {"iamcco/markdown-preview.nvim", run = "cd app && npm install"}

    -- Wiki and Note taking
    use 'lervag/vimtex'

    -- Utility
    use {"terrortylor/nvim-comment"}
    use {"tpope/vim-surround"}
    use {"sheerun/vim-polyglot"} -- Better syntax support
    use "junegunn/vim-easy-align"
    use "folke/which-key.nvim" -- Key binding support
    use "windwp/nvim-autopairs"
  end
)
