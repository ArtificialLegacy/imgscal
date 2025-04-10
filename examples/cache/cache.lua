---@param workflow imgscal_WorkflowInit
function init(workflow)
	workflow.import({
		"image",
		"cache",
	})
end

function main()
	local img1 = image.new_empty("image1", image.ENCODING_PNG)

	-- store image in cache and remove original image.
	local cached_img1 = cache.store(img1)
	image.remove(img1)

	-- retrieve image from cache
	local img2 = image.new_empty("image2", image.ENCODING_PNG)
	_ = cache.retrieve(cached_img1, img2)

	-- create cached image directly
	local cached_img2 = cache.new_empty()
	local img3 = image.new_empty("dummy", image.ENCODING_JPEG)
	_ = cache.retrieve_ext(cached_img2, img3, "image3", image.ENCODING_PNG, true)
end
