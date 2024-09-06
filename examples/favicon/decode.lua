function help(info)
    return [[
Usage:
 >  favicon/decode <inputPath>
    * Decodes all images in the favicon and saves them as PNG files.
    * The images are saved in the same directory as the input file, with the same name but the size appended.
    ]]
end

function init(workflow)
    workflow.import({
        "cmd",
        "ref",
        "cli",
        "io",
        "image",
    })
end

function main()
    local inRef = cmd.arg_string_pos()
    local ok, err = cmd.parse()

    if not ok then
        cli.print(cli.RED..err..cli.RESET)
        return
    end

    local inPath = ref.get(inRef)
    local cfg = io.decode_favicon_config(inPath)
    local imgs = nil

    if cfg.type == io.ICOTYPE_ICO then
        imgs = io.decode_favicon(inPath, image.ENCODING_PNG)
    elseif cfg.type == io.ICOTYPE_CUR then
        imgs = io.decode_favicon_cursor(inPath, image.ENCODING_PNG)
    end

    for i, img in ipairs(imgs) do
        io.encode(img, io.path_to(inPath))
    end
end