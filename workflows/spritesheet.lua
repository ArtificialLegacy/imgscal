
config({
    name= "Spritesheet",
    version= "0.1.0",
    author="Blub",
    requires= {
        "std",
        "io",
        "cli",
        "spritesheet",
        "collection"
    },

    desc="Splits a spritesheet up",
})

main(function ()
    pth = cli.question("Enter image to process")
    img = io.load_image(pth)

    task = collection.task("testing")

    collection.schedule(collection.TYPE_TASK, task, false, function ()
        std.log("task run")
    end)

    subimgs = spritesheet.to_frames(img, "frame", 8, 160, 160, 8)

    for k,v in pairs(subimgs) do
        io.out(v, "./output")
    end
end)