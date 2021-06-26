require'nvim-web-devicons'.setup {
 -- your personnal icons can go here (to override)
 -- DevIcon will be appended to `name`
 override = {
    tex = {
      icon = "",
      color = "#3D6117",
      name = "Tex" 
    }
 };
 -- globally enable default icons (default to false)
 -- will get overriden by `get_icons` option
 default = true;
}

-- Workaround for a bug in nerd-tree devicons where tex files move across the screen
vim.cmd('let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols = {}')
vim.cmd('let g:WebDevIconsUnicodeDecorateFileNodesExtensionSymbols["tex"] = "ƛ"')
