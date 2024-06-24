
config({
    name= "Spritesheet",
    version= "0.1.0",
    author="Blub",
    requires= {
        "io",
        "cli",
        "spritesheet"
    },

    desc="Splits a spritesheet up",
})

main(function ()
    pth = cli.question("Enter image to process")
    img = io.load_image(pth)

    subimgs = spritesheet.to_frames(img, "frame", 8, 160, 160, 8)

    for k,v in pairs(subimgs) do
        io.out(v, "./output")
    end
end)