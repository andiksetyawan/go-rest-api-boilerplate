package main

import "go-rest-api-boilerplate/cmd/commands"

func main() {
	commands.NewRootCommand().Execute()
}
