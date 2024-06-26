package mod

import (
	"fmt"
	"log"
)

func cmnd_cat_file(object string, typeObj string, args ...string) {
	repo := repo_find(".", true)
	cat_file(repo, object, typeObj)
}

func cat_file(repo GitRepository, typeObj string, header string) {
	obj_f, err := object_find(repo, "", typeObj, true)
	if err != nil {
		log.Fatal(err)
	}
	obj := object_read(repo, obj_f)
	if obj.header == "error" {
		log.Fatal("Unknow Type Object")
	}

	var response string

	switch obj.header {
	case "commit":
		response = obj.GitCommit.serialize()
	case "blob":
		response = obj.GitBlob.serialize()
	}

	fmt.Println(response)

}
