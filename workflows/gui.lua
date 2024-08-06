config({
    name="GUI Test",
    version="0.1.0",
    author="Blub",
    requires={
        "std",
        "image",
        "gui",
    },

    desc="test imgui"
})

main(function ()
    local win = gui.window("test", 512, 512)

    local labelText = "test"
    local checked = std.ref(false, std.REFTYPE_BOOL)
    local clr = std.ref(image.color_rgb(255, 0, 0), std.REFTYPE_RGBA)

    gui.window_run(win, function()
        gui.wg_single_window({
            gui.wg_label(labelText),
            gui.wg_separator(),
            gui.wg_button("test 2")
                :size(50, 50)
                :on_click(function()
                    labelText = "button pressed"
                end),
            gui.wg_child()
                :border(true)
                :layout({
                    gui.wg_bullet_text("bullet"),
                    gui.wg_checkbox("check", checked),
                    gui.wg_color_edit("color", clr),
                }),
        })
    end)
end)
