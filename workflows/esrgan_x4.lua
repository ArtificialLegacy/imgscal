
config({
    name= "ESRGAN_X4",
    version= "1.0.0",
    requires= {
        "imgscal",
        "esrgan",
    }
})

main(function (file)
    job("esrgan.x4", file)
    file = job("imgscal.rename", file, {prefix= "up_"})
    job("imgscal.output", file)
end)
