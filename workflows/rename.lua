
config({
    name= "Rename",
    version= "0.1.0",
    author="Blub",
    requires= {
        "std",
        "io",
        "cli",
        "image"
    },

    desc="Renames all files in a directory and saves them to ./output",
})

main(function ()
    pth = cli.question("Enter directory to process")
    pthList = io.dir_img(pth)

    for k,v in pairs(pthList) do
        local img = io.load_image(v)
        image.name_ext(img, {prefix="proc_", ext=".png"})
        io.out(img, "./output")
        image.collect(img)
    end
end)