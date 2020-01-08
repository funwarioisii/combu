package test

import (
	"combu/combu"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"testing"
	"time"
)

var containerId = "5c9f7c0482a0"

func initCombu() *client.Client {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
		return nil
	}
	return cli
}

func TestStats(t *testing.T) {
	cli := initCombu()
	c := struct {
		stream bool
		str    string
	}{
		false,
		"",
	}
	stats, _ := cli.ContainerStats(context.Background(), containerId, c.stream)
	defer stats.Body.Close()
	content, _ := ioutil.ReadAll(stats.Body)
	println(string(content))
}

// https://stackoverflow.com/questions/47154036/decode-json-from-stream-of-data-docker-go-sdk
type myStruct struct {
	Id       string `json:"id"`
	Read     string `json:"read"`
	Preread  string `json:"preread"`
	CpuStats cpu    `json:"cpu_stats"`
}

type cpu struct {
	Usage cpuUsage `json:"cpu_usage"`
}

type cpuUsage struct {
	Total float64 `json:"total_usage"`
}

func TestReadStatsJIT(t *testing.T) {
	cli := initCombu()
	stats, _ := cli.ContainerStats(context.Background(), containerId, false)
	//content, _ := ioutil.ReadAll(stats.Body)
	//println(string(content))

	decoder := json.NewDecoder(stats.Body)
	var containerStats combu.ContainerStats
	err := decoder.Decode(&containerStats)
	if err == io.EOF {
		fmt.Println("EOF")
		fmt.Println(containerStats.Name)
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("else")
		cpuTotal := float32(containerStats.CpuStats.CPUUsage.TotalUsage)
		preCpus := containerStats.CpuStats.CPUUsage.PercpuUsage
		for i, cpus := range preCpus {
			fmt.Printf("%d, %f, cpuTotal: %f\n", i, float32(cpus)/cpuTotal, cpuTotal)

		}
		cpuUsage := float32(containerStats.CpuStats.CPUUsage.TotalUsage) / (1024 * 1024 * 1024)
		fmt.Println(containerStats.Name, cpuUsage)
	}
}

func TestReadStatsStream(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cli, e := client.NewClientWithOpts(client.WithVersion("1.38"))
	if e != nil {
		panic(e)
	}
	stats, e := cli.ContainerStats(ctx, containerId, true)
	if e != nil {
		_ = fmt.Errorf("%s", e.Error())
	}
	decoder := json.NewDecoder(stats.Body)
	var containerStats myStruct
	for {
		select {
		case <-ctx.Done():
			e = stats.Body.Close()
			if e != nil {
				panic(e)
			}
			fmt.Println("Stop logging")
			return
		default:
			if err := decoder.Decode(&containerStats); err == io.EOF {
				return
			} else if err != nil {
				cancel()
			}
			fmt.Println(containerStats.CpuStats.Usage.Total / (1024 * 1024 * 1024))
		}
	}
}
