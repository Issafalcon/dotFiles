return {
	s("req", fmt("local {} = require('{}')", { i(1, "default"), rep(1) })),

	s({ trig = "for", dscr = "For loops in Lua" }, {
		t("for "),
		c(1, {
			sn(nil, {
				i(1, "k"),
				t(", "),
				i(2, "v"),
				t(" in "),
				c(3, { t("pairs"), t("ipairs") }),
				t("("),
				i(4),
				t(")"),
			}),
			sn(nil, { i(1, "i"), t(" = "), i(2), t(", "), i(3) }),
		}),
		t({ " do", "\t" }),
		i(0),
		t({ "", "end" }),
	}),

	s(
		{ trig = "plug", dscr = "Plugin config boilerplate" },
		fmt(
			[[
local status_ok, {} = pcall(require, "{}")
if not {} status_ok then
  return
end

local maps = require("custom_config").mappings
local mapper = require("utils.mapper")
local opts = {{ silent = true }}

{}.setup({{
  {}
}})
    ]],
			{ i(1), rep(1), rep(1), rep(1), i(0) }
		)
	),
}
