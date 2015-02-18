package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

func main() {
	flag.Parse()
	content, _ := ioutil.ReadFile(flag.Args()[0])
	fmt.Print(content)
}
