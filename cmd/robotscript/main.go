package main

import (
	"log"
	"os"

	"github.com/christopher-dG/robotscript"
)

func main() {
	script, err := robotscript.NewScript(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	script.Execute()
}
