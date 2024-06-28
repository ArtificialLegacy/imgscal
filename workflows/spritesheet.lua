
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

    subimgs = spritesheet.to_frames(img, "frame", 8, 160, 160, 8)

    for v,k in pairs(subimgs) do
        io.out(v, "./output")
    end

    ss = spritesheet.from_frames(subimgs, "ss", 160, 160, image.MODEL_NRGBA, image.ENCODING_PNG)
    io.out(ss, "./output")
end)