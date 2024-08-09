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

    local win = gui.window_master("test", 1024, 1024, masterFlags)

    gui.window_run(win, function()
        gui.window_single():layout({
            gui.wg_label(std.wd()),
            gui.wg_label(std.fmt("test fmt: %.1f", {5})),
        })
    end)
end)