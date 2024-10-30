package main

import (
	"log"

	"github.com/everettraven/crd-diff/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
