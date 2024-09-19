function init(workflow)
	workflow.import({
		"tui",
		"io",
	})
end

function main()
	local program = tui.new()
		:init(function(id)
			local model = {
				title = "testing program",
				spinner = tui.spinner(id, tui.SPINNER_MINIDOT),
				filepicker = tui.filepicker(id):current_directory_set(io.wd()),
			}

			return model, tui.cmd_batch({
				model.spinner.tick(),
				model.filepicker.init(),
			})
		end)
		:update(function(model, msg)
			if msg.msg == tui.MSG_SPINNERTICK then
				if model.spinner.id == msg.id then
					return model.spinner.update()
				end
			end
			return model.filepicker.update()
		end)
		:view(function(model)
			return model.title .. "\n\n" .. model.spinner.view() .. "\n\n" .. model.filepicker.view() .. "\n\n"
		end)

	tui.run(program)
end
