package combu

import (
	"github.com/docker/docker/api/types"
)

type VolumeParam struct {
	Host      string `json:"host"`
	Container string `json:"container"`
}

type PortParam struct {
	Host      string `json:"host"`
	Container string `json:"container"`
}

type DockerParam struct {
	Name     string        `json:"name"`
	Image    string        `json:"image"`
	Ports    []PortParam   `json:"ports"`
	Volumes  []VolumeParam `json:"volumes"`
	Networks []string      `json:"networks"`
	Depends  []string      `json:"depends"`
	Cmd      string        `json:"cmd"`
}

func (p DockerParam) String() string {
	text := "Name:"
	text += p.Name
	text += ", nodes["

	for _, dep := range p.Depends {
		text += dep
		text += ","
	}
	text += "]"
	return text
}

type DockerParams []DockerParam

func (params DockerParams) String() string {
	text := "["

	for _, p := range params {
		text += p.String()
		text += ", "
	}

	text += "]"
	return text
}

func (params DockerParams) ResolveDependency() DockerParams {
	m := make(map[string]DockerParam)

	nodes := DepNodes{}
	for _, param := range params {
		deps := DepNodes{}

		for _, dep := range param.Depends {
			deps = append(deps, DepNode{Name: dep})
		}
		nodes = append(nodes, DepNode{
			Name:  param.Name,
			Nodes: deps,
		})

		m[param.Name] = param
	}

	g := DepGraph{Nodes: nodes}
	nodes = g.Sort()

	var sortedParams []DockerParam
	for _, node := range nodes {
		sortedParams = append(sortedParams, m[node.Name])
	}

	return sortedParams
}

func (params DockerParams) Grouping() DockerGroup {
	m := make(map[string]DockerParam)

	nodes := DepNodes{}
	for _, param := range params {
		deps := DepNodes{}

		for _, dep := range param.Depends {
			deps = append(deps, DepNode{Name: dep})
		}
		nodes = append(nodes, DepNode{
			Name:  param.Name,
			Nodes: deps,
		})

		m[param.Name] = param
	}

	g := DepGraph{Nodes: nodes}
	gr := g.Grouping()

	var dg = DockerGroup{}

	for _, nodes := range gr {
		var params = DockerParams{}

		for _, node := range nodes {
			params = append(params, m[node.Name])
		}

		if params.String() != "[]" {
			dg = append(dg, params)
		}

	}

	return dg
}

type DockerGroup []DockerParams

type DepNode struct {
	Name  string
	Nodes DepNodes
}

func (n DepNode) Equal(node DepNode) bool {
	return n.Name == node.Name
}

func (n DepNode) String() string {
	text := "Name:"
	text += n.Name
	text += ", nodes["

	for _, node := range n.Nodes {
		text += node.Name
		text += ","
	}
	text += "]"
	return text
}

type DepNodes []DepNode

func (ns DepNodes) Equal(nodes DepNodes) bool {
	for x := range ns {
		for y := range nodes {
			if x != y {
				return false
			}
		}
	}

	return true
}

func (ns DepNodes) Contains(node DepNode) bool {
	for _, n := range ns {
		if n.Name == node.Name {
			return true
		}
	}

	return false
}

func (ns DepNodes) Remove(node DepNode) DepNodes {
	var depNodes DepNodes
	for _, n := range ns {
		if !n.Equal(node) {
			depNodes = append(depNodes, n)
		}
	}

	return depNodes
}

type DepGraph struct {
	Nodes DepNodes
}

func (g DepGraph) Search(name string) DepNode {
	for _, n := range g.Nodes {
		if n.Name == name {
			return n
		}
	}
	return DepNode{
		Name:  "",
		Nodes: nil,
	}
}

func (g DepGraph) sort(sorted DepNodes) (DepGraph, DepNodes) {
	nodes := g.SearchTop()

	for _, node := range nodes {
		g = g.RemoveNode(node)
		sorted = append(sorted, node)
	}
	return g, sorted
}

func (g DepGraph) Sort() DepNodes {
	var sorted DepNodes

	for range g.Nodes {
		g, sorted = g.sort(sorted)
	}
	return sorted
}

func (g DepGraph) SearchTop() DepNodes {
	/**
	1. 依存関係リスト
	2.
	*/
	var depNodes DepNodes
	for _, n := range g.Nodes {
		if len(n.Nodes) == 0 {
			depNodes = append(depNodes, n)
		}
	}
	return depNodes
}

func (g DepGraph) RemoveNode(node DepNode) DepGraph {
	g.Nodes = g.Nodes.Remove(node)

	var nodes DepNodes
	for _, depNode := range g.Nodes {
		if depNode.Nodes.Contains(node) {
			depNode.Nodes = depNode.Nodes.Remove(node)
		}
		nodes = append(nodes, depNode)
	}
	g.Nodes = nodes

	return g
}

func (g DepGraph) grouping(gr DepGroup) (DepGraph, DepGroup) {
	nodes := g.SearchTop()
	gr = append(gr, nodes)

	sorted := DepNodes{}

	for _, node := range nodes {
		g = g.RemoveNode(node)
		sorted = append(sorted, node)
	}

	return g, gr
}

func (g DepGraph) Grouping() DepGroup {
	gr := DepGroup{}

	for range g.Nodes {
		g, gr = g.grouping(gr)
	}
	return gr
}

type DepGroup []DepNodes

type NetworksStats map[string]types.NetworkStats

type ContainerStats struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Read        string           `json:"read"`
	Preread     string           `json:"preread"`
	NumProcs    uint             `json:"num_procs"`
	CpuStats    types.CPUStats   `json:"cpu_stats"`
	PreCpuStats types.CPUStats   `json:"pre_cpu_stats"`
	BlkioStats  types.BlkioStats `json:"blkio_stats"`

	PidsStats types.PidsStats `json:"pids_stats"`
	Networks  NetworksStats   `json:"networks"`
}
