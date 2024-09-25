function init(workflow)
	workflow.import({
		"tui",
		"lipgloss",
	})
end

function main()
	local program = tui.new()
		:init(function(id)
			local model = {
				title = "testing program",
				list = tui.list(id, {
					tui.list_item("test 1", "one"),
					tui.list_item("test 2", "two"),
					tui.list_item("test 3", "three"),
				}, 250, 10, tui.list_delegate(id):show_description_set(false))
					:help_show_set(true)
					:title_show_set(true)
					:statusbar_show_set(true)
					:style_titlebar_set(lipgloss.style():border_set(lipgloss.border_double(), true)),
				keymap = {
					A = tui.keybinding(id, tui.keybinding_option():keys("a", "A"):help("?", "a or A.")),
					B = tui.keybinding(id, tui.keybinding_option():keys("b", "B"):help("!", "b or B")),
				},
			}

			model.list:help_short_additional(function()
				return {
					model.keymap.A,
					model.keymap.B,
				}
			end)

			model.list:help_full_additional(function()
				return {
					model.keymap.A,
					model.keymap.B,
				}
			end)

			return model, tui.cmd_none()
		end)
		:update(function(model, msg)
			return model.list.update()
		end)
		:view(function(model)
			return model.title .. "\n\n" .. model.list.view() .. "\n\n"
		end)

	tui.run(program, tui.program_options():fps(1))
end
