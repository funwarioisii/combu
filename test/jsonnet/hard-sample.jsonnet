local uuid = "UUID";

local solver(id) = {
    name: "busybox-b-%d" % [id],
    image: "busybox",
    cmd: "echo %s" % [uuid],
    networks: ["sample"],
    depends: ["busybox-a"]
};

[
    solver(_i)
    for _i in std.range(2, 20)
] + [
    {
        name: "busybox-a",
        image: "busybox",
        networks: ["sample"],
    },
    {
        name: "busybox-c",
        image: "busybox",
        networks: ["sample"],
        depends: ["busybox-b-10"]
    }
]