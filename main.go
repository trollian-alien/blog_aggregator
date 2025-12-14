package main

import (
	"fmt"
	"os"
	"github.com/trollian-alien/blog_aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {fmt.Printf("%v", err)}
	s := &state{cfg: &c}

	//setting up the commands
	mainCommands := &commands{cmds: make(map[string]func(*state, command) error)}
	mainCommands.register("login", handlerLogin)
	userArgs := os.Args

	//reading user commands
	var cmd command
	if len(userArgs) < 2 {
		fmt.Println("No commands given")
		os.Exit(1)
	} else if len(userArgs) > 2 {
		cmd = command{name: userArgs[1], args: userArgs[2:]}
	} else {
		cmd = command{name: userArgs[1], args: []string{}}
	}
	err = mainCommands.run(s, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
