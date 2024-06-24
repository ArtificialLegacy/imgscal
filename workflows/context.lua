
config({
    name= "Context",
    version= "0.1.0",
    author="Blub",
    requires= {
        "io",
        "cli",
        "image",
        "collection",
        "context"
    },

    desc="Context lib testing.",
})

main(function ()
    pth = cli.question("Enter directory to process")
    pthList = io.dir_img(pth)

    for k,v in pairs(pthList) do
        local img = io.load_image(v)
        local ctx = context.new_image(img)
        context.color_hex(ctx, "#000000EE")

        context.draw_polygon(ctx, 3, 150, 220, 180, 0)
        context.dash_set(ctx, {20, 20, 15, 30})
        context.line_width(ctx, 10)
        context.stroke(ctx)

        local imgAfter = context.to_image(ctx, "output_"..k, image.ENCODING_PNG)
        io.out(imgAfter, "./output")
    end
end)