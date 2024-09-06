
function help(info)
    return [[
Usage:
 >  favicon <inputPath>
    * Decodes all images in the favicon and prints them to the console.
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

    cli.clear()

    cli.printlnf("Decoded Favicon: %s%s%s", cli.CYAN, inPath, cli.RESET)
    if cfg.type == io.ICOTYPE_ICO then
        cli.printlnf("Type: %sICO%s", cli.YELLOW, cli.RESET)
    elseif cfg.type == io.ICOTYPE_CUR then
        cli.printlnf("Type: %sCUR%s", cli.MAGENTA, cli.RESET)
    end

    cli.printlnf("Image Count: %s%d%s (%sLargest: %d%s)", cli.GREEN, cfg.count, cli.RESET, cli.GREEN, cfg.largest+1, cli.RESET)

    cli.println()

    for i, img in ipairs(imgs) do
        cli.printf("Image %d: %s%dx%d%s", i, cli.GREEN, cfg.entries[i].width, cfg.entries[i].height, cli.RESET)
        if i == cfg.largest+1 then
            cli.printf(" (%sLargest%s)", cli.GREEN, cli.RESET)
        end
        cli.println()

        if cfg.type == io.ICOTYPE_CUR then
            cli.printlnf("Hotspot: (%s%d,%d%s)", cli.GREEN, cfg.entries[i].data1, cfg.entries[i].data2, cli.RESET)
        end

        cli.print_image(img, true)
        cli.println()
    end
end