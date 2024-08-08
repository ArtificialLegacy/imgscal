config({
    name="GUI Test",
    version="0.1.0",
    author="Blub",
    requires={
        "gui",
        "bit",
        "image",
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
            gui.wg_label("canvas test"),
            gui.wg_custom(function()
                local pos = gui.cursor_screen_pos()

                local pos0 = image.point(pos.x, pos.y+50)
                local cp0 = image.point(pos.x+50, pos.y+50)
                local cp1 = image.point(pos.x+20, pos.y+150)
                local pos1 = image.point(pos.x+100, pos.y+150)
                local col = image.color_rgb(255, 0, 0)
                gui.canvas_bezier_cubic(pos0, cp0, cp1, pos1, col, 1, 5)

                gui.canvas_circle(pos1, 20, col, 10, 1)
                local pos2 = image.point(pos1.x+50, pos1.y)
                gui.canvas_circle_filled(pos2, 20, col)

                gui.canvas_line(pos0, pos1, col, 1)

                gui.canvas_quad(pos, pos0, pos1, pos2, col, 1)
                local pos3 = image.point(pos2.x+20, pos2.y+20)
                local pos4 = image.point(pos1.x-20, pos1.y-20)
                gui.canvas_quad_filled(pos, pos0, pos3, pos4, col)

                gui.canvas_rect(
                    image.point(pos.x+256, pos.y+256),
                    image.point(pos.x+512, pos.y+512),
                    col,
                    20, gui.FLAGDRAW_ROUNDCORNERSALL, 1
                )

                gui.canvas_rect_filled(
                    image.point(pos.x+300, pos.y+300),
                    image.point(pos.x+490, pos.y+490),
                    col,
                    20, bit.bitor(
                        gui.FLAGDRAW_ROUNDCORNERSTOPLEFT,
                        gui.FLAGDRAW_ROUNDCORNERSBOTTOMRIGHT
                    )
                )

                local col2 = image.color_rgb_gray(255)
                gui.canvas_text(
                    image.point(pos.x+100, pos.y+100),
                    col2, "canvas text"
                )

                gui.canvas_triangle(
                    image.point(pos.x+500, pos.y+30),
                    image.point(pos.x+100, pos.y+500),
                    image.point(pos.x+600, pos.y+500),
                    col2, 1
                )
                gui.canvas_triangle_filled(
                    image.point(pos.x+500, pos.y+50),
                    image.point(pos.x+450, pos.y+100),
                    image.point(pos.x+550, pos.y+100),
                    col2
                )

                gui.canvas_path_arc_to(
                    image.point(pos.x+700, pos.y+700),
                    30, 1, 3, 4 
                )

                gui.canvas_path_arc_to_fast(
                    image.point(pos.x+500, pos.y+500),
                    30, 1, 3
                )

                gui.canvas_path_bezier_cubic_to(
                    image.point(pos.x+600, pos.y+600),
                    image.point(pos.x+650, pos.y+650),
                    image.point(pos.x+600, pos.y+750),
                    5
                )

                gui.canvas_path_stroke(col2, 0, 1)

                gui.canvas_path_line_to(image.point(pos.x+400, pos.y+400))
                gui.canvas_path_line_to(image.point(pos.x+410, pos.y+400))
                gui.canvas_path_line_to(image.point(pos.x+400, pos.y+410))
                gui.canvas_path_line_to(image.point(pos.x+410, pos.y+410))
                
                local col3 = image.color_rgb_gray(100)
                gui.canvas_path_fill_convex(col3)
            end),
        })
    end)
end)