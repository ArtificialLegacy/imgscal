---@diagnostic disable:lowercase-global

function init(workflow)
	workflow.import({
		"gamemaker",
		"image",
		"std",
	})

	workflow.secrets({
		project_directory = "",
	})
end

function main()
	local secrets = std.secrets()
	if secrets.project_directory == "" then
		std.panic("Please set the 'project_directory' secret to the path of the GameMaker project.")
	end

	local proj = gamemaker.project_load(secrets.project_directory)

	local img1Small =
		image.new_filled("img1_small", image.ENCODING_PNG, 12, 12, image.color_rgb(255, 0, 0), image.MODEL_NRGBA)
	local img1 = image.new("img1", image.ENCODING_PNG, 16, 16, image.MODEL_NRGBA)
	image.draw(img1, img1Small, 2, 2)

	local img2 = image.new_filled("img2", image.ENCODING_PNG, 16, 16, image.color_rgb(255, 255, 255), image.MODEL_NRGBA)

	local sprite = gamemaker
		.sprite("sprImgScal", 16, 16, gamemaker.project_as_parent(proj), gamemaker.texgroup_default())
		:tags("test 1", "test 2", "test 3")
		:tile(false, true)
		:origin(gamemaker.SPRITEORIGIN_CUSTOM, 4, 4)
		:collision(gamemaker.BBOXMODE_MANUAL, gamemaker.COLLMASK_RECT, gamemaker.bbox(2, 2, 13, 13))
		:premultiply_alpha(true)
		:edge_filtering(true)
		:dynamic_texturepage(false)
		:nineslice(1, 1, 1, 1)
		:nineslice_tilemode(gamemaker.NINESLICESLICE_CENTER, gamemaker.NINESLICETILE_REPEAT)
		:nineslice_tilemode(gamemaker.NINESLICESLICE_LEFT, gamemaker.NINESLICETILE_MIRROR)
		:broadcast_message(0, "Hello World!")
		:broadcast_message(1, "Goodbye World!")
		:playback(5)
		:layers()
		:folder("Folder 1")
		:image("Layer Top")
		:back()
		:image("Layer Bottom")
		:back()
		:frames()
		:add({ img1, img2 })
		:add({ img1, img2 })
		:back()

	gamemaker.sprite_save(proj, sprite)
	gamemaker.project_save(proj)
end
