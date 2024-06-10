
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
    pth1 = imgscal.prompt("Enter file to rename")
    pth2 = imgscal.prompt("Enter file to rename")

    img1 = imgscal.load_image(pth1)
    imgscal.name(img1, "output_file1.png")
    imgscal.out(img1, "./output")

    img2 = imgscal.load_image(pth2)
    imgscal.name(img2, "output_file2.png")
    imgscal.out(img2, "./output")
end)