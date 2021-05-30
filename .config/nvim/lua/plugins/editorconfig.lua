require('nvim_utils')

vim.g.EditorConfig_exclude_patterns = {'fugitive://.*'}

local autoCmds = {
  editoConfig = {
    { "Filetype", "gitcommit", "let b:EditorConfig_disable = 1" };
  };
}

nvim_create_augroups(autoCmds)
