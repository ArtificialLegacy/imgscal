
config({
    name= "Context",
    version= "0.1.0",
    author="Blub",
    requires= {
        "io",
        "context"
    },

    desc="Context lib testing.",
})

main(function ()
    ctx = context.new(64, 64)
    img = context.to_image(ctx)
    io.out(img, "./output")
end)