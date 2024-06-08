
config({
    name= "Spritesheet Split",
    version= "1.0.0",
    author="Blub",
    requires= {
        "imgscal",
    }
})

main(function ()
    dir = imgscal.prompt_dir()
    imgs = imgscal.load_dir(dir)

    frames = imgscal.to_frames(imgs[0], 128, 128, 5, 0)
    walksheet = imgscal.to_sheet(frames)
    imgscal.name(walksheet, "sprWalk")

    imgscal.out(walksheet, "./output")
end)