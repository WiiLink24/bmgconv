package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		log.Println("Usage: bmgconv [toXML|toBMG] <input> <output>")
		os.Exit(1)
	}

	action := os.Args[1]
	input := os.Args[2]
	output := os.Args[3]

	inputData, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}

	var outputData []byte

	switch action {
	case "toXML":
		outputData, err = parseBMG(inputData)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(output, outputData, 0600)
		if err != nil {
			panic(err)
		}
	case "toBMG":
		err = createBMG(inputData, output)
		if err != nil {
			panic(err)
		}
	}
}
