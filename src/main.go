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

	ls_tree := parser.NewCommand("ls-tree", "Pretty print a tree object")
	recursive := ls_tree.Flag("r", "recursive", &argparse.Options{Required: false, Help: "Recurse into sub-trees", Default: nil})
	tree := ls_tree.String("t", "tree", &argparse.Options{Required: true, Help: "A tree object"})

	chckout := parser.NewCommand("checkout", "Checkout a commit inside of a directory")
	commit := chckout.String("c", "commit", &argparse.Options{Help: "Commit to checkout to"})
	path_c := chckout.String("p", "path", &argparse.Options{Help: "The Empty Directory to checkout on"})

	show_ref := parser.NewCommand("show-ref", "List References")

	tag := parser.NewCommand("tag", "List and create tags")
	a := tag.Flag("a", "-a", &argparse.Options{Help: "Whether to create a tag object"})
	name := tag.String("n", "name", &argparse.Options{Help: "The new tag's name"})
	object_tag := tag.String("o", "object", &argparse.Options{Default: "HEAD", Help: "The object tag will point to"})

	rev_parse := parser.NewCommand("rev-parse", "Parse revision identifiers")
	wt := rev_parse.Selector("wt", "wyag-type", []string{"blob", "commit", "tag", "tree"}, &argparse.Options{Help: "Specify the expected type"})
	name = rev_parse.String("n", "name", &argparse.Options{Help: "The name to parse"})
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
	} else if ls_tree.Happened() {
		fmt.Println(*recursive, *tree)
	} else if chckout.Happened() {
		fmt.Print(*commit, *path_c)
	} else if show_ref.Happened() {
		mod.Ref_Command()
	} else if tag.Happened() {
		fmt.Println(*a, *name, *object_tag)
	} else if rev_parse.Happened() {
		fmt.Println(*wt, *name)
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
