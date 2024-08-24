
function init(workflow)
    workflow.import({
        "gamemaker",
        "io",
        "image",
        "context",
    })
end

dir = "/home/joseph/dev/gm-proj-tool-testing"

function main()
    local proj = gamemaker.project_load(dir)

    local ctx = context.new(16, 16)

    context.draw_rect(ctx, 1, 1, 14, 14)
    context.color_rgb255(ctx, 255, 0, 0)
    context.fill(ctx)
    local img1 = context.to_image(ctx, "img1", image.ENCODING_PNG, image.MODEL_NRGBA, true)

    context.draw_rect(ctx, 0, 0, 16, 16)
    context.color_rgb255(ctx, 255, 255, 255)
    context.fill(ctx)
    local img2 = context.to_image(ctx, "img2", image.ENCODING_PNG, image.MODEL_NRGBA, true)

    local sprite = gamemaker.sprite("sprImgScal", 16, 16, gamemaker.project_as_parent(proj), gamemaker.texgroup_default())
        :layers()
            :folder("Folder 1")
                :image("Layer Top")
                :back()
            :image("Layer Bottom")
            :back()
        :frames()
            :add({img1, img2})
            :add({img1, img2})
            :back()
    
    gamemaker.sprite_save(proj, sprite)
    gamemaker.project_save(proj)
end