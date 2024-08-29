
function init(workflow)
    workflow.import({
        "gamemaker",
        "image",
        "cli",
    })
end

dir = "/home/joseph/dev/gm-proj-tool-testing"

function main()
    local proj = gamemaker.project_load(dir)

    local sprite = gamemaker.sprite_load(proj, "sprImgScal", image.ENCODING_PNG)
    cli.print(sprite.name)
    cli.print_number(sprite.__layerCount)
    cli.print_number(sprite.width, true)
    cli.print_number(sprite.height, true)

    sprite.name = "sprImgScal2"

    gamemaker.sprite_save(proj, sprite)
    gamemaker.project_save(proj)
end