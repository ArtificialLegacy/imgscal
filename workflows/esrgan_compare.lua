
config({
    name= "ESRGAN_COMPARE",
    version= "1.0.0",
    requires= {
        "imgscal",
        "esrgan",
    }
})

main(function (file)
    file1 = job('imgscal.copy', file, {prefix= "x4_2_"})
    file2 = job("imgscal.copy", file, {prefix= "animex4_2_"})
    file3 = job('imgscal.copy', file, {prefix= "x4_3_"})
    file4 = job("imgscal.copy", file, {prefix= "animex4_3_"})
    file5 = job('imgscal.copy', file, {prefix= "x4_4_"})
    file6 = job("imgscal.copy", file, {prefix= "animex4_4_"})

    job("esrgan.x4", file1, {scale= 2})
    job("esrgan.animex4", file2, {scale= 2})
    job("esrgan.x4", file3, {scale= 3})
    job("esrgan.animex4", file4, {scale= 3})
    job("esrgan.x4", file5, {scale= 4})
    job("esrgan.animex4", file6, {scale= 4})
    
    job("imgscal.output", file1)
    job("imgscal.output", file2)
    job("imgscal.output", file3)
    job("imgscal.output", file4)
    job("imgscal.output", file5)
    job("imgscal.output", file6)
end)
