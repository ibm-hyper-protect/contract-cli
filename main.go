package main

import (
	"log"
	"os"

	"github.com/ibm-hyper-protect/contract-cli-go/cli"
)

func main() {
	err := cli.CreateApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
