config({
    name="GUI Test",
    version="0.1.0",
    author="Blub",
    requires={
        "std",
        "gui",
        "time",
    },

    desc="test imgui"
})

main(function ()
    local win = gui.window("test", 512, 512)

    local infloat = std.ref(
        5.0,
        std.REFTYPE_FLOAT32
    )

    local inInt = std.ref(
        5,
        std.REFTYPE_INT32
    )

    local inString = std.ref(
        "",
        std.REFTYPE_STRING
    )

    local inString2 = std.ref(
        "",
        std.REFTYPE_STRING
    )

    gui.window_run(win, function()
        gui.wg_single_window({
            gui.wg_input_float(infloat)
                :label("input float")
                :size(100)
                :step_size(0.1)
                :format("%.1f"),

             gui.wg_input_int(inInt)
                :label("input int")
                :size(100)
                :step_size(1),

            gui.wg_input_text(inString)
                :label("input text")
                :hint("text here")
                :size(100)
                :autocomplete({
                    "text 1",
                    "text 2",
                    "text 3",
                    "text 4",
                    "text 5",
                }),

            gui.wg_input_text_multiline(inString2)
                :label("input text multiline")
                :size(100, 100)
                :autoscroll_to_bottom(false),

            gui.wg_progress_bar(0.5)
                :overlay("half")
                :size(500, 20),
        })
    end)
end)
