---@param workflow imgscal_WorkflowInit
function init(workflow)
	workflow.import({
		"image",
		"shader",
		"io",
		"cli",
	})
end

function main()
	local img = image.new_filled("image", image.ENCODING_PNG, 8, 8, image.color_rgb(200, 200, 200))

	local device = shader.device()
	shader.add_program_file(device, io.path_join(io.wd(), "shader.cl"))

	local buff_src = shader.buffer_image_from(device, img)
	local buff_dest = shader.buffer_image(device, shader.IMAGE_RGBA, 8, 8)

	local kernel = shader.kernel(device, "invert", { 8, 8 }, { 1, 1 })

	shader.run(device, kernel, buff_src, buff_dest)

	local result = buff_dest.data("result", image.ENCODING_PNG)

	cli.println("before:")
	cli.print_image(img, true)
	cli.println("after:")
	cli.print_image(result, true)
end
