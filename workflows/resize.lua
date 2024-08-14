config({
    name="Resize",
    version="1.0.0",
    author="Blub",
    requires={
        "cmd",
        "ref",
        "cli",
        "filter",
        "io",
        "image",
        "std",
    },
    cli_exclusive=true,

    desc="Resize an image from the command line.",
})

--[[
    Usage:
    resize <inputImage> <width> <height> [-o=outputImage{resized_..inputImage}] [-r=[box, cubic, lanczos, linear, nn]{box}]

    if -o is excluded, the output file will be the same as the input, but with the string "resized_" appended to the beginning.
    if -r is excluded, it will default to using box resampling.
    accepted -r values are box, cubic, lanczos, linear, and nn (nearest neighbor).
--]]

main(function ()
    local inRef = cmd.arg_string_pos()
    local widthRef = cmd.arg_int_pos()
    local heightRef = cmd.arg_int_pos()
    local outRef = cmd.arg_string("o",  "output")
    local rRef = cmd.arg_selector("r", "resampling", {"box", "cubic", "lanczos", "linear", "nn"})

    local ok, err = cmd.parse()

    if not ok then
        cli.print(cli.RED..err..cli.RESET)
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
    local inImg = io.load_image(inPath)

    -- check output file name, if not provided default to input file with "resized_" prefix.
    local outPath = ref.get(outRef)
    local outName = ""
    if outPath == "" then
        outPath = inPath
        outName = "resized_"..io.base(outPath)
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

    io.out(outImg, io.path_to(outPath))
end)