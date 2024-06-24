
config({
    name= "QRCode",
    version= "0.1.0",
    author="Blub",
    requires= {
        "qrcode",
        "io",
        "image"
    },

    desc="makes a qrcode",
})

main(function ()
    local qr = qrcode.new("blub", qrcode.RECOVERY_HIGHEST)
    qrcode.color_set_foreground(qr, {
        red=255,
        green=0,
        blue=0,
        alpha=255
    })
    qrcode.color_set_background(qr, {
        red=0,
        green=0,
        blue=255,
        alpha=255
    })
    qrcode.border_set(qr)
    local img = qrcode.to_image(qr, "qrcode", -1, image.ENCODING_PNG)
    io.out(img, "./output")
end)