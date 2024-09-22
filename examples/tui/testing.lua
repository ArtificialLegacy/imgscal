function init(workflow)
	workflow.import({
		"tui",
	})
end

function main()
	local program = tui.new()
		:init(function(id)
			local model = {
				title = "testing program",
				viewport = tui.viewport(id, 100, 1):content_set("test 1\ntest 2\ntest 3\ntest 4"),
			}

			local vkm = model.viewport.keymap()
			vkm.down:enabled_set(false)

			return model, tui.cmd_none()
		end)
		:update(function(model, msg)
			return model.viewport.update()
		end)
		:view(function(model)
			return model.title .. "\n\n" .. model.viewport.view() .. "\n\n"
		end)

	tui.run(program)
end
