package test

import (
	"combu/combu"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestSort(t *testing.T) {
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{a}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{a}}
	d := combu.DepNode{Name: "D", Nodes: combu.DepNodes{b}}
	e := combu.DepNode{Name: "E", Nodes: combu.DepNodes{b, c}}
	f := combu.DepNode{Name: "F", Nodes: combu.DepNodes{c, e}}

	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c, d, e, f}}

	for i, node := range g.Sort() {
		println(i, node.String())
	}
}

func TestSort2(t *testing.T) {
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{a}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{a}}
	e := combu.DepNode{Name: "E", Nodes: combu.DepNodes{b, c}}
	f := combu.DepNode{Name: "F", Nodes: combu.DepNodes{c, e}}
	d := combu.DepNode{Name: "D", Nodes: combu.DepNodes{c, f}}

	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c, d, e, f}}

	for i, node := range g.Sort() {
		println(i, node.String())
	}
}

func TestSort3(t *testing.T) {
	/**
	1つに対していっぱいパターン
	*/
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{a}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{a}}
	d := combu.DepNode{Name: "E", Nodes: combu.DepNodes{a}}
	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c, d}}

	for i, node := range g.Sort() {
		println(i, node.String())
	}
}

func TestSort4(t *testing.T) {
	/**
	Nodeの宣言順序が雑パターン
	*/
	d := combu.DepNode{Name: "E", Nodes: combu.DepNodes{combu.DepNode{Name: "B"}}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{combu.DepNode{Name: "A"}}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{combu.DepNode{Name: "A"}}}
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}

	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c, d}}

	for i, node := range g.Sort() {
		println(i, node.String())
	}
}

func TestSort5(t *testing.T) {
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{a}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{a}}
	d := combu.DepNode{Name: "D", Nodes: combu.DepNodes{a}}
	e := combu.DepNode{Name: "E", Nodes: combu.DepNodes{a, b}}
	f := combu.DepNode{Name: "F", Nodes: combu.DepNodes{a, c}}

	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c, d, e, f}}

	for i, node := range g.Sort() {
		println(i, node.String())
	}
}

func TestSort6(t *testing.T) {
	/**
	依存nodeが存在しないパターン
	*/
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{}}

	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c}}

	for i, node := range g.Sort() {
		println(i, node.String())
	}
}

func TestResolveDep(t *testing.T) {
	buf, err := ioutil.ReadFile("complex-dep.jsonnet")
	if err != nil {
		panic(err)
	}

	params, err := combu.ReadJsonnetConfig(buf)
	if err != nil {
		panic(err)
	}

	params = params.ResolveDependency()

	fmt.Println("resolved")
	for _, param := range params {
		println(param.Name)
	}
}

func TestSearch(t *testing.T) {
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{a}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{a}}
	d := combu.DepNode{Name: "D", Nodes: combu.DepNodes{b}}
	e := combu.DepNode{Name: "E", Nodes: combu.DepNodes{b, c}}
	f := combu.DepNode{Name: "F", Nodes: combu.DepNodes{c, e}}

	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c, d, e, f}}

	node := g.SearchTop()[0]
	println("search 1", node.String())
	println("remove node")
	g = g.RemoveNode(node)
	for _, depNode := range g.Nodes {
		println("after removed", depNode.String())
	}

	for _, node := range g.SearchTop() {
		println("search 2", node.String())
	}
}

func TestSearch2(t *testing.T) {
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{}}
	d := combu.DepNode{Name: "D", Nodes: combu.DepNodes{b}}
	e := combu.DepNode{Name: "E", Nodes: combu.DepNodes{b, c}}
	f := combu.DepNode{Name: "F", Nodes: combu.DepNodes{c, e}}

	g := combu.DepGraph{Nodes: combu.DepNodes{b, c, d, e, f}}

	for _, node := range g.SearchTop() {
		println(node.String())
	}
}

func TestQueue(t *testing.T) {
	buf, err := ioutil.ReadFile("irl-uuid.jsonnet")
	if err != nil {
		panic(err)
	}

	params, err := combu.ReadJsonnetConfig(buf)
	if err != nil {
		panic(err)
	}

	var queue []combu.DockerParam
	for _, param := range params {
		queue = append(queue, param)
	}

	for _, param := range queue {
		for _, dep := range param.Depends {
			println(dep)
		}
	}
}

func TestNodeOperation(t *testing.T) {
	a := combu.DepNode{Name: "A"}
	b := combu.DepNode{Name: "B", Nodes: []combu.DepNode{a}}
	c := combu.DepNode{Name: "C", Nodes: []combu.DepNode{b}}
	d := combu.DepNode{Name: "D", Nodes: []combu.DepNode{b}}
	e := combu.DepNode{Name: "E", Nodes: []combu.DepNode{c}}

	println(d.String())
	println(e.String())
}

func TestGraphOperation(t *testing.T) {
	a := combu.DepNode{Name: "A"}
	b := combu.DepNode{Name: "B", Nodes: []combu.DepNode{a}}
	c := combu.DepNode{Name: "C", Nodes: []combu.DepNode{b}}
	d := combu.DepNode{Name: "D", Nodes: []combu.DepNode{b}}
	e := combu.DepNode{Name: "E", Nodes: []combu.DepNode{c}}

	g := combu.DepGraph{Nodes: []combu.DepNode{a, b, c, d, e}}
	println(g.Search("A").String())
}

func TestGraphGroup(t *testing.T) {
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{a}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{a}}
	d := combu.DepNode{Name: "D", Nodes: combu.DepNodes{b}}
	e := combu.DepNode{Name: "E", Nodes: combu.DepNodes{b, c}}
	f := combu.DepNode{Name: "F", Nodes: combu.DepNodes{c, e}}

	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c, d, e, f}}

	gr := g.Grouping()

	for i, nodes := range gr {
		for j, node := range nodes {
			println(i, j, node.String())
		}
	}
}

func TestGraphGroup2(t *testing.T) {
	/**

	 */
	a := combu.DepNode{Name: "A", Nodes: combu.DepNodes{}}
	b := combu.DepNode{Name: "B", Nodes: combu.DepNodes{a}}
	c := combu.DepNode{Name: "C", Nodes: combu.DepNodes{a}}
	d := combu.DepNode{Name: "D", Nodes: combu.DepNodes{a}}
	e := combu.DepNode{Name: "E", Nodes: combu.DepNodes{a}}
	f := combu.DepNode{Name: "F", Nodes: combu.DepNodes{a}}

	g := combu.DepGraph{Nodes: combu.DepNodes{a, b, c, d, e, f}}

	gr := g.Grouping()

	for i, nodes := range gr {
		for j, node := range nodes {
			println(i, j, node.String())
		}
	}
}
