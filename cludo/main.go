package main

import (
	"flag"
	"log"
	"os"

	"github.com/oxtopus/cludo"
)

func main() {
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cliArgs := os.Args[1:]

	cludo.Run(cwd, cliArgs)
}
