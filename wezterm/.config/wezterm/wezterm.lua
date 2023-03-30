local wezterm = require("wezterm")
local colors = require("colors")
local fonts = require("fonts")

local config = {}

if wezterm.config_builder then
    config = wezterm.config_builder()
end

-- Custom Config
colors.apply_to_config(config)
fonts.apply_to_config(config)

return config
