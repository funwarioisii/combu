# combu

## About

The `combu` is named as docker-**com**pose **bu**rst.

This is an easy container orchestration tool powered by jsonnet.

`combu run`; You can create and run many containers with one command.

This target user is amateur machine-learning developer and run experiments on docker container.

 
## Getting started

```terminal
$ make 
$ ./build/combu -file config/config.jsonnet run
```


## Usage

### create and start containers.

`combu -f <filename> run`

### destroy containers.

Attention:  It will remove all specified containers.

`combu -f <filename> kill`

### config file

This is a most simple configuration.

```jsonnet
[
    {
        name: "sample-container-a",
        image: "busybox",
        networks: ["sample"],
        ports: [
            {
                host:'2222/tcp',
                container: '80/tcp'
            },
        ],
    },
    {
        name: "sample-container-b",
        image: "busybox",
        depends: ["sample-container-a"],
        networks: ["sample"],
        cmd: "echo combu"
    }           
]
```

When started, 
1. create network `sample`
2. create and start `sample-container-a`
3. create and start `sample-container-b`
4. 

You have to fill `name` and `image`.

|key|explain|
|---|---|
|name|container name|
|image|docker image|
|ports|supply port by host and expected port by container|
|networks|docker network name|
|depends| depended container|
|cmd|override command|


<details>
<summary>This is a practical sample</summary>

In many cases, we need to use an unique id for experiments.

When you declare uuid, such as `local uuid ="UUID"`, `combu` set UUID with 12 random character.

```jsonnet
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
        ports: [
            {
                host:'3000/tcp',
                container: '3000/tcp'
            },
        ],
        networks: ["sample"],
    },
    {
        name: "busybox-c",
        image: "busybox",
        networks: ["sample"],
        depends: ["busybox-b-10", "busybox-b-20"],
    }
]
```

read the official [tutorial](https://jsonnet.org/learning/tutorial.html) of jsonnet, very helpful. 
</details>


## Development

- Go (>= 1.11)
- Docker (api>=1.38)

and some packages...(i don't know how to show like requirements.txt)


## Tips and my operation
### before run and kill

Before `combu -f <filename> run`, you have better to check with `jsonnet` command.

### set config your experiment projects
I put jsonnet files under `<project directory>/config/`.
