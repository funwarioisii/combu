package combu

import (
	"encoding/json"
	"github.com/google/go-jsonnet"
	"gopkg.in/yaml.v2"
	"math/rand"
	"regexp"
	"time"
)

func ReadConfig(fileBuf []byte) ([]DockerParam, error) {
	data := make([]DockerParam, 20)
	err := yaml.Unmarshal(fileBuf, &data)
	if err != nil {
		panic(err)
	}
	return data, nil
}

func randomText()  string {
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	text := make([]rune, 16)
	for i := range text {
		rand.Seed(time.Now().UnixNano())
		text[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(text)
}


func ReadJsonnetConfig(fileBuf []byte) (DockerParams, error) {
	snippet := string(fileBuf)

	rep := regexp.MustCompile("(['\"])UUID(['\"])")
	snippet = rep.ReplaceAllString(snippet, "\"" + randomText() + "\"")

	vm := jsonnet.MakeVM()

	rawOutput, err := vm.EvaluateSnippet("file", snippet)
	if err != nil {
		panic(err)
	}

	var dockerParams []DockerParam
	err = json.Unmarshal([]byte(rawOutput), &dockerParams)
	if err != nil {
		panic(err)
	}
	return dockerParams, nil
}
