
config({
    name= "Context",
    version= "0.1.0",
    author="Blub",
    requires= {
        "io",
        "context",
        "std"
    },

    desc="Context lib testing.",
})

main(function ()
    ctx = context.new(64, 64)
    context.color_hex(ctx, "#FFF")
    context.draw_circle(ctx, 32, 32, 16)
    context.stroke(ctx)
    img = context.to_image(ctx)
    io.out(img, "./output")
end)