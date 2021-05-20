require("vimp")

local function custom_attach(client)
  local function buf_set_keymap(...)
    vim.api.nvim_buf_set_keymap(vim.fn.bufnr, ...)
  end
  local function buf_set_option(...)
    vim.api.nvim_buf_set_option(vim.fn.bufnr, ...)
  end

  buf_set_option("omnifunc", "v:lua.vim.lsp.omnifunc")

  -- Mappings.
  local opts = {noremap = true, silent = true}

  -- Native lsp
  if vim.bo.filetype == "cs" then
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>ost <Plug>(omnisharp_type_lookup)]]
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>osd <Plug>(omnisharp_documentation)]]
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>osfs <Plug>(omnisharp_find_symbol)]]
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>osfx <Plug>(omnisharp_fix_usings)]]
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <C-\> <Plug>(omnisharp_signature_help)]]
    --     vim.api.nvim_command [[autocmd FileType cs imap <silent> <buffer> <C-\> <Plug>(omnisharp_signature_help)]]
    --
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>osgcc <Plug>(omnisharp_global_code_check)]]
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>osca <Plug>(omnisharp_code_actions)]]
    --     vim.api.nvim_command [[autocmd FileType cs xmap <silent> <buffer> <Leader>osca <Plug>(omnisharp_code_actions)]]
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>os. <Plug>(omnisharp_code_action_repeat)]]
    --     vim.api.nvim_command [[autocmd FileType cs xmap <silent> <buffer> <Leader>os. <Plug>(omnisharp_code_action_repeat)]]
    --
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>os= <Plug>(omnisharp_code_format)]]
    --
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>osnm <Plug>(omnisharp_rename)]]
    --
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>osre <Plug>(omnisharp_restart_server)]]
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>osst <Plug>(omnisharp_start_server)]]
    --     vim.api.nvim_command [[autocmd FileType cs nmap <silent> <buffer> <Leader>ossp <Plug>(omnisharp_stop_server)]]
    --     vim.api.nvim_command [[augroup END]]
    vimp.add_buffer_maps(
      function()
        vimp.nmap({"silent"}, "gd", "<Plug>(omnisharp_go_to_definition)")
        vimp.nmap({"silent"}, "gr", "<Plug>(omnisharp_find_usages)")
        vimp.nmap({"silent"}, "<leader>gi", "<Plug>(omnisharp_find_implementations)")
        vimp.nmap({"silent"}, "gpd", "<Plug>(omnisharp_preview_definition)")
        vimp.nmap({"silent"}, "<leader>gpi", "<Plug>(omnisharp_preview_implementations)")
        vimp.nmap({"silent"}, "gtd", "<Plug>(omnisharp_type_lookup)")
        vimp.nmap({"silent"}, "<leader>ac", "<Plug>(omnisharp_global_code_check)")
      end
    )
  else
    buf_set_keymap("n", "gD", "<Cmd>lua vim.lsp.buf.declaration()<CR>", opts)
    buf_set_keymap("n", "gd", "<Cmd>lua vim.lsp.buf.definition()<CR>", opts)
    buf_set_keymap("n", "gtd", "<cmd>lua vim.lsp.buf.type_definition()<CR>", opts)
    buf_set_keymap("n", "<leader>gi", "<cmd>lua vim.lsp.buf.implementation()<CR>", opts)
    buf_set_keymap("n", "gr", "<cmd>lua vim.lsp.buf.references()<CR>", opts)
    buf_set_keymap("n", "rn", "<cmd>lua vim.lsp.buf.rename()<CR>", opts)
    buf_set_keymap("n", "<Leader>e", "<cmd>lua vim.lsp.diagnostic.show_line_diagnostics()<CR>", opts)
    buf_set_keymap("n", "<Leader>l", "<cmd>lua vim.lsp.diagnostic.set_loclist()<CR>", opts)

    -- LspSaga
    buf_set_keymap("n", "<leader>gh", ":Lspsaga lsp_finder<CR>", opts)
    buf_set_keymap("n", "K", ":Lspsaga hover_doc<CR>", opts)
    buf_set_keymap("n", "[g", ":Lspsaga diagnostic_jump_prev<CR>", opts)
    buf_set_keymap("n", "]g", ":Lspsaga diagnostic_jump_next<CR>", opts)
    buf_set_keymap("n", "<leader>ca", ":Lspsaga code_action<CR>", opts)
    buf_set_keymap("n", "<C-f>", "<cmd>lua require('lspsaga.action').smart_scroll_with_saga(1)<CR>", opts)
    buf_set_keymap("n", "<C-b>", "<cmd>lua require('lspsaga.action').smart_scroll_with_saga(-1)<CR>", opts)
  end
  if client.resolved_capabilities.document_formatting then
    buf_set_keymap("n", "<Leader>f", "<cmd>lua vim.lsp.buf.formatting()<CR>", opts)
  end
  if client.resolved_capabilities.document_range_formatting then
    buf_set_keymap("v", "<Leader>f", "<cmd>lua vim.lsp.buf.range_formatting()<CR>", opts)
  end

  -- Set autocommands conditional on server_capabilities
  if client.resolved_capabilities.document_highlight then
    vim.api.nvim_exec(
      [[
        hi LspReferenceRead cterm=bold ctermbg=red guibg=#464646
        hi LspReferenceText cterm=bold ctermbg=red guibg=#464646
        hi LspReferenceWrite cterm=bold ctermbg=red guibg=#464646
      augroup lsp_document_highlight
        autocmd! * <buffer>
        autocmd CursorHold <buffer> lua vim.lsp.buf.document_highlight()
        autocmd CursorMoved <buffer> lua vim.lsp.buf.clear_references()
      augroup END    ]],
      false
    )
  end
end

return custom_attach
