require('auto-session').setup {
    post_save_cmds = {"VimspectorMkSession"},
    post_restore_cmds = {"VimspectorLoadSession .vimspector.session"}
}
