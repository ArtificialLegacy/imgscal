config({
    name="Example Filters",
    version="1.0.0",
    author="Blub",
    requires={
        "gui",
        "image",
        "ref",
        "collection",
        "io",
        "std",
        "bit",
        "filter",
    },

    desc="GUI example for showcasing different filters"
})

--[[
    Helper functions for creating the sliders.
    These use a dummy widget that gurantees the width of the label+dummy is 100 pixels,
    this is for visuals to make sure the sliders are vertically aligned.

    It also adds an arrow button widget that resets the slider's ref to a default value.
--]]
function wg_slider_float(str, rf, min, max, dflt)
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

function wg_slider_int(str, rf, min, max, dflt)
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

function wg_filter(name, widgets)
    return gui.wg_tree_node(name)
        :flags(bit.bitor_many({
            gui.FLAGTREENODE_FRAMED,
            gui.FLAGTREENODE_NOTREEPUSHONOPEN,
        }))
        :layout(widgets)
end

function wg_ksize(rf)
    return gui.wg_row({
        gui.wg_label("Kernel Size:"),
        gui.wg_button_radio("3", ref.get(rf) == 3)
            :on_change(function()
                ref.set(rf, 3)
            end),
        gui.wg_button_radio("5", ref.get(rf) == 5)
            :on_change(function()
                ref.set(rf, 5)
            end),
        gui.wg_button_radio("7", ref.get(rf) == 7)
            :on_change(function()
                ref.set(rf, 7)
            end),
        gui.wg_button_radio("9", ref.get(rf) == 9)
            :on_change(function()
                ref.set(rf, 9)
            end),
    })
end

convolutionFilters = {
    "None",
    "Identity",
    "Ridge",
    "Emboss",
}

convolutionKernels = {
    { -- identity
        0, 0, 0,
        0, 1, 0,
        0, 0, 0,
    },
    { -- ridge
        0, -1, 0,
        -1, 4, -1,
        0, -1, 0,
    },
    { -- emboss
        -1, -1, 0,
        -1, 1, 1,
        0, 1, 1,
    }
}

cropAnchors = {
    "Center",
    "Top Left",
    "Top",
    "Top Right",
    "Left",
    "Right",
    "Bottom Left",
    "Bottom",
    "Bottom Right",
}

interpolation = {
    "Nearest Neighbor",
    "Linear",
    "Cubic",
}

main(function ()
    local win = gui.window_master("Filter Example", 512, 512, 0)
    gui.window_set_icon_imgscal(win, true)

    local imgSrc = io.load_image(std.wd().."/example_image.png")
    local imgDst = image.new("dst_img", image.ENCODING_PNG, 200, 200)

    -- window padding used for setting the size of the child widget that wraps the images.
    local wpx, wpy = gui.window_padding()

    -- general refs
    local ready = ref.new(true, ref.BOOL)

    -- basic filter refs
    local brightnessEnabled = ref.new(false, ref.BOOL)
    local brightnessPercent = ref.new(0, ref.FLOAT32)
    local contrastEnabled = ref.new(false, ref.BOOL)
    local contrastPercent = ref.new(0, ref.FLOAT32)
    local gammaEnabled = ref.new(false, ref.BOOL)
    local gammaPercent = ref.new(1, ref.FLOAT32)
    local hueEnabled = ref.new(false, ref.BOOL)
    local huePercent = ref.new(0, ref.FLOAT32)
    local saturationEnabled = ref.new(false, ref.BOOL)
    local saturationPercent = ref.new(0, ref.FLOAT32)
    local sepiaEnabled = ref.new(false, ref.BOOL)
    local sepiaPercent = ref.new(0, ref.FLOAT32)

    -- color balance filter refs
    local colorBalanceEnabled = ref.new(false, ref.BOOL)
    local colorBalancePercentR = ref.new(0, ref.FLOAT32)
    local colorBalancePercentG = ref.new(0, ref.FLOAT32)
    local colorBalancePercentB = ref.new(0, ref.FLOAT32)

    -- colorize filter refs
    local colorizeEnabled = ref.new(false, ref.BOOL)
    local colorizeHue = ref.new(0, ref.FLOAT32)
    local colorizeSaturation = ref.new(0, ref.FLOAT32)
    local colorizePercent = ref.new(0, ref.FLOAT32)

    -- colorspace filter refs
    local colorspaceToSRGB = ref.new(false, ref.BOOL)
    local colorspaceToLinear = ref.new(false, ref.BOOL)
    local colorGrayscale = ref.new(false, ref.BOOL)
    local colorInvert = ref.new(false, ref.BOOL)

    -- convolution filter refs
    local convolutionSelected = ref.new(0, ref.INT32)
    local convolutionNormalized = ref.new(false, ref.BOOL)
    local convolutionAlpha = ref.new(false, ref.BOOL)
    local convolutionABS = ref.new(false, ref.BOOL)
    local convolutionDelta = ref.new(0, ref.FLOAT32)
    local maximumEnabled = ref.new(false, ref.BOOL)
    local maximumDisk = ref.new(false, ref.BOOL)
    local maximumKSize = ref.new(3, ref.INT32)
    local meanEnabled = ref.new(false, ref.BOOL)
    local meanDisk = ref.new(false, ref.BOOL)
    local meanKSize = ref.new(3, ref.INT32)
    local medianEnabled = ref.new(false, ref.BOOL)
    local medianDisk = ref.new(false, ref.BOOL)
    local medianKSize = ref.new(3, ref.INT32)
    local minimumEnabled = ref.new(false, ref.BOOL)
    local minimumDisk = ref.new(false, ref.BOOL)
    local minimumKSize = ref.new(3, ref.INT32)

    -- crop filter refs
    local cropEnabled = ref.new(false, ref.BOOL)
    local cropToSizeEnabled = ref.new(false, ref.BOOL)
    local cropxmin = ref.new(0, ref.INT32)
    local cropymin = ref.new(0, ref.INT32)
    local cropxmax = ref.new(200, ref.INT32)
    local cropymax = ref.new(200, ref.INT32)
    local cropwidth = ref.new(200, ref.INT32)
    local cropheight = ref.new(200, ref.INT32)
    local cropanchor = ref.new(0, ref.INT32)

    -- transformation refs
    local flipH = ref.new(false, ref.BOOL)
    local flipV = ref.new(false, ref.BOOL)
    local rotType = ref.new(0, ref.INT32)
    local rotAngle = ref.new(0, ref.FLOAT32)
    local rotInterp = ref.new(0, ref.INT32)
    local etranspose = ref.new(false, ref.BOOL)
    local etransverse = ref.new(false, ref.BOOL)

    -- advanced refs
    local gaussEnabled = ref.new(false, ref.BOOL)
    local gaussSigma = ref.new(1, ref.FLOAT32)
    local pixelateEnabled = ref.new(false, ref.BOOL)
    local pixelateSize = ref.new(0, ref.INT32)
    local thresholdEnabled = ref.new(false, ref.BOOL)
    local thresholdPercent = ref.new(50, ref.FLOAT32)
    local sobelEnabled = ref.new(false, ref.BOOL)
    local sobelBefore = ref.new(true, ref.BOOL)
    local sigmoidEnabled = ref.new(false, ref.BOOL)
    local sigmoidMid = ref.new(0.5, ref.FLOAT32)
    local sigmoidFactor = ref.new(0, ref.FLOAT32)

    gui.window_run(win, function()
        gui.window_single():layout({
            gui.wg_align(gui.ALIGN_CENTER):to({
                gui.wg_style()
                    :set_style_float(gui.STYLEVAR_CHILDROUNDING, 10)
                    :to({
                        gui.wg_child()
                            :size((200 + wpx * 2) * 2 + 50, 200 + wpy * 2)
                            :layout({
                                -- using wg_image_sync here allows it to display while the image is being processed.
                                -- otherwise the main goroutine would be blocked here
                                gui.wg_row({
                                    gui.wg_image_sync(imgSrc)
                                        :size(200, 200),
                                    gui.wg_dummy(50, 1), -- the 50 pixel space between images
                                    gui.wg_image_sync(imgDst)
                                        :size(200, 200),
                                })
                            }),
                    })
            }),
            gui.wg_align(gui.ALIGN_CENTER):to({
                gui.wg_button("Apply Filters")
                    :size(100, 50)
                    :disabled(not ref.get(ready))
                    :on_click(function()
                        ref.set(ready, false)
                        local filters = {}

                        if ref.get(sobelBefore) and ref.get(sobelEnabled) then
                            table.insert(filters, filter.sobel())
                        end

                        if ref.get(colorspaceToSRGB) then
                            table.insert(filters, filter.colorspace_linear_to_srgb())
                        end

                        if ref.get(brightnessEnabled) then
                            table.insert(filters, filter.brightness(ref.get(brightnessPercent)))
                        end

                        if ref.get(contrastEnabled) then
                            table.insert(filters, filter.contrast(ref.get(contrastPercent)))
                        end

                        if ref.get(sigmoidEnabled) then
                            table.insert(filters, filter.sigmoid(
                                ref.get(sigmoidMid),
                                ref.get(sigmoidFactor)
                            ))
                        end

                        if ref.get(gammaEnabled) then
                            table.insert(filters, filter.gamma(ref.get(gammaPercent)))
                        end

                        if ref.get(hueEnabled) then
                            table.insert(filters, filter.hue(ref.get(huePercent)))
                        end

                        if ref.get(saturationEnabled) then
                            table.insert(filters, filter.saturation(ref.get(saturationPercent)))
                        end

                        if ref.get(sepiaEnabled) then
                            table.insert(filters, filter.sepia(ref.get(sepiaPercent)))
                        end

                        if ref.get(colorBalanceEnabled) then
                            table.insert(filters, filter.color_balance(
                                ref.get(colorBalancePercentR),
                                ref.get(colorBalancePercentG),
                                ref.get(colorBalancePercentB)
                            ))
                        end

                        if ref.get(colorizeEnabled) then
                            table.insert(filters, filter.colorize(
                                ref.get(colorizeHue),
                                ref.get(colorizeSaturation),
                                ref.get(colorizePercent)
                            ))
                        end

                        if ref.get(convolutionSelected) > 0 then
                            table.insert(filters, filter.convolution(
                                convolutionKernels[ref.get(convolutionSelected)],
                                ref.get(convolutionNormalized),
                                ref.get(convolutionAlpha),
                                ref.get(convolutionABS),
                                ref.get(convolutionDelta)
                            ))
                        end

                        if ref.get(maximumEnabled) then
                            table.insert(filters, filter.maximum(
                                ref.get(maximumKSize),
                                ref.get(maximumDisk)
                            ))
                        end

                        if ref.get(meanEnabled) then
                            table.insert(filters, filter.mean(
                                ref.get(meanKSize),
                                ref.get(meanDisk)
                            ))
                        end

                        if ref.get(medianEnabled) then
                            table.insert(filters, filter.median(
                                ref.get(medianKSize),
                                ref.get(medianDisk)
                            ))
                        end

                        if ref.get(minimumEnabled) then
                            table.insert(filters, filter.minimum(
                                ref.get(minimumKSize),
                                ref.get(minimumDisk)
                            ))
                        end

                        if ref.get(cropEnabled) then
                            table.insert(filters, filter.crop_xy(
                                ref.get(cropxmin),
                                ref.get(cropymin),
                                ref.get(cropxmax),
                                ref.get(cropymax)
                            ))
                        end

                        if ref.get(cropToSizeEnabled) then
                            table.insert(filters, filter.crop_to_size(
                                ref.get(cropwidth),
                                ref.get(cropheight),
                                ref.get(cropanchor)
                            ))
                        end

                        if ref.get(gaussEnabled) then
                            table.insert(filters, filter.gaussian_blur(ref.get(gaussSigma)))
                        end

                        if ref.get(pixelateEnabled) then
                            table.insert(filters, filter.pixelate(ref.get(pixelateSize)))
                        end

                        if ref.get(flipH) then
                            table.insert(filters, filter.flip_horizontal())
                        end

                        if ref.get(flipV) then
                            table.insert(filters, filter.flip_vertical())
                        end

                        if ref.get(etranspose) then
                            table.insert(filters, filter.transpose())
                        end

                        if ref.get(etransverse) then
                            table.insert(filters, filter.transverse())
                        end

                        if ref.get(rotType) == 90 then
                            table.insert(filters, filter.rotate_90())
                        end

                        if ref.get(rotType) == 180 then
                            table.insert(filters, filter.rotate_180())
                        end

                        if ref.get(rotType) == 270 then
                            table.insert(filters, filter.rotate_270())
                        end

                        if ref.get(rotType) == -1 then
                            table.insert(filters, filter.rotate(
                                ref.get(rotAngle),
                                image.color_rgba_gray(0, 0),
                                ref.get(rotInterp)
                            ))
                        end

                        if ref.get(colorInvert) then
                            table.insert(filters, filter.invert())
                        end

                        if ref.get(colorspaceToLinear) then
                            table.insert(filters, filter.colorspace_srgb_to_linear())
                        end

                        if ref.get(colorGrayscale) then
                            table.insert(filters, filter.grayscale())
                        end

                        if ref.get(thresholdEnabled) then
                            table.insert(filters, filter.threshold(ref.get(thresholdPercent)))
                        end

                        if not ref.get(sobelBefore) and ref.get(sobelEnabled) then
                            table.insert(filters, filter.sobel())
                        end

                        image.clear(imgDst)
                        filter.draw(imgSrc, imgDst, filters)

                        collection.schedule(collection.TYPE_IMAGE, imgDst, function()
                            ref.set(ready, true)
                            gui.update()
                        end)
                    end),
            }),
            wg_filter("Colorspace", {
                gui.wg_checkbox("Linear -> sRGB", colorspaceToSRGB),
                gui.wg_checkbox("sRGB -> Linear", colorspaceToLinear),
                gui.wg_separator(),
                gui.wg_checkbox("Grayscale", colorGrayscale),
                gui.wg_checkbox("Invert", colorInvert),
            }),
            wg_filter("Basic Filters", {
                gui.wg_checkbox("Brightness Enabled", brightnessEnabled),
                wg_slider_float("Percentage:", brightnessPercent, -100, 100, 0),
                gui.wg_separator(),
                gui.wg_checkbox("Contrast Enabled", contrastEnabled),
                wg_slider_float("Percentage:", contrastPercent, -100, 100, 0),
                gui.wg_separator(),
                gui.wg_checkbox("Gamma Enabled", gammaEnabled),
                wg_slider_float("Gamma:", gammaPercent, 0, 2, 1),
                gui.wg_separator(),
                gui.wg_checkbox("Hue Enabled", hueEnabled),
                wg_slider_float("Shift:", huePercent, -180, 180, 0),
                gui.wg_separator(),
                gui.wg_checkbox("Saturation Enabled", saturationEnabled),
                wg_slider_float("Saturation:", saturationPercent, -100, 500, 0),
                gui.wg_separator(),
                gui.wg_checkbox("Sepia Enabled", sepiaEnabled),
                wg_slider_float("Percentage:", sepiaPercent, 0, 100, 0),
            }),
            wg_filter("Color Balance", {
                gui.wg_checkbox("Enabled", colorBalanceEnabled),
                wg_slider_float("Red %:", colorBalancePercentR, -100, 500, 0),
                wg_slider_float("Green %:", colorBalancePercentG, -100, 500, 0),
                wg_slider_float("Blue %:", colorBalancePercentB, -100, 500, 0),
            }),
            wg_filter("Colorize", {
                gui.wg_checkbox("Enabled", colorizeEnabled),
                wg_slider_float("Hue:", colorizeHue, 0, 360, 0),
                wg_slider_float("Saturation:", colorizeSaturation, 0, 100, 0),
                wg_slider_float("Percent:", colorizePercent, 0, 100, 0),
            }),
            wg_filter("Convolution Filter", {
                gui.wg_combo_preview(
                    "Kernel",
                    convolutionFilters,
                    convolutionSelected
                ),
                gui.wg_row({
                    gui.wg_checkbox("Normalized", convolutionNormalized),
                    gui.wg_checkbox("Alpha", convolutionAlpha),
                    gui.wg_checkbox("ABS", convolutionABS),
                }),
                wg_slider_float("Delta:", convolutionDelta, -2, 2, 0),

                gui.wg_separator(),
                
                gui.wg_row({
                    gui.wg_checkbox("Maximum Enabled", maximumEnabled),
                    gui.wg_checkbox("Disk", maximumDisk),
                }),
                wg_ksize(maximumKSize),
                gui.wg_separator(),
                gui.wg_row({
                    gui.wg_checkbox("Mean Enabled", meanEnabled),
                    gui.wg_checkbox("Disk", meanDisk),
                }),
                wg_ksize(meanKSize),
                gui.wg_separator(),
                gui.wg_row({
                    gui.wg_checkbox("Median Enabled", medianEnabled),
                    gui.wg_checkbox("Disk", medianDisk),
                }),
                wg_ksize(medianKSize),
                gui.wg_separator(),
                gui.wg_row({
                    gui.wg_checkbox("Minimum Enabled", minimumEnabled),
                    gui.wg_checkbox("Disk", minimumDisk),
                }),
                wg_ksize(minimumKSize),
            }),
            wg_filter("Crop", {
                gui.wg_tab_bar():tab_items({
                    gui.wg_tab_item("Crop Rectangle"):layout({
                        gui.wg_checkbox("Enabled", cropEnabled)
                            :on_change(function(b)
                                if b then
                                    ref.set(cropToSizeEnabled, false)
                                end
                            end),
                        gui.wg_label(std.fmt("Cropping between points: (%d,%d)-(%d,%d)", {
                            ref.get(cropxmin),
                            ref.get(cropymin),
                            ref.get(cropxmax),
                            ref.get(cropymax),
                        })),
                        wg_slider_int("Min X:", cropxmin, 0, 200, 0),
                        wg_slider_int("Max X:", cropxmax, 0, 200, 200),
                        wg_slider_int("Min Y:", cropymin, 0, 200, 0),
                        wg_slider_int("Max Y:", cropymax, 0, 200, 200),
                    }),
                    gui.wg_tab_item("Crop Size"):layout({
                        gui.wg_checkbox("Enabled", cropToSizeEnabled)
                            :on_change(function(b)
                                if b then
                                    ref.set(cropEnabled, false)
                                end
                            end),
                        wg_slider_int("Width:", cropwidth, 1, 200, 200),
                        wg_slider_int("Height:", cropheight, 1, 200, 200),
                        gui.wg_combo_preview(
                            "Anchor",
                            cropAnchors,
                            cropanchor
                        ),
                    }),
                }),
            }),
            wg_filter("Transformations", {
                gui.wg_row({
                    gui.wg_label("Flip:"),
                    gui.wg_checkbox("Horizontal", flipH),
                    gui.wg_checkbox("Vertical", flipV),
                }),
                gui.wg_row({
                    gui.wg_checkbox("Transpose", etranspose),
                    gui.wg_checkbox("Transverse", etransverse),
                }),
                gui.wg_separator(),
                gui.wg_row({
                    gui.wg_label("Rotations:"),
                    gui.wg_button_radio("None", ref.get(rotType) == 0)
                        :on_change(function()
                            ref.set(rotType, 0)
                        end),
                    gui.wg_button_radio("90", ref.get(rotType) == 90)
                        :on_change(function()
                            ref.set(rotType, 90)
                        end),
                    gui.wg_button_radio("180", ref.get(rotType) == 180)
                        :on_change(function()
                            ref.set(rotType, 180)
                        end),
                    gui.wg_button_radio("270", ref.get(rotType) == 270)
                        :on_change(function()
                            ref.set(rotType, 270)
                        end),
                    gui.wg_button_radio("Custom", ref.get(rotType) == -1)
                        :on_change(function()
                            ref.set(rotType, -1)
                        end),
                }),
                gui.wg_combo_preview(
                    "Interpolation",
                    interpolation,
                    rotInterp
                ),
                wg_slider_float("Angle:", rotAngle, 0, 360, 0),
            }),
            wg_filter("Advanced Filters", {
                gui.wg_checkbox("Gaussian Blur Enabled", gaussEnabled),
                wg_slider_float("Sigma:", gaussSigma, 0.1, 5, 1),

                gui.wg_separator(),
                
                gui.wg_checkbox("Pixelate Enabled", pixelateEnabled),
                wg_slider_int("Size:", pixelateSize, 0, 15, 0),
                
                gui.wg_separator(),
                
                gui.wg_checkbox("Threshold Enabled", thresholdEnabled),
                wg_slider_float("Percentage:", thresholdPercent, 0, 100, 50),
                
                gui.wg_separator(),
                
                gui.wg_row({
                    gui.wg_checkbox("Sobel Enabled", sobelEnabled),
                    gui.wg_button_radio("Before", ref.get(sobelBefore))
                        :on_change(function()
                            ref.set(sobelBefore, true)
                        end),
                    gui.wg_button_radio("After", not ref.get(sobelBefore))
                        :on_change(function()
                            ref.set(sobelBefore, false)
                        end),
                }),
                
                gui.wg_separator(),
                
                gui.wg_checkbox("Sigmoid Enabled", sigmoidEnabled),
                wg_slider_float("Midpoint:", sigmoidMid, 0, 1, 0.5),
                wg_slider_float("Factor:", sigmoidFactor, -10, 10, 0),
            })
        })
    end)
end)