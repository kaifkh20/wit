package mod

import (
	"fmt"
	"log"
	"path/filepath"
)

func cmnd_ls_tree(tree string, recursive bool) {
	repo := repo_find("", false)
	ls_tree(repo, tree, recursive, "")
}

func ls_tree(repo GitRepository, tree string, recursive bool, prefix string) {
	sha, err := object_find(repo, tree, "tree", true)
	if err != nil {
		log.Fatal(err)
	}
	obj := object_read(repo, sha)
	for _, items := range obj.items {
		var type_ string
		if len(items.mode) == 5 {
			type_ = items.mode[0:1]
		} else {
			type_ = items.mode[0:2]
		}

		switch type_ {
		case "04":
			type_ = "tree"
		case "10":
			type_ = "blob"
		case "12":
			type_ = "blob"
		case "16":
			type_ = "commit"
		default:
			log.Fatal("Weird tree leaf mode", type_)
		}

		if !recursive && type_ == "tree" {
			fmt.Printf("%s %s %s \t%s\n",
				items.mode,
				type_,
				items.sha,
				filepath.Join(prefix, items.path),
			)
		} else {
			ls_tree(repo, items.sha, recursive, filepath.Join(prefix, items.path))
		}
	}
}
