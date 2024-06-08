package main

import (
	"fmt"
	"log"
	"os"
	"wit/mod"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("wit", "Wit executes command")
	// cmnd := parser.String("c", "command", &argparse.Options{Required: true, Help: "Enter a Command"})
	// arguments := parser.String("a", "arguments", &argparse.Options{Required: true, Help: "Enter an argument if not then enter an empty string"})

	init := parser.NewCommand("init", "Initialize A Repository")
	path := init.String("p", "path", &argparse.Options{Required: true, Help: "Path where repo should be initialized"})

	test := parser.NewCommand("test", "Test this app")

	cat_file := parser.NewCommand("cat-file", "Provide content of repository objects")
	type_cmnd := cat_file.Selector("t", "type", []string{"blob", "commit", "tag", "tree"}, &argparse.Options{Required: true, Help: "Specify type"})
	object := cat_file.String("o", "object", &argparse.Options{Required: true, Help: "The object to display"})

	hash_obj := parser.NewCommand("hash-object", "Compute Object ID and optionally creates a blob from a file")
	// write_cmnd := hash_obj.String("w", "write", &argparse.Options{Help: "Actual write"})
	type_cmnd_hash := hash_obj.Selector("t", "type", []string{"blob", "commit", "tag", "tree"}, &argparse.Options{Required: true, Help: "Specify type"})
	hash_path := hash_obj.String("p", "path", &argparse.Options{Required: true, Help: "Actual Path"})

	logCmnd := parser.NewCommand("log", "Display log of commit")
	commitLog := logCmnd.String("c", "commit", &argparse.Options{Default: "HEAD", Help: "Commit to start at"})

	err := parser.Parse(os.Args)

	if err != nil {
		log.Fatal(parser.Usage(err))
	}

	if init.Happened() {
		mod.Cmnd_Init(*path)
	} else if cat_file.Happened() {
		fmt.Print(cat_file.GetName(), *type_cmnd, *object)
	} else if hash_obj.Happened() {
		mod.Cmnd_Hash(true, *hash_path, *type_cmnd_hash)
	} else if logCmnd.Happened() {
		fmt.Println(*commitLog)
	} else if test.Happened() {
		mod.Test()
	}

	// fmt.Println(, path)
	// fmt.Println(*cmnd, *arguments)
	// switch *cmnd {
	// case "init":
	// 	mod.Cmnd_Init(*arguments)
	// case "test":
	// 	mod.Test()
	// }

}
