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
    local win = gui.window_master("test", 512, 512)

    gui.window_run(win, function()
        gui.window_single({
            
        })
    end)
end)
