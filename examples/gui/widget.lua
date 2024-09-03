
local widget = {}

--[[
    Helper functions for creating the sliders.
    These use a dummy widget that gurantees the width of the label+dummy is 100 pixels,
    this is for visuals to make sure the sliders are vertically aligned.

    It also adds an arrow button widget that resets the slider's ref to a default value.
--]]

function widget.slider_float(str, rf, min, max, dflt)
    return gui.wg_row({
        gui.wg_label(str),
        gui.wg_dummy(
            100-gui.calc_text_size_width(str),
            1
        ),
        gui.wg_slider_float(rf, min, max)
            :size(325),
        gui.wg_button_arrow(gui.DIR_LEFT)
            :on_click(function()
                ref.set(rf, dflt)
            end)
    })
end

function widget.slider_float_step(str, rf, min, max, dflt)
    return gui.wg_row({
        gui.wg_label(str),
        gui.wg_dummy(
            100-gui.calc_text_size_width(str),
            1
        ),
        gui.wg_slider_float(rf, min, max)
            :size(325)
            :format("%.1f"),
        gui.wg_button_arrow(gui.DIR_LEFT)
            :on_click(function()
                ref.set(rf, dflt)
            end)
    })
end

function widget.slider_int(str, rf, min, max, dflt)
    return gui.wg_row({
        gui.wg_label(str),
        gui.wg_dummy(
            100-gui.calc_text_size_width(str),
            1
        ),
        gui.wg_slider_int(rf, min, max)
            :size(325),
        gui.wg_button_arrow(gui.DIR_LEFT)
            :on_click(function()
                ref.set(rf, dflt)
            end)
    })
end

return widget