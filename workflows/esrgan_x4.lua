
config({
    name= "ESRGAN_X4",
    version= "0.1.0",
    requires= {
        "imgscal",
        "esrgan",
    }
})

main(function (file)
    job("esrgan.x4", file)
    job("imgscal.rename", file, {prefix= "x4_"})
    job("imgscal.output", file)
end)