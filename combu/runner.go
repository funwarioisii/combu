package combu

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
	"strings"
	"sync"
)

var hostIp = "0.0.0.0"

func initCombu() *client.Client {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
		return nil
	}
	return cli
}

func getMounts(volParam []VolumeParam) []mount.Mount {
	var mounts []mount.Mount
	for _, vol := range volParam {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: vol.Host,
			Target: vol.Container,
		})
	}
	return mounts
}

func getPortsConfig(hostIp string, netParam []PortParam) (nat.PortMap, nat.PortSet) {
	pBinds := make(nat.PortMap)
	exPorts := make(nat.PortSet)
	for _, por := range netParam {
		pBinds[nat.Port(por.Container)] = []nat.PortBinding{
			{
				HostPort: por.Host,
				HostIP:   hostIp,
			},
		}
		exPorts[nat.Port(por.Host)] = struct{}{}
	}

	for key, pbs := range pBinds {
		for _, pb := range pbs {
			log.Println(key, pb.HostPort, pb.HostIP)
		}
	}
	return pBinds, exPorts
}

func getContainerConfig(cmdParam string, imgParam string, exPorts nat.PortSet) *container.Config {
	cmd := strings.Split(cmdParam, " ")
	log.Println(cmd)

	var cc *container.Config
	if cmd[0] == "" {
		cc = &container.Config{Image: imgParam, ExposedPorts: exPorts}
	} else {
		cc = &container.Config{Image: imgParam, ExposedPorts: exPorts, Cmd: cmd}
	}
	return cc
}

func Kill(params []DockerParam) {
	cli := initCombu()
	containers, _ := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Quiet:   true,
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filters.Args{},
	})
	//log.Println(params)

	for _, param := range params {
		for _, cont := range containers {
			for _, name := range cont.Names {
				if "/"+param.Name == name {
					err := cli.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{Force: true})
					if err != nil {
						log.Println("failed remove", name)
					} else {
						log.Println("remove", name)
					}
				}
			}
		}
	}
}

func createNetwork(param DockerParam) {
	cli := initCombu()
	for _, net := range param.Networks {
		log.Println("create network...", net)

		_, _ = cli.NetworkCreate(context.Background(), net, types.NetworkCreate{
			CheckDuplicate: true,
			Driver:         "",
			Scope:          "",
			EnableIPv6:     false,
			IPAM:           nil,
			Internal:       false,
			Attachable:     false,
			Ingress:        false,
			ConfigOnly:     false,
			ConfigFrom:     nil,
			Options:        nil,
			Labels:         nil,
		})
	}
}

func createContainer(cli *client.Client, param DockerParam) container.ContainerCreateCreatedBody {
	rp, err := cli.ImagePull(context.Background(), param.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	createNetwork(param)

	payload := struct {
		ID             string `json:"id"`
		Status         string `json:"status"`
		Progress       string `json:"progress"`
		ProgressDetail struct {
			Current uint16 `json:"current"`
			Total   uint16 `json:"total"`
		} `json:"progressDetail"`
	}{}

	scanner := bufio.NewScanner(rp)
	for scanner.Scan() {
		_ = json.Unmarshal(scanner.Bytes(), &payload)
		//fmt.Printf("\r\t%+v\n", payload)
	}

	mounts := getMounts(param.Volumes)
	pBinds, exPorts := getPortsConfig(hostIp, param.Ports)

	cc := getContainerConfig(param.Cmd, param.Image, exPorts)

	hc := &container.HostConfig{
		//AutoRemove:   true,
		Mounts:       mounts,
		PortBindings: pBinds,
	}

	ec := map[string]*network.EndpointSettings{
		param.Networks[0]: {
			IPAMConfig: &network.EndpointIPAMConfig{
				IPv4Address:  "",
				IPv6Address:  "",
				LinkLocalIPs: nil,
			},
		},
	}

	nc := &network.NetworkingConfig{
		EndpointsConfig: ec,
	}
	con, err := cli.ContainerCreate(context.Background(), cc, hc, nc, param.Name)
	for warn := range con.Warnings {
		log.Println(warn)
	}
	if err != nil {
		panic(err)
		return container.ContainerCreateCreatedBody{}
	}
	log.Printf("%v", con.ID) // ContainerID

	return con
}

func startContainer(cli *client.Client, con container.ContainerCreateCreatedBody) error {
	err := cli.ContainerStart(context.Background(), con.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func runContainer(cli *client.Client, p DockerParam) error {
	con := createContainer(cli, p)
	return startContainer(cli, con)
}

func Run(params DockerParams) {
	cli := initCombu()

	params = params.ResolveDependency()

	for _, param := range params {
		log.Printf("create: %v\n", param.Name)
		err := runContainer(cli, param)
		if err != nil {
			panic(err)
		}
	}
}

func runParallelWithoutOrder(params DockerParams) {
	cli := initCombu()

	var wg sync.WaitGroup
	for _, param := range params {
		wg.Add(1)
		go func(c *client.Client, p DockerParam) {
			defer wg.Done()
			log.Printf("create: %v\n", p.Name)
			_ = runContainer(c, p)
		}(cli, param)
	}

	wg.Wait()
}

func RunParallel(param DockerParams) {
	g := param.Grouping()
	for _, params := range g {
		runParallelWithoutOrder(params)
	}
}
