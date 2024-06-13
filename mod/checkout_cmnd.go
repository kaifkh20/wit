package mod

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func chck_cmnd(commit string, path string, args ...string) {
	repo := repo_find("", false)
	obj_f, err := object_find(repo, commit, "", true)
	if err != nil {
		log.Fatal(err)
	}
	obj := object_read(repo, obj_f)
	si, _ := obj.kvlm.Get("tree")
	tree_val := si.(string)
	if obj.header == "commit" {
		obj = object_read(repo, tree_val)
	}
	if f, err := os.Stat(path); os.IsExist(err) {

		if !f.IsDir() {
			log.Fatal("Not a directory", path)

		}
	} else {
		fmt.Println("No such path exists", path)
		os.Mkdir(path, 0755)
	}
	path, _ = filepath.Abs(path)
	tree_checkout(repo, obj.GitTree, path)

}

func tree_checkout(repo GitRepository, tree GitTree, path string) {
	for _, item := range tree.items {
		obj := object_read(repo, item.sha)
		dest := filepath.Join(path, item.path)

		if obj.header == "tree" {
			os.Mkdir(dest, 07)
			tree_checkout(repo, obj.GitTree, dest)
		} else if obj.header == "blob" {
			fi, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			fi.Write([]byte(obj.GitBlob.blobData))
		}

	}
}
