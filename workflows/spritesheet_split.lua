
config({
    name= "Spritesheet Split",
    version= "0.1.0",
    author="Blub",
    requires= {
        "imgscal",
        "imgscal_sheet"
    },

    desc="Extracts an animation spritesheet from a larger spritesheet.",
})

main(function ()
    img = imgscal.prompt_file()

    walksheet = imgscal_sheet.crop_frames(img, 128, 128, 5, 0)
    imgscal.name(walksheet, "sprWalk")

    imgscal.out(walksheet, "./output")
end)