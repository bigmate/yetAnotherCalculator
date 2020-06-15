package main

import (
	"fmt"
	"os"
	"os/user"
	"yac/repl"
)

func main() {
	var u, err = user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, welcome to yac!\n", u.Username)
	repl.Start(os.Stdin, os.Stdout)
}
