---@param workflow imgscal_WorkflowInfo
function help(info)
	return [[
Usage:
 >  resize <inputPath> <width> <height> [-o=outputPath] [-r=resampling]
    * If output is omitted, it will append 'resized_' to the inputted file name.
    * If resampling is omitted, it will default to 'box'.
        * Valid resampling values are: ['box', 'cubic', 'lanczos', 'linear', 'nn'].
        * 'nn' is shorthand for nearest neighbor.
    ]]
end

---@param workflow imgscal_WorkflowInit
function init(workflow)
	workflow.import({
		"cmd",
		"ref",
		"cli",
		"filter",
		"io",
		"image",
		"std",
	})
end

function main()
	local inRef = cmd.arg_string_pos()
	local widthRef = cmd.arg_int_pos()
	local heightRef = cmd.arg_int_pos()
	local outRef = cmd.arg_string("o", "output")
	local rRef = cmd.arg_selector("r", "resampling", { "box", "cubic", "lanczos", "linear", "nn" })

	local ok, err = cmd.parse()

	if not ok then
		cli.print(cli.RED .. err .. cli.RESET)
		return
	end

	local width = ref.get(widthRef)
	local height = ref.get(heightRef)

	if width == 0 then
		std.panic("image width must not be 0")
	end
	if height == 0 then
		std.panic("image height must not be 0")
	end

	local inPath = ref.get(inRef)
	local inImg = io.decode(inPath)

	-- check output file name, if not provided default to input file with "resized_" prefix.
	local outPath = ref.get(outRef)
	local outName = ""
	if outPath == "" then
		outPath = inPath
		outName = "resized_" .. io.base(outPath)
	else
		outName = io.base(outPath)
	end

	local outImg = image.new(outName, image.path_to_encoding(outPath), width, height)

	-- get resampling method to use, defaulting to box.
	local r = ref.get(rRef)
	local resampling = filter.RESAMPLING_BOX

	if r == "cubic" then
		resampling = filter.RESAMPLING_CUBIC
	elseif r == "lanczos" then
		resampling = filter.RESAMPLING_LANCZOS
	elseif r == "linear" then
		resampling = filter.RESAMPLING_LINEAR
	elseif r == "nn" then
		resampling = filter.RESAMPLING_NEARESTNEIGHBOR
	end

	filter.draw(inImg, outImg, {
		filter.resize(width, height, resampling),
	})

	io.encode(outImg, io.path_to(outPath))
end

