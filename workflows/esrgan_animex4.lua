
config({
    name= "esrgan_animex4",
    version= "0.1.0",
    requires= {
        "esrgan",
    }
})

function main (file)
    process("esrgan.animex4", file)
    process("core.output", file, "up_")
end