---@param workflow imgscal_WorkflowInit
function init(workflow)
	workflow.import({
		"gui",
		"guiplot",
		"image",
		"imger",
		"ref",
		"collection",
		"io",
		"bit",
		"noise",
	})
end

ready = true

function main()
	local widget = require("widget")

	local win = gui.window_master("Imger Example", 512, 512, 0)
	gui.window_set_icon_imgscal(win, true)

	local imgOrigin = io.decode(io.path_join(io.wd(), "example_image.png"))
	local imgEmpty = image.new("empty", image.ENCODING_PNG, 200, 200)
	local imgSrc = image.copy(imgOrigin, "src", image.MODEL_RGBA)
	local imgDst = image.new("dest", image.ENCODING_PNG, 200, 200)
	local imgHistogram =
		image.new_filled("histogram", image.ENCODING_PNG, 256, 256, image.color_rgb(0, 0, 0), image.MODEL_RGBA)

	local wpx, wpy = gui.window_padding()
	local spx, _ = gui.item_spacing()
	local isx, _ = gui.item_inner_spacing()

	-- source refs
	local srcOption = ref.new(0, ref.INT32)
	local srcGradientStart = ref.new(image.color_rgb(0, 0, 0), ref.RGBA) --[[@as ref_RGBA]]
	local srcGradientEnd = ref.new(image.color_rgb(255, 255, 255), ref.RGBA) --[[@as ref_RGBA]]
	local srcGradientDir = ref.new(0, ref.INT32)
	local srcNoiseSeed = ref.new(0, ref.INT32)
	local srcNoiseScale = ref.new(0.5, ref.FLOAT32)
	local srcNormalize = ref.new(true, ref.BOOL)
	local histogramMode = ref.new(0, ref.INT32)
	local histogramUseGray = ref.new(false, ref.BOOL)
	local histogramUseRGB = ref.new(false, ref.BOOL)
	local histogramPlotGray = ref.new(false, ref.BOOL)
	local histogramPlotRed = ref.new(false, ref.BOOL)
	local histogramPlotGreen = ref.new(false, ref.BOOL)
	local histogramPlotBlue = ref.new(false, ref.BOOL)

	local histogramDataGray = { 0 }
	local histogramDataRed = { 0 }
	local histogramDataGreen = { 0 }
	local histogramDataBlue = { 0 }

	-- grayscale refs
	local grayscaleUse = ref.new(false, ref.BOOL)
	local grayscaleSelfEnabled = ref.new(false, ref.BOOL)
	local grayscaleWeightedEnabled = ref.new(false, ref.BOOL)
	local grayscaleWeight1 = ref.new(0.5, ref.FLOAT32)
	local grayscaleWeight2 = ref.new(0.5, ref.FLOAT32)
	local grayscaleScalarEnabled = ref.new(false, ref.BOOL)
	local grayscaleScalar = ref.new(127, ref.INT32)
	local grayscaleEnabled = ref.new(false, ref.BOOL)
	local grayscaleTo16 = ref.new(false, ref.BOOL)

	-- blur refs
	local blurBoxEnabled = ref.new(false, ref.BOOL)
	local blurBoxWidth = ref.new(3, ref.INT32)
	local blurBoxHeight = ref.new(3, ref.INT32)
	local blurBoxAX = ref.new(1, ref.INT32)
	local blurBoxAY = ref.new(1, ref.INT32)
	local blurBoxBorder = ref.new(0, ref.INT32)
	local blurBoxGray = ref.new(false, ref.BOOL)
	local gaussEnabled = ref.new(false, ref.BOOL)
	local gaussRadius = ref.new(3.0, ref.FLOAT32)
	local gaussSigma = ref.new(1.0, ref.FLOAT32)
	local gaussBorder = ref.new(0, ref.INT32)
	local gaussGray = ref.new(false, ref.BOOL)

	-- edge detection refs
	local edgeCannyEnabled = ref.new(false, ref.BOOL)
	local edgeCannyLower = ref.new(0.5, ref.FLOAT32)
	local edgeCannyUpper = ref.new(1.0, ref.FLOAT32)
	local edgeCannyKernel = ref.new(3, ref.INT32)
	local edgeCannyGray = ref.new(false, ref.BOOL)
	local edgeSobelEnabled = ref.new(false, ref.BOOL)
	local edgeSobelDir = ref.new(0, ref.INT32)
	local edgeSobelBorder = ref.new(0, ref.INT32)
	local edgeSobelGray = ref.new(false, ref.BOOL)
	local edgeLaplacianEnabled = ref.new(false, ref.BOOL)
	local edgeLaplacianKernel = ref.new(0, ref.INT32)
	local edgeLaplacianBorder = ref.new(0, ref.INT32)
	local edgeLaplacianGray = ref.new(false, ref.BOOL)

	-- effects refs
	local effectEmbossEnabled = ref.new(false, ref.BOOL)
	local effectEmbossGray = ref.new(false, ref.BOOL)
	local effectInvertEnabled = ref.new(false, ref.BOOL)
	local effectInvertGray = ref.new(false, ref.BOOL)
	local effectPixelateEnabled = ref.new(false, ref.BOOL)
	local effectPixelateSize = ref.new(5.0, ref.FLOAT32)
	local effectPixelateGray = ref.new(false, ref.BOOL)
	local effectSharpenEnabled = ref.new(false, ref.BOOL)
	local effectSharpenGray = ref.new(false, ref.BOOL)
	local effectSepiaEnabled = ref.new(false, ref.BOOL)

	-- padding refs
	local paddingEnabled = ref.new(false, ref.BOOL)
	local paddingMode = ref.new(0, ref.INT32)
	local paddingKernelWidth = ref.new(3, ref.INT32)
	local paddingKernelHeight = ref.new(3, ref.INT32)
	local paddingKernelAX = ref.new(1, ref.INT32)
	local paddingKernelAY = ref.new(1, ref.INT32)
	local paddingBorder = ref.new(0, ref.INT32)
	local paddingGray = ref.new(false, ref.BOOL)
	local paddingTop = ref.new(1, ref.INT32)
	local paddingBottom = ref.new(1, ref.INT32)
	local paddingLeft = ref.new(1, ref.INT32)
	local paddingRight = ref.new(1, ref.INT32)

	-- threshold refs
	local thresholdEnabled = ref.new(false, ref.BOOL)
	local thresholdMode = ref.new(0, ref.INT32)
	local thresholdValue = ref.new(127, ref.INT32)
	local thresholdMethod = ref.new(0, ref.INT32)
	local thresholdTo16 = ref.new(false, ref.BOOL)

	-- transform refs
	local rotateEnabled = ref.new(false, ref.BOOL)
	local rotateAngle = ref.new(180.0, ref.FLOAT32)
	local rotateAX = ref.new(100, ref.INT32)
	local rotateAY = ref.new(100, ref.INT32)
	local rotateResize = ref.new(false, ref.BOOL)
	local rotateGray = ref.new(false, ref.BOOL)
	local resizeEnabled = ref.new(false, ref.BOOL)
	local resizeScalex = ref.new(1.0, ref.FLOAT32)
	local resizeScaley = ref.new(1.0, ref.FLOAT32)
	local resizeInterp = ref.new(0, ref.INT32)
	local resizeGray = ref.new(false, ref.BOOL)

	_ = gui.window_run(win, function()
		gui.window_single():layout({
			gui.wg_align(gui.ALIGN_CENTER):to({
				gui.wg_style():set_style_float(gui.STYLEVAR_CHILDROUNDING, 10):to({
					gui.wg_child():size((200 + wpx * 2) * 2 + 50, 200 + wpy * 2 + 4):layout({
						-- using wg_image_sync here allows it to display while the image is being processed.
						-- otherwise the main goroutine would be blocked here
						gui.wg_align(gui.ALIGN_CENTER):to({
							gui.wg_table()
								:flags(bit.bitor_many(gui.FLAGTABLE_BORDERSINNERV, gui.FLAGTABLE_NOPADOUTERX))
								:size(450, 200)
								:rows({
									gui.wg_table_row({
										gui.wg_image_sync(imgSrc):size(200, 200),
										gui.wg_row({
											gui.wg_dummy(12, 1),
											gui.wg_image_sync(imgDst):size(200, 200),
										}),
									}),
								}),
						}),
					}),
				}),
			}),
			gui.wg_custom(function()
				if ref.get(histogramMode) == 1 then
					gui.layout({
						gui.wg_align(gui.ALIGN_CENTER):to({
							gui.wg_style():set_style_float(gui.STYLEVAR_CHILDROUNDING, 10):to({
								gui.wg_child():size((256 + wpx * 2), 256 + wpy * 2 + 4):layout({
									-- using wg_image_sync here allows it to display while the image is being processed.
									-- otherwise the main goroutine would be blocked here
									gui.wg_image_sync(imgHistogram):size(256, 256),
								}),
							}),
						}),
					})
				elseif ref.get(histogramMode) == 2 then
					gui.layout({
						gui.wg_align(gui.ALIGN_CENTER):to({
							gui.wg_style():set_style_float(gui.STYLEVAR_CHILDROUNDING, 10):to({
								gui.wg_child():size((256 + wpx * 2), 256 + wpy * 2 + 4):layout({
									guiplot
										.wg_plot("Histogram")
										:size(256, 256)
										:axis_limits(0, 256, 0, 100, gui.COND_ALWAYS)
										:x_axeflags(guiplot.FLAGPLOTAXIS_NOHIGHLIGHT)
										:y_axeflags(
											guiplot.FLAGPLOTAXIS_NOHIGHLIGHT,
											guiplot.FLAGPLOTAXIS_NOHIGHLIGHT,
											guiplot.FLAGPLOTAXIS_NOHIGHLIGHT
										)
										:plots({
											guiplot.pt_line("Gray", histogramDataGray),
											guiplot.pt_line("Red", histogramDataRed),
											guiplot.pt_line("Green", histogramDataGreen),
											guiplot.pt_line("Blue", histogramDataBlue),
										}),
								}),
							}),
						}),
					})
				end
			end),
			gui.wg_align(gui.ALIGN_CENTER):to({
				gui.wg_button("Apply Filters"):size(100, 50):disabled(not ready):on_click(function()
					ready = false

					if ref.get(grayscaleUse) then
						image.clone(imgDst, imgSrc, image.MODEL_GRAY)
					else
						image.clone(imgDst, imgSrc, image.MODEL_RGBA)
					end

					if ref.get(paddingEnabled) then
						if ref.get(paddingMode) == 0 then
							imger.padding_xy(
								imgDst,
								ref.get(paddingKernelWidth),
								ref.get(paddingKernelHeight),
								ref.get(paddingKernelAX),
								ref.get(paddingKernelAY),
								ref.get(paddingBorder),
								ref.get(paddingGray)
							)
						else
							imger.padding_size(
								imgDst,
								ref.get(paddingTop),
								ref.get(paddingBottom),
								ref.get(paddingLeft),
								ref.get(paddingRight),
								ref.get(paddingBorder),
								ref.get(paddingGray)
							)
						end
					end

					if ref.get(effectSepiaEnabled) then
						imger.sepia(imgDst)
						if ref.get(grayscaleUse) then
							image.convert(imgDst, image.MODEL_GRAY)
						end
					end

					if ref.get(grayscaleEnabled) then
						imger.grayscale_inplace(imgDst, ref.get(grayscaleTo16))
					end

					if ref.get(grayscaleSelfEnabled) then
						imger.gray_add(imgDst, imgSrc)
					end

					if ref.get(grayscaleWeightedEnabled) then
						imger.gray_add_weighted(imgDst, imgSrc, ref.get(grayscaleWeight1), ref.get(grayscaleWeight2))
					end

					if ref.get(grayscaleScalarEnabled) then
						imger.gray_add_scalar(imgDst, ref.get(grayscaleScalar))
					end

					if ref.get(effectEmbossEnabled) then
						imger.emboss(imgDst, ref.get(effectEmbossGray))
					end

					if ref.get(effectInvertEnabled) then
						imger.invert(imgDst, ref.get(effectInvertGray))
					end

					if ref.get(effectPixelateEnabled) then
						imger.pixelate(imgDst, ref.get(effectPixelateSize), ref.get(effectPixelateGray))
					end

					if ref.get(effectSharpenEnabled) then
						imger.sharpen(imgDst, ref.get(effectSharpenGray))
					end

					if ref.get(blurBoxEnabled) then
						local width = ref.get(blurBoxWidth)
						local height = ref.get(blurBoxHeight)
						local ax = ref.get(blurBoxAX)
						local ay = ref.get(blurBoxAY)

						if ax >= width then
							ax = width - 1
						end
						if ay >= height then
							ay = height - 1
						end

						imger.blur_box_xy(imgDst, width, height, ax, ay, ref.get(blurBoxBorder), ref.get(blurBoxGray))
					end

					if ref.get(gaussEnabled) then
						imger.blur_gaussian(
							imgDst,
							ref.get(gaussRadius),
							ref.get(gaussSigma),
							ref.get(gaussBorder),
							ref.get(gaussGray)
						)
					end

					if ref.get(thresholdEnabled) then
						if ref.get(thresholdMode) == 0 then
							local t = ref.get(thresholdValue)
							if ref.get(thresholdTo16) then
								t = t / 255 * 65535
							end

							imger.threshold(imgDst, t, ref.get(thresholdMethod), ref.get(thresholdTo16))
						else
							imger.threshold_otsu(imgDst, ref.get(thresholdMethod))
						end
					end

					if ref.get(edgeCannyEnabled) then
						imger.edge_canny_inplace(
							imgDst,
							ref.get(edgeCannyLower),
							ref.get(edgeCannyUpper),
							ref.get(edgeCannyKernel),
							ref.get(edgeCannyGray)
						)
					end

					if ref.get(edgeSobelEnabled) then
						local dir = ref.get(edgeSobelDir)
						if dir == 0 then
							imger.edge_sobel_inplace(imgDst, ref.get(edgeSobelBorder), ref.get(edgeSobelGray))
						elseif dir == 1 then
							imger.edge_sobel_horizontal_inplace(
								imgDst,
								ref.get(edgeSobelBorder),
								ref.get(edgeSobelGray)
							)
						elseif dir == 2 then
							imger.edge_sobel_vertical_inplace(imgDst, ref.get(edgeSobelBorder), ref.get(edgeSobelGray))
						end
					end

					if ref.get(edgeLaplacianEnabled) then
						imger.edge_laplacian_inplace(
							imgDst,
							ref.get(edgeLaplacianKernel),
							ref.get(edgeLaplacianBorder),
							ref.get(edgeLaplacianGray)
						)
					end

					if ref.get(rotateEnabled) then
						imger.rotate_xy(
							imgDst,
							ref.get(rotateAngle),
							ref.get(rotateAX),
							ref.get(rotateAY),
							ref.get(rotateResize),
							ref.get(rotateGray)
						)
					end

					if ref.get(resizeEnabled) then
						imger.resize(
							imgDst,
							ref.get(resizeScalex),
							ref.get(resizeScaley),
							ref.get(resizeInterp),
							ref.get(resizeGray)
						)
					end

					if ref.get(histogramMode) == 1 then
						local gray = ref.get(histogramUseGray)
						if gray then
							image.clone(imgHistogram, imgDst, image.MODEL_GRAY)
						else
							image.clone(imgHistogram, imgDst, image.MODEL_RGBA)
						end

						imger.histogram_draw_inplace_xy(imgHistogram, 1, 1, gray)
					elseif ref.get(histogramMode) == 2 then
						histogramDataGray = { 0 }
						histogramDataRed = { 0 }
						histogramDataGreen = { 0 }
						histogramDataBlue = { 0 }

						if ref.get(histogramPlotGray) then
							histogramDataGray = imger.histogram_gray(imgDst, 100, true)
						end

						if ref.get(histogramUseRGB) then
							local red, green, blue = imger.histogram_rgb(imgDst, 100, true)

							if ref.get(histogramPlotRed) then
								histogramDataRed = red
							end

							if ref.get(histogramPlotGreen) then
								histogramDataGreen = green
							end

							if ref.get(histogramPlotBlue) then
								histogramDataBlue = blue
							end
						else
							if ref.get(histogramPlotRed) then
								histogramDataRed = imger.histogram_red(imgDst, 100, true)
							end

							if ref.get(histogramPlotGreen) then
								histogramDataGreen = imger.histogram_green(imgDst, 100, true)
							end

							if ref.get(histogramPlotBlue) then
								histogramDataBlue = imger.histogram_blue(imgDst, 100, true)
							end
						end
					end

					collection.schedule(collection.IMAGE, imgDst, function()
						_ = collection.wait(collection.IMAGE, imgHistogram)
						ready = true
						gui.update()
					end)
				end),
			}),
			wg_filter("Source", {
				gui.wg_row({
					gui.wg_label("Histogram"),
					gui.wg_button_radio("None", ref.get(histogramMode) == 0):on_change(function()
						ref.set(histogramMode, 0)
					end),
					gui.wg_button_radio("Image", ref.get(histogramMode) == 1):on_change(function()
						ref.set(histogramMode, 1)
					end),
					gui.wg_button_radio("Plot", ref.get(histogramMode) == 2):on_change(function()
						ref.set(histogramMode, 2)
					end),
				}),
				gui.wg_checkbox("Image use GRAY.", histogramUseGray),
				gui.wg_checkbox("Plot Gray.", histogramPlotGray),
				gui.wg_checkbox("Plot use RGB.", histogramUseRGB),
				gui.wg_row({
					gui.wg_checkbox("Plot Red.", histogramPlotRed),
					gui.wg_checkbox("Plot Green.", histogramPlotGreen),
					gui.wg_checkbox("Plot Blue.", histogramPlotBlue),
				}),
				gui.wg_separator(),
				gui.wg_row({
					gui.wg_button("Refresh Source"):disabled(not ready):on_click(function()
						ready = false
						local temp = imgSrc
						imgSrc = imgEmpty
						gui.update()
						collection.collect(collection.IMAGE, temp)

						if ref.get(srcOption) == 0 then
							imgSrc = image.copy(imgOrigin, "src", image.MODEL_RGBA)
						elseif ref.get(srcOption) == 1 then
							imgSrc = imger.gradient_linear_xy(
								200,
								200,
								ref.get(srcGradientStart),
								ref.get(srcGradientEnd),
								ref.get(srcGradientDir),
								"src",
								image.ENCODING_PNG
							)
						elseif ref.get(srcOption) == 2 then
							imgSrc = imger.gradient_sigmoidal_xy(
								200,
								200,
								ref.get(srcGradientStart),
								ref.get(srcGradientEnd),
								ref.get(srcGradientDir),
								"src",
								image.ENCODING_PNG
							)
						elseif ref.get(srcOption) == 3 then
							imgSrc = noise.simplex_image_new(
								ref.get(srcNoiseSeed),
								ref.get(srcNoiseScale),
								ref.get(srcNormalize),
								"src",
								image.ENCODING_PNG,
								200,
								200,
								image.MODEL_RGBA,
								false,
								true
							)
						elseif ref.get(srcOption) == 4 then
							imgSrc = image.new_random("src", image.ENCODING_PNG, 200, 200, false, image.MODEL_RGBA)
						end

						collection.schedule(collection.IMAGE, imgSrc, function()
							ready = true
							gui.update()
						end)
					end),
					gui.wg_combo_preview("Image", {
						"Crystal",
						"Gradient Linear",
						"Gradient Sigmoidal",
						"Noise Map",
						"Random",
					}, srcOption),
				}),
				gui.wg_separator(),
				gui.wg_row({
					gui.wg_label("Gradient"),
					gui.wg_row({
						gui.wg_button_radio("Horizontal", ref.get(srcGradientDir) == 0):on_change(function()
							ref.set(srcGradientDir, 0)
						end),
						gui.wg_button_radio("Vertical", ref.get(srcGradientDir) == 1):on_change(function()
							ref.set(srcGradientDir, 1)
						end),
					}),
				}),
				gui.wg_color_edit("Start", srcGradientStart),
				gui.wg_color_edit("End", srcGradientEnd),
				gui.wg_separator(),
				gui.wg_row({
					gui.wg_label("Noise Map"),
					gui.wg_label("Seed:"),
					gui.wg_input_int(srcNoiseSeed):size(100),
					gui.wg_button("Randomize"):on_click(function()
						ref.set(srcNoiseSeed, math.random(0, 9999))
					end),
					gui.wg_checkbox("Normalize", srcNormalize),
				}),
				widget.slider_float("Coef", srcNoiseScale, 0.0, 1.0),
			}),
			wg_filter("Grayscale", {
				gui.wg_checkbox("Use GRAY over RGBA", grayscaleUse),
				gui.wg_row({
					gui.wg_checkbox("Grayscale Enabled", grayscaleEnabled),
					gui.wg_checkbox("To 16-bit", grayscaleTo16),
				}),
				gui.wg_checkbox("Add Self", grayscaleSelfEnabled),
				gui.wg_row({
					gui.wg_checkbox("Add Weighted", grayscaleWeightedEnabled),
					gui.wg_slider_float(grayscaleWeight1, 0, 1):size(125):label("Weight 1"),
					gui.wg_slider_float(grayscaleWeight2, 0, 1):size(125):label("Weight 2"),
					gui.wg_button_arrow(gui.DIR_LEFT):on_click(function()
						ref.set(grayscaleWeight1, 0.5)
						ref.set(grayscaleWeight2, 0.5)
					end),
				}),
				gui.wg_row({
					gui.wg_checkbox("Scalar", grayscaleScalarEnabled),
					gui.wg_dummy(
						gui.calc_text_size_width("Add Weighted") - gui.calc_text_size_width("Scalar") - spx,
						1
					),
					gui.wg_slider_int(grayscaleScalar, 0, 255):size(
						250
							+ gui.calc_text_size_width("Weight 1")
							+ gui.calc_text_size_width("Weight 2")
							+ (spx + (isx * 2))
					),
					gui.wg_button_arrow(gui.DIR_LEFT):on_click(function()
						ref.set(grayscaleScalar, 127)
					end),
				}),
			}),
			wg_filter("Transform", {
				gui.wg_row({
					gui.wg_checkbox("Rotate Enabled", rotateEnabled),
					gui.wg_checkbox("Resize To Fit", rotateResize),
					gui.wg_checkbox("Use GRAY", rotateGray),
				}),
				widget.slider_float_step("Angle", rotateAngle, 0.0, 360.0, 180.0),
				gui.wg_row({
					gui.wg_label("Anchor"),
					gui.wg_dummy(100 - gui.calc_text_size_width("Anchor"), 1),
					gui.wg_slider_int(rotateAX, 0, 200):size(125),
					gui.wg_slider_int(rotateAY, 0, 200):size(125),
					gui.wg_button_arrow(gui.DIR_LEFT):on_click(function()
						ref.set(rotateAX, 100)
						ref.set(rotateAY, 100)
					end),
				}),
				gui.wg_separator(),
				gui.wg_row({
					gui.wg_checkbox("Resize Enabled", resizeEnabled),
					gui.wg_checkbox("Use GRAY", resizeGray),
				}),
				widget.slider_float_step("Scale X", resizeScalex, 0.1, 5.0, 1.0),
				widget.slider_float_step("Scale Y", resizeScaley, 0.1, 5.0, 1.0),
				gui.wg_combo_preview("Interpolation", {
					"Nearest",
					"Linear",
					"Catmull-Rom",
					"Lanczos",
				}, resizeInterp),
			}),
			wg_filter("Blur", {
				gui.wg_row({
					gui.wg_checkbox("Box Blur Enabled", blurBoxEnabled),
					gui.wg_checkbox("Use GRAY", blurBoxGray),
				}),
				gui.wg_row({
					gui.wg_label("Size"),
					gui.wg_dummy(gui.calc_text_size_width("Anchor") - gui.calc_text_size_width("Size") - spx, 1),
					gui.wg_slider_int(blurBoxWidth, 1, 25):size(125),
					gui.wg_slider_int(blurBoxHeight, 1, 25):size(125),
					gui.wg_button_arrow(gui.DIR_LEFT):on_click(function()
						ref.set(blurBoxWidth, 3)
						ref.set(blurBoxHeight, 3)
					end),
				}),
				gui.wg_row({
					gui.wg_label("Anchor"),
					gui.wg_slider_int(blurBoxAX, 0, 24):size(125),
					gui.wg_slider_int(blurBoxAY, 0, 24):size(125),
					gui.wg_button_arrow(gui.DIR_LEFT):on_click(function()
						ref.set(blurBoxAX, 1)
						ref.set(blurBoxAY, 1)
					end),
					gui.wg_button_small("Center"):on_click(function()
						ref.set(blurBoxAX, ref.get(blurBoxWidth) / 2)
						ref.set(blurBoxAY, ref.get(blurBoxHeight) / 2)
					end),
				}),
				wg_border(blurBoxBorder),
				gui.wg_separator(),
				gui.wg_row({
					gui.wg_checkbox("Gaussian Blur Enabled", gaussEnabled),
					gui.wg_checkbox("Use GRAY", gaussGray),
				}),
				widget.slider_float("Radius", gaussRadius, 0.1, 25, 3.0),
				widget.slider_float("Sigma", gaussSigma, 0.1, 5, 1.0),
				wg_border(gaussBorder),
			}),
			wg_filter("Edge Detection", {
				gui.wg_row({
					gui.wg_checkbox("Canny Enabled", edgeCannyEnabled),
					gui.wg_checkbox("Use GRAY", edgeCannyGray),
				}),
				widget.slider_float("Lower Threshold", edgeCannyLower, 0.1, 1.0, 0.5),
				widget.slider_float("Upper Threshold", edgeCannyUpper, 0.1, 1.0, 1.0),
				widget.slider_int("Kernel Size", edgeCannyKernel, 3, 25, 3),
				gui.wg_separator(),
				gui.wg_row({
					gui.wg_checkbox("Sobel Enabled", edgeSobelEnabled),
					gui.wg_checkbox("Use GRAY", edgeSobelGray),
				}),
				gui.wg_row({
					gui.wg_button_radio("Both", ref.get(edgeSobelDir) == 0):on_change(function()
						ref.set(edgeSobelDir, 0)
					end),
					gui.wg_button_radio("Horizontal", ref.get(edgeSobelDir) == 1):on_change(function()
						ref.set(edgeSobelDir, 1)
					end),
					gui.wg_button_radio("Vertical", ref.get(edgeSobelDir) == 2):on_change(function()
						ref.set(edgeSobelDir, 2)
					end),
				}),
				wg_border(edgeSobelBorder),
				gui.wg_separator(),
				gui.wg_row({
					gui.wg_checkbox("Laplacian Enabled", edgeLaplacianEnabled),
					gui.wg_checkbox("Use GRAY", edgeLaplacianGray),
				}),
				gui.wg_combo_preview("Kernel", {
					"K4",
					"K8",
				}, edgeLaplacianKernel),
				wg_border(edgeLaplacianBorder),
			}),
			wg_filter("Effects", {
				gui.wg_checkbox("Sepia Enabled", effectSepiaEnabled),
				gui.wg_row({
					gui.wg_checkbox("Emboss Enabled", effectEmbossEnabled),
					gui.wg_checkbox("Use GRAY", effectEmbossGray),
				}),
				gui.wg_row({
					gui.wg_checkbox("Invert Enabled", effectInvertEnabled),
					gui.wg_checkbox("Use GRAY", effectInvertGray),
				}),
				gui.wg_row({
					gui.wg_checkbox("Pixelate Enabled", effectPixelateEnabled),
					gui.wg_checkbox("Use GRAY", effectPixelateGray),
				}),
				widget.slider_float_step("Factor", effectPixelateSize, 1.0, 15.0, 5.0),
				gui.wg_row({
					gui.wg_checkbox("Sharpen Enabled", effectSharpenEnabled),
					gui.wg_checkbox("Use GRAY", effectSharpenGray),
				}),
			}),
			wg_filter("Padding", {
				gui.wg_row({
					gui.wg_checkbox("Padding Enabled", paddingEnabled),
					gui.wg_button_radio("XY", ref.get(paddingMode) == 0):on_change(function()
						ref.set(paddingMode, 0)
					end),
					gui.wg_button_radio("Size", ref.get(paddingMode) == 1):on_change(function()
						ref.set(paddingMode, 1)
					end),
					gui.wg_checkbox("Use GRAY", paddingGray),
				}),
				wg_border(paddingBorder),
				gui.wg_tab_bar():tab_items({
					gui.wg_tab_item("XY"):layout({
						gui.wg_row({
							gui.wg_label("Size"),
							gui.wg_dummy(
								gui.calc_text_size_width("Anchor") - gui.calc_text_size_width("Size") - spx,
								1
							),
							gui.wg_slider_int(paddingKernelWidth, 1, 25):size(125),
							gui.wg_slider_int(paddingKernelHeight, 1, 25):size(125),
							gui.wg_button_arrow(gui.DIR_LEFT):on_click(function()
								ref.set(paddingKernelWidth, 3)
								ref.set(paddingKernelHeight, 3)
							end),
						}),
						gui.wg_row({
							gui.wg_label("Anchor"),
							gui.wg_slider_int(paddingKernelAX, 0, 24):size(125),
							gui.wg_slider_int(paddingKernelAY, 0, 24):size(125),
							gui.wg_button_arrow(gui.DIR_LEFT):on_click(function()
								ref.set(paddingKernelAX, 1)
								ref.set(paddingKernelAY, 1)
							end),
							gui.wg_button_small("Center"):on_click(function()
								ref.set(paddingKernelAX, ref.get(paddingKernelWidth) / 2)
								ref.set(paddingKernelAY, ref.get(paddingKernelHeight) / 2)
							end),
						}),
					}),
					gui.wg_tab_item("Size"):layout({
						widget.slider_int("Top", paddingTop, 0, 25, 1),
						widget.slider_int("Bottom", paddingBottom, 0, 25, 1),
						widget.slider_int("Left", paddingLeft, 0, 25, 1),
						widget.slider_int("Right", paddingRight, 0, 25, 1),
					}),
				}),
			}),
			wg_filter("Threshold", {
				gui.wg_row({
					gui.wg_checkbox("Threshold Enabled", thresholdEnabled),
					gui.wg_button_radio("Normal", ref.get(thresholdMode) == 0):on_change(function()
						ref.set(thresholdMode, 0)
					end),
					gui.wg_button_radio("Otsu", ref.get(thresholdMode) == 1):on_change(function()
						ref.set(thresholdMode, 1)
					end),
				}),
				gui.wg_combo_preview("Method", {
					"Binary",
					"Binary Inverted",
					"Truncate",
					"To Zero",
					"To Zero Inverted",
				}, thresholdMethod),
				gui.wg_separator(),
				gui.wg_label("Normal Only:"),
				gui.wg_checkbox("To 16-bit", thresholdTo16),
				widget.slider_int("Value", thresholdValue, 0, 255, 127),
			}),
		})
	end)
end

function wg_filter(name, widgets)
	return gui.wg_tree_node(name)
		:flags(bit.bitor_many(gui.FLAGTREENODE_FRAMED, gui.FLAGTREENODE_NOTREEPUSHONOPEN))
		:layout(widgets)
end

function wg_border(selectRef)
	return gui.wg_combo_preview("Border", {
		"Constant",
		"Replicate",
		"Reflect",
	}, selectRef)
end
