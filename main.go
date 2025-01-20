package main

import (
	"github.com/compliance-framework/framework/cmd"
	"log"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		log.Fatal(err)
	}
}
