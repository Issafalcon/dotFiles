local ls = require("luasnip")
-- some shorthands...

-- Snippet Creator
-- s(<trigger>, <nodes>)
local s = ls.snippet
local sn = ls.snippet_node
local t = ls.text_node

-- Insert Node
-- Takes position (like $1) and optionally some default text
-- i(<position>, [default_text])
local i = ls.insert_node

local f = ls.function_node
local c = ls.choice_node
local d = ls.dynamic_node
local r = ls.restore_node
local l = require("luasnip.extras").lambda

-- Repeats a Node
-- rep(<position>)
local rep = require("luasnip.extras").rep
local p = require("luasnip.extras").partial
local m = require("luasnip.extras").match
local n = require("luasnip.extras").nonempty
local dl = require("luasnip.extras").dynamic_lambda

-- Format Node
-- Takes a format string and a list of nodes
-- fmt(<fmt_string>, {...nodes})
local fmt = require("luasnip.extras.fmt").fmt
local fmta = require("luasnip.extras.fmt").fmta
local types = require("luasnip.util.types")
local conds = require("luasnip.extras.expand_conditions")

ls.snippets = {
  lua = {
    s("req",
      fmt("local {} = require('{}')",
        { i(1, "default"), rep(1) })
    )
  }
}
