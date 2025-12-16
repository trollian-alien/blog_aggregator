package main

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/lib/pq"
	"github.com/trollian-alien/blog_aggregator/internal/config"
	"github.com/trollian-alien/blog_aggregator/internal/database"
)

func main() {
	c, err := config.Read()
	if err != nil {fmt.Printf("%v", err)}
	db, err := sql.Open("postgres", c.DbURL)
	dbQueries := database.New(db)
	if err != nil {fmt.Printf("%v", err)}
	s := &state{cfg: &c, db: dbQueries}

	//setting up the commands
	mainCommands := &commands{cmds: make(map[string]func(*state, command) error)}
	mainCommands.register("login", handlerLogin)
	mainCommands.register("register", handlerRegister)
	mainCommands.register("reset", handlerReset)
	mainCommands.register("users", handlerUsers)
	mainCommands.register("agg", handlerAgg)
	mainCommands.register("addfeed", handlerAddfeed)
	mainCommands.register("feeds", handlerFeeds)
	mainCommands.register("follow", handlerFollow)
	mainCommands.register("following", handlerFollowing)

	//reading user commands
	userArgs := os.Args
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
