local uuid = "UUID";

[
    {
        name: "mzld",
        image: "registry-sd.com/kazuyuki-hashimoto/mzld-image-run:latest",
        networks: ["irl"],
    }
]
  +
[
    {
        name: "irl-hard-%d" % [i_1],
        image: "registry-sd.com/kazuyuki-hashimoto/rl_navi-image-run:latest",
        networks: ["irl"],
        cmd: "make parallel-simulation MODE=hard MAZE_ID=%d" % [i_1] + " --uuid=" +uuid,
        depends: ["mzld"]
    }  for i_1 in std.range(2, 4)
]
  +
[
    {
        name: "irl-soft-%d" % [i_1],
        image: "registry-sd.com/kazuyuki-hashimoto/rl_navi-image-run:latest",
        networks: ["irl"],
        cmd: "make parallel-simulation MODE=soft MAZE_ID=%d" % [i_1] + " --uuid=" +uuid,
        depends: ["mzld"]
    }  for i_1 in std.range(2, 4)
]