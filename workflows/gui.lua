config({
    name="GUI Test",
    version="0.1.0",
    author="Blub",
    requires={
        "gui",
    },

    desc="test imgui"
})

main(function ()
    local win = gui.window("test", 512, 512)

    local labelText = "test"

    gui.window_run(win, function()
        gui.wg_single_window({
            gui.wg_label(labelText),
            gui.wg_button("test 2")
                :size(50, 50)
                :on_click(function()
                    labelText = "button pressed"
                end),
        })
    end)
end)
