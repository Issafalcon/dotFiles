local wezterm = require("wezterm")
local fonts = {}

function fonts.apply_to_config(config)
    config.font = wezterm.font("JetBrains Mono")
end

return fonts
