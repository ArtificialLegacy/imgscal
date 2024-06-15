
config({
    name= "ASCII",
    version= "0.1.0",
    author="Blub",
    requires= {
        "io",
        "cli",
        "ascii",
        "txt",
    },

    desc="Converts all images in a dir to ascii art.",
})

main(function ()
    pth = cli.question("Enter directory to process")
    pthList = io.dir_img(pth)

    for k,v in pairs(pthList) do
        local img = io.load_image(v)
        local str = ascii.to_string_size(img, 64, 64, false, false)
        local file = txt.file_open("./output", k..".txt")
        txt.write(file, str)
    end
end)