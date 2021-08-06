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
    vimp.add_buffer_maps(
      function()
        vimp.nmap({"silent"}, "gd", "<Plug>(omnisharp_go_to_definition)")
        vimp.nmap({"silent"}, "gr", "<Plug>(omnisharp_find_usages)")
        vimp.nmap({"silent"}, "<leader>gi", "<Plug>(omnisharp_find_implementations)")
        vimp.nmap({"silent"}, "gpd", "<Plug>(omnisharp_preview_definition)")
        vimp.nmap({"silent"}, "<leader>gpi", "<Plug>(omnisharp_preview_implementations)")
        vimp.nmap({"silent"}, "gtd", "<Plug>(omnisharp_type_lookup)")
        vimp.nmap({"silent"}, "<leader>ac", "<Plug>(omnisharp_global_code_check)")
        vimp.nmap({"silent"}, "K", "<Plug>(omnisharp_documentation)")
        vimp.nmap({"silent"}, "rn", "<Plug>(omnisharp_rename)")
        vimp.nmap({"silent"}, "<leader>f", "<Plug>(omnisharp_code_format)")
        vimp.nmap({"silent"}, "gm", "<Plug>(omnisharp_find_members)")
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

    buf_set_keymap("n", "K", ":Lspsaga hover_doc<CR>", opts)

    if client.resolved_capabilities.document_formatting then
      buf_set_keymap("n", "<Leader>f", "<cmd>lua vim.lsp.buf.formatting()<CR>", opts)
    end

    if client.resolved_capabilities.document_range_formatting then
      buf_set_keymap("v", "<Leader>f", "<cmd>lua vim.lsp.buf.range_formatting()<CR>", opts)
    end
  end

  buf_set_keymap("n", "<A-s>", "<cmd>lua vim.lsp.buf.signature_help()<CR>", opts)

  -- LspSaga
  buf_set_keymap("n", "<leader>gh", ":Lspsaga lsp_finder<CR>", opts)
  buf_set_keymap("n", "[g", ":Lspsaga diagnostic_jump_prev<CR>", opts)
  buf_set_keymap("n", "]g", ":Lspsaga diagnostic_jump_next<CR>", opts)
  buf_set_keymap("n", "<leader>ca", ":Lspsaga code_action<CR>", opts)
  buf_set_keymap("n", "<C-f>", "<cmd>lua require('lspsaga.action').smart_scroll_with_saga(1)<CR>", opts)
  buf_set_keymap("n", "<C-b>", "<cmd>lua require('lspsaga.action').smart_scroll_with_saga(-1)<CR>", opts)

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

  require("lsp_signature").on_attach(
    {
      bind = true,
      use_lspsaga = false,
      fix_pos = false
    }
  )
end

return custom_attach
