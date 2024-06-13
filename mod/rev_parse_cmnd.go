package mod

import (
	"fmt"
	"log"
)

func cmnd_rev(wt string, name string) {
	var header string
	header = wt

	repo := repo_find("", false)

	obj, err := object_find(repo, name, header, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(obj)
}
