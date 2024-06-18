
config({
    name= "Rename",
    version= "0.1.0",
    author="Blub",
    requires= {
        "io",
        "cli",
        "image",
        "collection"
    },

    desc="Renames all files in a directory and saves them to ./output",
})

main(function ()
    pth = cli.question("Enter directory to process")
    pthList = io.dir_img(pth)

    for k,v in pairs(pthList) do
        local img = io.load_image(v)
        local width, height = image.size(img)
        image.name_ext(img, {prefix=width.."_"..height, ext=".png"})
        io.out(img, "./output")
    end
end)