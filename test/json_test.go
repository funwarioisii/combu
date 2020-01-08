package test

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-jsonnet"
	"testing"
)

var snipped = `
local fibonacci(n) = if n <= 1 then 1
else fibonacci(n - 1) + fibonacci(n - 2);
local nodes = ["Node A", "Node B", "Node C"];
[
	{
		id: i,
		node: nodes[i],
		backoff: fibonacci(i+1)
	} for i in std.range(0, std.length(nodes)-1)
]`

type SampleConfig struct {
	Backoff int    `json:"backoff"`
	Id      int    `json:"id"`
	Node    string `json:"node"`
}

func TestExampleSuccess(t *testing.T) {
	vm := jsonnet.MakeVM()

	rawOutput, err := vm.EvaluateSnippet("file", snipped)
	if err != nil {
		panic(err)
	}
	//println(rawOutput)

	var configs []SampleConfig
	err = json.Unmarshal([]byte(rawOutput), &configs)

	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(configs)
}
