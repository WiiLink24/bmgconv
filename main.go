package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("./message.bmg")
	if err != nil {
		panic(err)
	}

	bmg, err := parseBMG(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bmg))
}
