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
				progress = tui.progress(id),
			}

			return model, model.progress.percent_set(0.5)
		end)
		:update(function(model, msg)
			return model.progress.update()
		end)
		:view(function(model)
			return model.title .. "\n\n" .. model.progress.view() .. "\n\n"
		end)

	tui.run(program)
end
