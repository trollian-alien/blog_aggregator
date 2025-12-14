package main

import (
	"fmt"

	"github.com/trollian-alien/blog_aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {fmt.Printf("%v", err)}
	user := "me"
	err = c.SetUser(user)
	if err != nil {fmt.Printf("%v", err)}
	c, err = config.Read()
	if err != nil {fmt.Printf("%v", err)}
	fmt.Println(c)
}
