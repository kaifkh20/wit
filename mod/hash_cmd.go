package mod

import (
	"fmt"
	"log"
	"os"
)

func Cmnd_Hash(write bool, path string, typeObj string) {
	var repo GitRepository
	if write {
		// repo := repo_find()
	}
	file, err := os.Open(path)

	if err != nil {
		log.Fatal("Error Opening a file")
	}

	hash_str := object_hash(file, typeObj, repo)
	fmt.Println(hash_str)
}

func object_hash(file *os.File, typeObj string, repo GitRepository) string {
	var raw []byte

	var obj GitObject
	_, err := file.Read(raw)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}

	obj.data = string(raw)
	obj.header = typeObj

	return obj.obj_write(&repo)
}
