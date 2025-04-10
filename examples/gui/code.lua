---@param workflow imgscal_WorkflowInit
function init(workflow)
	workflow.import({
		"gui",
		"guicode",
		"cli",
	})
end

example = [[select *
     from greeting
     where date > current_timestamp
     order by date;]]

function main()
	local win = gui.window_master("Code Editor Example", 512, 512, 0)
	gui.window_set_icon_imgscal(win, true)

	local code =
		guicode.wg_code_editor():language_definition(guicode.LANG_SQL):border(true):text(example):show_whitespace(false)

	local font, fontok = gui.fontatlas_add_font("FiraMonoNerdFont-Regular", 18)

	local style = gui.wg_style()
	if fontok then
		style:font(font)
	end

	_ = gui.window_run(win, function()
		gui.window_single():layout({
			gui.wg_align(gui.ALIGN_CENTER):to({
				gui.wg_label("SQL Code Editor"),
			}),
			style:to({ code }),
		})
	end)
end
