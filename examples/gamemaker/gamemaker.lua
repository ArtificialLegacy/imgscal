
function init(workflow)
    workflow.import({
        "gamemaker",
        "image",
        "io",
    })
end

dir = "/home/joseph/dev/gm-proj-tool-testing"

function main()
    local proj = gamemaker.project_load(dir)

    local img = image.new("datafile", image.ENCODING_PNG, 16, 16)

    local datafile1 = gamemaker.datafile_from_string("data.txt", "datafiles", "hello world!")
    local datafile2 = gamemaker.datafile_from_file("filedata.txt", "datafiles", io.path_join({io.wd(), "filedata.txt"}))
    local datafile3 = gamemaker.datafile_from_image("data.png", "datafiles", img)
 
    gamemaker.datafile_save(proj, datafile1)
    gamemaker.datafile_save(proj, datafile2)
    gamemaker.datafile_save(proj, datafile3)
    gamemaker.project_save(proj)
end