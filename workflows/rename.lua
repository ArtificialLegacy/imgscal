
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
    img1 = imgscal.prompt_file("Enter file to rename")
    img2 = imgscal.prompt_file("Enter file to rename")
    imgscal.name(img1, "output_file1.png")
    imgscal.name(img2, "output_file2.png")
    imgscal.out(img1, "./output")
    imgscal.out(img2, "./output")
end)