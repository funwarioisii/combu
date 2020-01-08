[
    {
        name: "mzld",
        image: "registry-sd.com/kazuyuki-hashimoto/mzld-image-run",
        networks: ["irl"],
    }
]
  +
[
    {
        name: "irl-hard-%d" % [i_1],
        image: "registry-sd.com/kazuyuki-hashimoto/rl_navi-image-run",
        networks: ["irl"],
        cmd: "make simulate MODE=hard MAZE_ID=%d" % [i_1],
        depends: ["mzld"]
    }  for i_1 in std.range(0, 3)
]
  +
[
    {
        name: "irl-soft-%d" % [i_1],
        image: "registry-sd.com/kazuyuki-hashimoto/rl_navi-image-run",
        networks: ["irl"],
        cmd: "make simulate MODE=soft MAZE_ID=%d" % [i_1],
        depends: ["mzld"] + ["irl-hard-%d" % [i_1]]
    }  for i_1 in std.range(0, 3)
]