package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/christopher-dG/robotscript"
)

func main() {
	var path string

	if len(os.Args) < 2 {
		fmt.Print("Enter the path to your script file: ")
		bytes, _, err := bufio.NewReader(os.Stdin).ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		path = string(bytes)
	} else {
		path = os.Args[1]
	}

	script, err := robotscript.NewScript(path)
	if err != nil {
		log.Fatal(err)
	}

	script.Execute()
}
