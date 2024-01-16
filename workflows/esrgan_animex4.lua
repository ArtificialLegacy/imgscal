
config({
    name= "ESRGAN_AnimeX4",
    version= "0.1.0",
    requires= {
        "imgscal",
        "esrgan",
    }
})

main(function (file)
    job("esrgan.animex4", file)
    job("imgscal.rename", file, {prefix= "animex4_"})
    job("imgscal.output", file)
end)