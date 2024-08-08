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

function barPlots()
    return gui.wg_plot("test bar plots"):plots({
        gui.pt_bar("bar", {
            1,
            2,
            3,
            4,
            5,
        }),
        gui.pt_bar_h("bar h", {
            1,
            2,
            3,
            4,
            5,
        })
    })
end

function linePlots()
    return gui.wg_plot("test line plots"):plots({
        gui.pt_line("line", {
            6,
            5,
            10,
            7,
            1.5,
        }),
        gui.pt_line("line 2", {
            6,
            5,
            10,
            7,
            1.5,
        }):x0(1),
        gui.pt_line_xy("line xy", {
            1,
            2,
            3,
            4,
            5,
            6
        }, {
            3,
            5,
            2,
            7,
            8,
            1,
        }),
    })
end

function piePlots()
    return gui.wg_plot("test pie charts")
        :flags(gui.FLAGPLOT_EQUAL)
        :x_axeflags(gui.FLAGPLOTAXIS_NODECORATIONS)
        :y_axeflags(gui.FLAGPLOTAXIS_NODECORATIONS, 0, 0)
        :axis_limits(0, 1, 0, 1, gui.COND_ALWAYS)
        :plots({
            gui.pt_pie_chart({
                "test 1",
                "test 2",
                "test 3",
            }, {
                0.2,
                0.2,
                0.6,
            }, 0.5, 0.5, 0.4)
        })
end

function scatterPlots()
    return gui.wg_plot("test scatter plots"):plots({
        gui.pt_scatter("scatter", {
            6,
            5,
            10,
            7,
            1.5,
        }),
        gui.pt_scatter_xy("scatter xy", {
            1,
            2,
            3,
            4,
            5,
            6
        }, {
            3,
            5,
            2,
            7,
            8,
            1,
        }),
    })
end

main(function ()
    local masterFlags = bit.bitor_many({
        gui.FLAGMASTERWINDOW_FLOATING,
    })

    local win = gui.window_master("test", 1024, 1024, masterFlags)

    local splitref1 = std.ref(512, std.REFTYPE_FLOAT32)
    local splitref2 = std.ref(512, std.REFTYPE_FLOAT32)

    gui.window_run(win, function()
        gui.window_single():layout({
            gui.wg_layout_split(gui.SPLITDIRECTION_VERTICAL, splitref1, {
                gui.wg_layout_split(gui.SPLITDIRECTION_HORIZONTAL, splitref2, {
                    barPlots(),
                }, {
                    linePlots(),
                })
            }, {
                gui.wg_layout_split(gui.SPLITDIRECTION_HORIZONTAL, splitref2, {
                    piePlots(),
                }, {
                    scatterPlots(),
                })
            }),
        })
    end)
end)