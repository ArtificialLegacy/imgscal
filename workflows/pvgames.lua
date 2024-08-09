config({
    name="PVGames",
    version="0.1.0",
    author="Blub",
    requires={
        "io",
        "spritesheet",
        "json",
        "image",
        "cli",
        "std"
    },

    desc="Extracts all animations from pvgames assets"
})

portraits = {
    "Neutral_1",
    "Neutral_2",
    "Happy_1",
    "Happy_2",
    "Sad_1",
    "Sad_2",
    "Angry_1",
    "Angry_2",
    "Nervous_1",
    "Nervous_2",
    "Scared",
    "Injured",
    "Thoughtful_1",
    "Thoughtful_2",
    "Annoyed_1",
    "Annoyed_2",
}

anims = {
    {
        name="walk",
        l=64,
    },
    {
        name="run",
        l=64,
    },
    {
        name="idle_1",
        l=40,
    },
    {
        name="idle_2",
        l=40,
    },
    {
        name="idle_3",
        l=40,
    },
    {
        name="idle_4",
        l=40,
    },
    {
        name="fidget_1",
        l=24,
    },
    {
        name="fidget_2",
        l=24,
    },
    {
        name="fidget_3",
        l=24,
    },
    {
        name="talking_1",
        l=40,
    },
    {
        name="talking_2",
        l=40,
    },
    {
        name="interact",
        l=40,
    },
    {
        name="use_item",
        l=24,
    },
    {
        name="sitting",
        l=24,
    },
    {
        name="climb",
        l=40,
    },
    {
        name="praying",
        l=24,
    },
    {
        name="jump",
        l=40,
    },
    {
        name="sneaking",
        l=64,
    },
    {
        name="crouch",
        l=24,
    },
    {
        name="casting",
        l=24,
    },
    {
        name="die_forward",
        l=40,
    },
    {
        name="die_backward",
        l=40,
    },
    {
        name="down",
        l=32,
    },
    {
        name="roll",
        l=40,
    },
    {
        name="hit_1",
        l=24,
    },
    {
        name="hit_2",
        l=24,
    },
    {
        name="idle_critical_1",
        l=24,
    },
    {
        name="idle_critical_2",
        l=24,
    },
    {
        name="block",
        l=24,
    },
    {
        name="drink",
        l=24,
    },
    {
        name="riding",
        l=24,
    },
    {
        name="walk_1H",
        l=64,
    },
    {
        name="idle_1H",
        l=40,
    },
    {
        name="attack_1H_1",
        l=24,
    },
    {
        name="attack_1H_2",
        l=24,
    },
    {
        name="attack_1H_3",
        l=24,
    },
    {
        name="fidget_1H",
        l=24,
    },
    {
        name="walk_2H",
        l=64,
    },
    {
        name="idle_2H",
        l=40,
    },
    {
        name="run_2H",
        l=64,
    },
    {
        name="attack_2H_1",
        l=24,
    },
    {
        name="attack_2H_2",
        l=24,
    },
    {
        name="walk_dual",
        l=64,
    },
    {
        name="idle_dual",
        l=40,
    },
    {
        name="attack_dual_1",
        l=24,
    },
    {
        name="attack_dual_2",
        l=24,
    },
    {
        name="fidget_dual",
        l=24,
    },
    {
        name="walk_bow",
        l=64,
    },
    {
        name="idle_bow",
        l=40,
    },
    {
        name="attack_bow",
        l=40,
    },
    {
        name="fidget_bow",
        l=24,
    },
    {
        name="idle_unarmed",
        l=40,
    },
    {
        name="attack_unarmed_1",
        l=24,
    },
    {
        name="attack_unarmed_2",
        l=24,
    },
    {
        name="fidget_unarmed",
        l=24,
    },
    {
        name="walk_staff",
        l=64,
    },
    {
        name="run_staff",
        l=64,
    },
    {
        name="idle_staff",
        l=40,
    },
    {
        name="attack_staff_1",
        l=24,
    },
    {
        name="attack_staff_2",
        l=24,
    },
    {
        name="walk_pistol",
        l=64,
    },
    {
        name="idle_pistol",
        l=40,
    },
    {
        name="attack_pistol_1",
        l=24,
    },
    {
        name="attack_pistol_2",
        l=24,
    },
    {
        name="fidget_pistol",
        l=24,
    },
    {
        name="walk_rifle",
        l=64,
    },
    {
        name="idle_rifle",
        l=40,
    },
    {
        name="attack_rifle",
        l=24,
    },
    {
        name="fidget_rifle",
        l=24,
    },
}

directions = {
    "s",
    "w",
    "e",
    "n",
    "sw",
    "nw",
    "se",
    "ne"
}

downPoses = {
    "sitting",
    "stomach",
    "back",
    "prepared"
}

function parsePortraits(body)
    local p = io.load_image(body.path.."/Portraits.png", image.MODEL_RGBA)
    local ps = spritesheet.to_frames(p, "portrait_"..body.name, 16, 2000, 2000, 4)

    for k2,v2 in pairs(ps) do
        io.mkdir("./output/bodies/"..body.name, true)
        image.name(v2, "portrait_"..portraits[k2])
        io.out(v2,"./output/bodies/"..body.name)
    end
end

function parseAnims(body)
    local s = io.load_image(body.path.."/Spritesheet.png", image.MODEL_RGBA)
    local ss = spritesheet.to_frames(s, "portrait_"..body.name, 2496, 200, 200, 50)

    local pos=1
    for k,v in pairs(anims) do
        local a = sliceTable(ss, pos, pos+v.l-1)
        pos = pos + v.l

        if v.name == "down" then
            for ind,p in pairs(downPoses) do
                for d,dv in pairs(directions) do
                    local index = ((#directions) * (ind-1)) + d
                    image.name(a[index], "down"..p.."_"..dv)
                    io.out(a[index], "./output/bodies/"..body.name)
                end
            end
        else
            local dpos = 1
            local per = (#a) / 8
            for d,dv in pairs(directions) do
                local dir = sliceTable(a, dpos, dpos+per-1)
                local an = spritesheet.from_frames(dir, v.name.."_"..dv, 200, 200, image.MODEL_RGBA, image.ENCODING_PNG, per)
                io.out(an, "./output/bodies/"..body.name)
            end
        end
    end
end

function sliceTable(tab, s, e)
    local new = {}
    local pos = 1

    for i=s, e do
        new[pos] = tab[i]
        pos = pos + 1
    end

    return new
end

main(function()
    local pth = cli.question("JSON path for image data.")
    local data = json.parse(pth)

    if data.bodies ~= nil then
        for k,v in pairs(data.bodies) do
            parsePortraits(v)
            parseAnims(v)
        end
    end
end)