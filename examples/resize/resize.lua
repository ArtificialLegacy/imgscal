
function init(workflow)
    workflow.import({
        "cli",
        "filter",
        "io",
        "image",
    })
end

function main()
    local inpath = cli.question("Enter input image to resize")
    local inImg = io.decode(inpath)

    local width = tonumber(cli.question("Enter new width"))
    local height = tonumber(cli.question("Enter new height"))

    local outName = "resized_"..io.base(inpath)
    local outImg = image.new(outName, image.path_to_encoding(inpath), width, height)

    local resampling = cli.select("Select resampling method to use", {
        "Box",
        "Cubic",
        "Lanczos",
        "Linear",
        "Nearest Neighbor",
    })

    filter.draw(inImg, outImg, {
        filter.resize(width, height, resampling-1),
    })

    io.encode(outImg, io.path_to(inpath))
    cli.print("Resized image saved to: "..io.path_join(inpath, outName)..".")
end