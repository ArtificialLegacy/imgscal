
config({
    name= "ASCII",
    version= "0.1.0",
    author="Blub",
    requires= {
        "io",
        "cli",
        "ascii"
    },

    desc="Converts all images in a dir to ascii art.",
})

main(function ()
    pth = cli.question("Enter directory to process")
    pthList = io.dir_img(pth)

    for k,v in pairs(pthList) do
        local img = io.load_image(v)
        ascii.to_file_size(img, "./output/"..k..".txt", 64, 64, false, false)
    end
end)