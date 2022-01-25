package main

import (
	"fmt"
	"frugal-hero/outputs"
	"frugal-hero/services"
	"log"
	"os"
	"time"
)

func track(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s, execution time %s\n", name, time.Since(start))
	}
}

func main() {
	defer track("main")()
	service, err := services.GetService(os.Args[1])

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	outputName := ""
	if len(os.Args) > 2 {
		outputName = os.Args[2]
	}
	service.Inspect(outputs.GetOutput(outputName))
}
