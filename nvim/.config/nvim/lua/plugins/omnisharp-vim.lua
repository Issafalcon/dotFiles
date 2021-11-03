require("vimp")

vim.g.OmniSharp_highlight_groups = {
    ClassName = 'ClassName',
    StaticSymbol = 'StaticSymbol',
    DelegateName = 'DelegateName',
    EnumName = 'EnumName',
    EnumMemberName = 'EnumMemberName',
    ConstantName = 'ConstantName',
    PropertyName = 'PropertyName',
    InterfaceName = 'Include',
    StructName = 'StructName',
    FieldName = 'FieldName',
    ParameterName = 'ParameterName',
    LocalName = 'LocalName',
    MethodName = 'MethodName',
    NamespaceName = 'NamespaceName',
    ExtensionMethodName = 'ExtensionMethodName',
    EventName = 'EventName'
}

-- Mappings for omnisharp only functions (functions not available in LSP)
vimp.nmap({"silent"}, "<leader>odt", "<Plug>(omnisharp_debug_test)")
vimp.nmap({"silent"}, "<leader>ort", "<Plug>(omnisharp_run_test)")
vimp.nmap({"silent"}, "<leader>ot", "<Plug>(omnisharp_run_tests_in_file)")
