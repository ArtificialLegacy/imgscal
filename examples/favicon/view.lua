---@param info imgscal_WorkflowInfo
function help(info)
	return [[
Usage:
 >  favicon <inputPath>
    * Decodes all images in the favicon and prints them to the console.
    ]]
end

---@param workflow imgscal_WorkflowInit
function init(workflow)
	workflow.import({
		"cmd",
		"ref",
		"cli",
		"io",
		"image",
		"tui",
		"lipgloss",
		"std",
	})
end

function main()
	local inRef = cmd.arg_string_pos()
	local ok, err = cmd.parse()

	if not ok then
		cli.print(cli.RED .. err .. cli.RESET)
		return
	end

	local inPath = ref.get(inRef)
	local cfg = io.decode_favicon_config(inPath)
	local imgs = {}

	if cfg.type == io.ICOTYPE_ICO then
		imgs = io.decode_favicon(inPath, image.ENCODING_PNG)
	elseif cfg.type == io.ICOTYPE_CUR then
		imgs = io.decode_favicon_cursor(inPath, image.ENCODING_PNG)
	end

	cli.clear()

	local imgStrings = {}

	for i, v in pairs(imgs) do
		imgStrings[i] = cli.string_image(v, true, 10)
	end

	local program = tui.new()
		:init(function(id)
			local model = {
				width = 100,
				height = 100,
				paginator = tui.paginator(id, 0, cfg.count)
					:type_set(tui.PAGINATOR_DOT)
					:format_dot_set(" ● ", " ◌ "),
				img = tui.image(id, true, true):image_string_set(imgStrings[1]),
				current = 0,
			}

			return model, tui.cmd_batch({
				tui.cmd_window_size(),
				tui.cmd_window_title("Favicon Viewer"),
			})
		end)
		:update(function(model, msg)
			if msg.msg == tui.MSG_WINDOWSIZE then
				---@cast msg tui.MSGWindowSize
				model.width = msg.width
				model.height = msg.height
			end

			local pagecmd = model.paginator.update()
			local imgcmd = model.img.update()

			local cmds = { pagecmd, imgcmd }

			local index = model.paginator.page() + 1
			if index ~= model.current then
				model.img:image_string_set(imgStrings[index])
				local imgsize = model.img.size_set(cfg.entries[index].width * 2 + 4, cfg.entries[index].height)
				model.current = index

				cmds[3] = imgsize
			end

			return tui.cmd_batch(cmds)
		end)
		:view(function(model)
			local index = model.paginator.page() + 1

			local titlestr = image_name(inPath, model)
			local cfgstr = favicon_data(cfg, model)
			local datastr = image_data(cfg.entries[index], model)
			local pagestr = page_view(model)
			local imgstr = image_view(model)

			return lipgloss.join_vertical(lipgloss.POSITION_CENTER, titlestr, cfgstr, datastr, imgstr, pagestr)
		end)

	tui.run(program)
end

function image_name(str, model)
	return lipgloss.style_string(
		str,
		lipgloss
			.style()
			:width_set(model.width)
			:align_horizontal_set(lipgloss.POSITION_CENTER)
			:bold_set(true)
			:background_set(lipgloss.color("#282A36"))
			:foreground_set(lipgloss.color("#BD93F9"))
			:padding_top_set(1)
			:padding_bottom_set(1)
	)
end

function favicon_data(cfg, model)
	local largest = cfg.entries[cfg.largest + 1]

	return lipgloss.style_string(
		std.fmt("\nImage Count: %d - Largest: (%dpx, %dpx)", cfg.count, largest.width, largest.height),
		lipgloss
			.style()
			:width_set(model.width)
			:align_horizontal_set(lipgloss.POSITION_CENTER)
			:foreground_set(lipgloss.color("#F8F8F2"))
			:background_set(lipgloss.color("#44475A"))
	)
end

function image_data(entry, model)
	return lipgloss.style_string(
		std.fmt("Current: (%dpx, %dpx)\n", entry.width, entry.height),
		lipgloss
			.style()
			:width_set(model.width)
			:align_horizontal_set(lipgloss.POSITION_CENTER)
			:foreground_set(lipgloss.color("#F8F8F2"))
			:background_set(lipgloss.color("#44475A"))
	)
end

function page_view(model)
	return lipgloss.style_string(
		model.paginator.view() .. "\n" .. "← Prev  ●  Next →",
		lipgloss
			.style()
			:width_set(model.width)
			:align_horizontal_set(lipgloss.POSITION_CENTER)
			:foreground_set(lipgloss.color("#F8F8F2"))
			:background_set(lipgloss.color("#44475A"))
			:bold_set(true)
	)
end

function image_view(model)
	return lipgloss.style_string(
		model.img.view(),
		lipgloss.style():height_set(model.height - 9):align_vertical_set(lipgloss.POSITION_CENTER)
	)
end
