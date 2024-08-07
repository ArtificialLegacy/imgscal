config({
    name="GUI Test",
    version="0.1.0",
    author="Blub",
    requires={
        "gui",
        "bit",
        "std",
    },

    desc="test imgui"
})

main(function ()
    local masterFlags = bit.bitor_many({
        gui.FLAGMASTERWINDOW_FLOATING,
    })
    
    local win = gui.window_master("test", 512, 512, masterFlags)

    local fref = std.ref(256, std.REFTYPE_FLOAT32)

    gui.window_run(win, function()
        gui.window_single():layout({
            gui.wg_layout_split(
                gui.SPLITDIRECTION_VERTICAL,
                fref,
                {
                    gui.wg_label("layout1")
                },
                {
                    gui.wg_label("layout2")
                }
            ):border(false)
        })

        
    end)
end)
