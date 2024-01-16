
config({
    name= "ESRGAN_AnimeX4",
    version= "1.0.0",
    requires= {
        "imgscal",
        "esrgan",
    }
})

main(function (file)
    job("esrgan.animex4", file)
    file = job("imgscal.rename", file, {prefix= "up_"})
    job("imgscal.output", file)
end)
