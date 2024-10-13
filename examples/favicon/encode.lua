---@diagnostic disable:lowercase-global

function help(info)
	return [[
Usage:
 >  favicon/encode <inputPath> [-c]
    * Encodes an image as a favicon. Will automatically create 16x16, 32x32, and 48x48 images.
    * Optionally, the -c flag can be used to encode the image as a cursor.
    * For cursors the hotspot will be set to the center of the image.
    ]]
end

function init(workflow)
	workflow.import({
		"cmd",
		"ref",
		"cli",
		"io",
		"image",
		"filter",
	})
end

function main()
	local inRef = cmd.arg_string_pos()
	local asCur = cmd.arg_flag("c", "cursor")
	local ok, err = cmd.parse()

	if not ok then
		cli.print(cli.RED .. err .. cli.RESET)
		return
	end

	local inPath = ref.get(inRef)
	local inImg = io.decode(inPath)

	local size1 = image.new("16x16", image.ENCODING_PNG, 16, 16, image.MODEL_NRGBA)
	local size2 = image.new("32x32", image.ENCODING_PNG, 32, 32, image.MODEL_NRGBA)
	local size3 = image.new("48x48", image.ENCODING_PNG, 48, 48, image.MODEL_NRGBA)

	filter.draw(inImg, size1, {
		filter.resize(16, 16, filter.RESAMPLING_NEARESTNEIGHBOR),
	})
	filter.draw(inImg, size2, {
		filter.resize(32, 32, filter.RESAMPLING_NEARESTNEIGHBOR),
	})
	filter.draw(inImg, size3, {
		filter.resize(48, 48, filter.RESAMPLING_NEARESTNEIGHBOR),
	})

	local outName = io.base(inPath)
	local outPath = io.path_to(inPath)

	if ref.get(asCur) then
		_ = io.encode_favicon_cursor(outName, { size1, size2, size3 }, { 8, 8, 16, 16, 24, 24 }, outPath)
	else
		_ = io.encode_favicon(outName, { size1, size2, size3 }, outPath)
	end
end

