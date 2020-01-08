package test

import (
	"combu/combu"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

func runContainer(param combu.DockerParam) {
	time.Sleep(1 * time.Second)
	fmt.Println("container create...", param.Name)
}

func runContainers2(params combu.DockerParams) {
	var wg sync.WaitGroup
	for i, param := range params {
		wg.Add(1)
		go func(p combu.DockerParam) {
			defer wg.Done()
			fmt.Println("create:", p.Name, i)
			runContainer(p)
		}(param)
	}

	wg.Wait()
	return
}

func TestParallelize(t *testing.T) {
	dParams := combu.DockerParams{
		combu.DockerParam{Name: "containerId1"},
		combu.DockerParam{Name: "containerId2"},
		combu.DockerParam{Name: "containerId3"},
		combu.DockerParam{Name: "containerId4"},
		combu.DockerParam{Name: "containerId5"},
		combu.DockerParam{Name: "containerId6"},
	}

	//c := make(chan int, len(params))
	//
	runContainers2(dParams)
}

func TestParallelize2(t *testing.T) {
	dp := combu.DockerParams{
		combu.DockerParam{Name: "a"},
		combu.DockerParam{Name: "b", Depends: []string{"a"}},
		combu.DockerParam{Name: "c", Depends: []string{"a"}},
		combu.DockerParam{Name: "d", Depends: []string{"a"}},
		combu.DockerParam{Name: "e"},
		combu.DockerParam{Name: "f"},
	}

	dg := dp.Grouping()
	for i, nodes := range dg {
		for j, node := range nodes {
			println(i, j, node.String())
		}
	}
}

func TestParallelize3(t *testing.T) {
	dp := combu.DockerParams{
		combu.DockerParam{Name: "a"},
		combu.DockerParam{Name: "b", Depends: []string{"a"}},
		combu.DockerParam{Name: "c", Depends: []string{"a"}},
		combu.DockerParam{Name: "d", Depends: []string{"a"}},
		combu.DockerParam{Name: "e"},
		combu.DockerParam{Name: "f"},
	}

	dg := dp.Grouping()

	for i, params := range dg {
		runContainers2(params)
		println(i, "done")

		if i > 1 {
			for _, p := range params {
				print(p.String())
			}
		}
	}

	for i, ps := range dg {
		println(i, ps.String())
	}
}

func TestParallelize4(t *testing.T) {
	filename := "../test/hard-sample.jsonnet"
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	dp, err := combu.ReadJsonnetConfig(buf)
	if err != nil {
		panic(err)
		return
	}

	dg := dp.Grouping()

	for i, params := range dg {
		runContainers2(params)
		println(i, "done")

		if i > 1 {
			for _, p := range params {
				print(p.String())
			}
		}
	}

	for i, ps := range dg {
		println(i, ps.String())
	}
}
