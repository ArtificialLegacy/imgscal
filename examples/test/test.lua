
function init(workflow)
    workflow.import({
        "cli",
        "image",
        "imger",
        "noise",
    })
end

function main()
    local img = noise.simplex_image_new(
        0, 0.1, true,
        "test", image.ENCODING_PNG,
        20, 20,
        image.MODEL_RGBA,
        false, true
    )

    local img2 = imger.edge_laplacian(
        img, imger.LAPLACIAN_K8,
        "test_edge", image.ENCODING_PNG,
        imger.BORDER_REFLECT
    )

    imger.threshold(img2, 5, imger.THRESHOLD_BINARY)
    imger.invert(img2, true)

    cli.print_image(img, true)
    cli.print("")
    cli.print_image(img2, true)
end