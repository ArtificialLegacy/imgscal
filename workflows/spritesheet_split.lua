
config({
    name= "Spritesheet Split",
    version= "0.1.0",
    author="Blub",
    requires= {
        "imgscal",
        "imgscal_sheet"
    }
})

main(function ()
    img = imgscal.prompt_file()

    frames = imgscal_sheet.to_frames(img, 128, 128, 5, 0)
    walksheet = imgscal_sheet.to_sheet(frames)
    imgscal.name(walksheet, "sprWalk")

    imgscal.out(walksheet, "./output")
end)