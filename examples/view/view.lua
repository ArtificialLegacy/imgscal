function help(info)
	return [[
Usage:
 >  view <inputPath>
    * Displays the image in the terminal. 
    ]]
end

function init(workflow)
	workflow.import({
		"cmd",
		"ref",
		"cli",
		"io",
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

	cli.clear()

	local img = io.decode(inPath)
	cli.print_image(img, true)

	cli.println()
end
