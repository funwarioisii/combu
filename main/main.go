package main

import (
	"combu/combu"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	var filename string
	var parallel string
	flag.StringVar(&filename, "f", "default", "Message")
	flag.StringVar(&parallel, "p", "true", "parallel")
	flag.Parse()

	args := flag.Args()
	cmd := args[0]

	log.Println(cmd)
	log.Println(filename)

	if filename == "" || filename == "default" {
		log.Println(filename)
		log.Fatal("specify config filename")
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := combu.ReadJsonnetConfig(buf)
	if err != nil {
		panic(err)
		return
	}

	if cmd == "run" {
		if parallel == "true" {
			combu.RunParallel(data)
		} else {
			combu.Run(data)
		}
	} else if cmd == "kill" {
		combu.Kill(data)
	}

}
