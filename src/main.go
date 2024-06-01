package main

import (
	"fmt"
	"os"

	"wit/mod"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("wit", "Wit executes command")
	cmnd := parser.String("c", "command", &argparse.Options{Required: true, Help: "Enter a Command"})
	arguments := parser.String("a", "arguments", &argparse.Options{Required: true, Help: "Enter an argument if not then enter an empty string"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	fmt.Println(*cmnd, *arguments)
	switch *cmnd {
	case "init":
		mod.Cmnd_Init(*arguments)
	}
}
