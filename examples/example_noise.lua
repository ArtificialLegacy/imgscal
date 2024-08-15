config({
    name="Example Noise",
    version="1.0.0",
    author="Blub",
    requires={
        "gui",
        "image",
        "noise",
        "ref",
        "collection",
        "cli",
    },

    desc="GUI example for generating noise maps with opensimplex noise."
})

function newSeed(seed)
    ref.set(seed, math.random(0, 999999))
end

--[[
    Noise generation functions that are run when the 'Generate' button is clicked.

    These schedule the noise generation tasks onto the image first,
    then schedule an addition custom task on the image.

    This custom task is added to the queue to be run after the first task is finished
    to update the ready ref and manually refresh the gui to ensure the new image is displayed.

    Do note that manually updating the gui does not always work due to upstream issues
    related to windows compatibility and mainthread.
    This is due to it being run within the goroutine for the image, instead of the main goroutine.
--]]
function genNoise2(img, seed, scale, direct, ecolor, ealpha, normal, ready)
    ref.set(ready, false)

    if direct then
        noise.simplex_image_map(seed, scale, img, normal, not ecolor, not ealpha)
    else
        image.map(img, function(x, y, c)
            local v = noise.simplex_2d(seed, x * scale, y * scale, normal)
            local i = math.floor(255 * v)
            return image.color_gray(i)
        end)
    end

    collection.schedule(collection.IMAGE, img, function()
        ref.set(ready, true)
        gui.update()
    end)
end

function genNoise3(img, seed, scale, aoff, normal, ready)
    ref.set(ready, false)
    image.map(img, function(x, y, c)
        local v1 = noise.simplex_3d(seed, x * scale, y * scale, 0, normal)
        local v2 = noise.simplex_3d(seed, x * scale, y * scale, aoff, normal)
        local i1 = math.floor(255 * v1)
        local i2 = math.floor(255 * v2)
        return image.color_graya(i1, i2)
    end)

    collection.schedule(collection.IMAGE, img, function()
        ref.set(ready, true)
        gui.update()
    end)
end

function genNoise4(img, seed, scale, goff, boff, normal, ready)
    ref.set(ready, false)
    image.map(img, function(x, y, c)
        local v1 = noise.simplex_4d(seed, x * scale, y * scale, 0, 0, normal)
        local v2 = noise.simplex_4d(seed, x * scale, y * scale, goff, 0, normal)
        local v3 = noise.simplex_4d(seed, x * scale, y * scale, 0, boff, normal)
        local i1 = math.floor(255 * v1)
        local i2 = math.floor(255 * v2)
        local i3 = math.floor(255 * v3)
        return image.color_rgb(i1, i2, i3)
    end)

    collection.schedule(collection.IMAGE, img, function()
        ref.set(ready, true)
        gui.update()
    end)
end

simplex2_desc = [[
r,g,b = simplex2d(x*scale, y*scale)
Enabling direct mode will use simplex_image_map instead, this is implemented in go and as such is much faster.
Toggling color and alpha only affect direct mode, as they are a part of the simplex_image_map function.
]]

simplex3_desc = [[
r,g,b = simplex3d(x*scale, y*scale, 0)
alpha = simplex3d(x*scale, y*scale, offset)
]]

simplex4_desc = [[
r = simplex4d(x*scale, y*scale, 0, 0)
g = simplex4d(x*scale, y*scale, green_offset, 0)
b = simplex4d(x*scale, y*scale, 0, blue_offset)
]]

--[[
    Helper function for building widgets for the seed settings.
--]]
function wg_seed(seed, rnd, normal)
    return gui.wg_row({
        gui.wg_label("Seed:"),
        gui.wg_input_int(seed):size(100),
        gui.wg_button("Randomize")
            :on_click(function ()
                newSeed(seed)
            end),
        gui.wg_checkbox("On Generate", rnd),
        gui.wg_checkbox("Normalize", normal),
    })
end

--[[
    Helper functions that create the `Generate` buttons for each of the tabs.
    Also calls `newSeed` if the randomize on generate checkbox is checked.

    These get the values from each ref on click, as this is a copy 
    the related controls do not need to be disabled while the image is being processed.
--]]
function wg_gen2(img, ready, seed, rnd, scale, direct, ecolor, ealpha, normal)
    return gui.wg_align(gui.ALIGN_CENTER):to({
        gui.wg_button("Generate")
            :disabled(not ref.get(ready))
            :size(100, 50)
            :on_click(function()
                genNoise2(img,
                    ref.get(seed), 
                    ref.get(scale),
                    ref.get(direct),
                    ref.get(ecolor),
                    ref.get(ealpha),
                    ref.get(normal),
                    ready
                )

                if ref.get(rnd) then
                    newSeed(seed)
                end
            end)
    })
end

function wg_gen3(img, ready, seed, rnd, scale, aoff, normal)
    return gui.wg_align(gui.ALIGN_CENTER):to({
        gui.wg_button("Generate")
            :disabled(not ref.get(ready))
            :size(100, 50)
            :on_click(function()
                genNoise3(img,
                    ref.get(seed),
                    ref.get(scale),
                    ref.get(aoff),
                    ref.get(normal),
                    ready
                )

                if ref.get(rnd) then
                    newSeed(seed)
                end
            end)
    })
end

function wg_gen4(img, ready, seed, rnd, scale, goff, boff, normal)
    return gui.wg_align(gui.ALIGN_CENTER):to({
        gui.wg_button("Generate")
            :disabled(not ref.get(ready))
            :size(100, 50)
            :on_click(function()
                genNoise4(img,
                    ref.get(seed),
                    ref.get(scale),
                    ref.get(goff),
                    ref.get(boff),
                    ref.get(normal),
                    ready
                )

                if ref.get(rnd) then
                    newSeed(seed)
                end
            end)
    })
end

--[[
    A helper function for creating the sliders.
    These use a dummy widget that gurantees the width of the label+dummy is 100 pixels,
    this is for visuals to make sure the sliders are vertically aligned.

    It also adds an arrow button widget that resets the slider's ref to a default value.
--]]
function wg_slider(str, rf, min, max, dflt)
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

main(function ()
    local win = gui.window_master("Noise Example", 512, 512, 0)
    gui.window_set_icon_imgscal(win, true)

    local img = image.new("noise_img", image.ENCODING_PNG, 100, 100)

    -- refs used for general controls
    local seed = ref.new(0, ref.INT32)
    local scale = ref.new(0.25, ref.FLOAT32)
    local rnd = ref.new(true, ref.BOOL)
    local ready = ref.new(true, ref.BOOL)
    local normal = ref.new(true, ref.BOOL)

    -- refs used for the 2d controls
    local direct = ref.new(false, ref.BOOL)
    local ecolor = ref.new(true, ref.BOOL)
    local ealpha = ref.new(false, ref.BOOL)

    -- refs used for the 3d controls
    local aoff = ref.new(0.5, ref.FLOAT32)

    -- refs used for the 4d controls
    local goff = ref.new(0.5, ref.FLOAT32)
    local boff = ref.new(0.5, ref.FLOAT32)
    
    -- window padding used for setting the size of the child widget that wraps the image.
    local wpx, wpy = gui.window_padding()

    gui.window_run(win, function()
        gui.window_single():layout({gui.wg_column({
            gui.wg_align(gui.ALIGN_CENTER):to({
                gui.wg_style()
                    :set_style_float(gui.STYLEVAR_CHILDROUNDING, 10)
                    :to({
                        gui.wg_child()
                            :size(200 + wpx * 2, 200 + wpy * 2)
                            :layout({
                                -- using wg_image_sync here allows it to display while the image is being processed.
                                -- otherwise the main goroutine would be blocked here
                                gui.wg_image_sync(img)
                                    :size(200, 200),
                            }),
                    }),
            }),
            gui.wg_style()
                :set_style_float(gui.STYLEVAR_CHILDROUNDING, 10)
                :to({
                    gui.wg_child():layout({
                        gui.wg_tab_bar():tab_items({
                            gui.wg_tab_item("simplex_2d"):layout({
                                gui.wg_label(simplex2_desc)
                                    :wrapped(true),
                                gui.wg_separator(),
                                gui.wg_spacing(),
                                gui.wg_spacing(),
                                wg_seed(seed, rnd, normal),
                                wg_slider("Scale Factor:", scale, 0.01, 0.5, 0.25),
                                gui.wg_row({
                                    gui.wg_checkbox("Direct", direct),
                                    gui.wg_checkbox("Enable Color", ecolor),
                                    gui.wg_checkbox("Enable Alpha", ealpha),
                                }),
                                wg_gen2(img, ready, seed, rnd, scale, direct, ecolor, ealpha, normal),
                            }),
                            gui.wg_tab_item("simplex_3d"):layout({
                                gui.wg_label(simplex3_desc)
                                    :wrapped(true),
                                gui.wg_separator(),
                                gui.wg_spacing(),
                                gui.wg_spacing(),
                                wg_seed(seed, rnd, normal),
                                wg_slider("Scale Factor:", scale, 0.01, 0.5, 0.25),
                                wg_slider("Alpha Offset:", aoff, 0.1, 1.5, 0.5),
                                wg_gen3(img, ready, seed, rnd, scale, aoff, normal),
                            }),
                            gui.wg_tab_item("simplex_4d"):layout({
                                gui.wg_label(simplex4_desc)
                                    :wrapped(true),
                                gui.wg_separator(),
                                gui.wg_spacing(),
                                gui.wg_spacing(),
                                wg_seed(seed, rnd, normal),
                                wg_slider("Scale Factor:", scale, 0.01, 0.5, 0.25),
                                wg_slider("Green Offset:", goff, 0.1, 1.5, 0.5),
                                wg_slider("Blue Offset:", boff, 0.1, 1.5, 0.5),
                                wg_gen4(img, ready, seed, rnd, scale, goff, boff, normal),
                            }),
                    }),
                }),
            }),
        })})
    end)
end)