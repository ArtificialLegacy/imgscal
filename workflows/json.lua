
config({
    name= "JSON",
    version= "0.1.0",
    author="Blub",
    requires= {
        "cli",
        "json",
        "std"
    },

    desc="Converts all images in a dir to ascii art.",
})

main(function ()
    pth = cli.question("Enter directory to process")
    t = json.parse(pth)
    json.save(t, "./output/test.json")
end)