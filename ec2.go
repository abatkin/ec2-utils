package main

import (
	"ec2-utils/commands"
	"log"
)

func main() {

	rootCommand := commands.BuildRootCommand()

	if err := rootCommand.Execute(); err != nil {
		log.Fatalf("error running ec2: %v", err)
	}
}

