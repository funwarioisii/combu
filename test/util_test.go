package test

import (
	"combu/combu"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func TestReadConfigJsonnet(t *testing.T) {
	buf, err := ioutil.ReadFile("config.jsonnet")
	if err != nil {
		panic(err)
	}

	params, err := combu.ReadJsonnetConfig(buf)
	if err != nil {
		panic(err)
	}
	println(params[0].Image, params[0].Networks[0])
}

func TestReadConfigJsonnetIRL(t *testing.T) {
	buf, err := ioutil.ReadFile("irl-uuid.jsonnet")
	if err != nil {
		panic(err)
	}

	params, err := combu.ReadJsonnetConfig(buf)
	if err != nil {
		panic(err)
	}
	println(params[0].Image, params[0].Networks[0], params[0].Name)

	println(len(params))
}

func TestReadConfigJsonnetIncludeUUID(t *testing.T) {
	buf, err := ioutil.ReadFile("jsonnet/irl-uuid.jsonnet")
	if err != nil {
		panic(err)
	}

	snippet := string(buf)
	rep := regexp.MustCompile("('|\")UUID('|\")")
	snippet = rep.ReplaceAllString(snippet, "\"" + randomText() + "\"")
	println(snippet)

	params, err := combu.ReadJsonnetConfig(buf)
	if err != nil {
		panic(err)
	}
	println(params[2].Image, params[2].Networks[0], params[2].Cmd)

	println(len(params))
}




func randomText()  string {
	var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	text := make([]rune, 16)
	for i := range text {
		rand.Seed(time.Now().UnixNano())
		text[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(text)
}

func TestRandomText(t *testing.T) {
	println(randomText())
	println(randomText())
	println(randomText())
	println(randomText())

	fmt.Println(rand.Intn(100))
	fmt.Println(rand.Intn(100))
	fmt.Println(rand.Intn(100))
}