
config({
    name= "Rename",
    version= "0.1.0",
    author="Blub",
    requires= {
        "imgscal",
    },

    desc="Renames a file and saves it to ./output",
})

main(function ()
    img = imgscal.prompt_file("Enter file to rename")
    imgscal.name(img, "output_file.png")
    imgscal.out(img, "./output")
end)