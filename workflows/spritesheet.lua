
config({
    name= "Spritesheet",
    version= "0.1.0",
    author="Blub",
    requires= {
        "std",
        "io",
        "cli",
        "spritesheet",
        "collection",
        "image"
    },

    desc="Splits a spritesheet up",
})

main(function ()
    pth = cli.question("Enter image to process")
    img = io.load_image(pth)

    subimgs = spritesheet.to_frames(img, "frame", 8, 8, 8, 8)

    for v,k in pairs(subimgs) do
        image.map(k, function (x, y, c)
            return {red=0, green=0, blue=0, alpha=0}
        end)

        io.out(k, "./output")
    end

    ss = spritesheet.from_frames(subimgs, "ss", 8, 8, image.MODEL_RGBA, image.ENCODING_PNG, 2, nil, 2, 2)
    io.out(ss, "./output")
end)