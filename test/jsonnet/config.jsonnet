[
    {
        name: "nginx-%d" % [i],
        image: "nginx",
        ports: [
            {
                host:'%d/tcp' % [22220+i*2],
                container: '80/tcp'
            },
            {
                host: '%d/tcp' % [22220+i*2+1],
                container: '443/tcp'
            }
        ],
        networks: ["nginxies"]
    } for i in std.range(0, 3)  # include 3!
]