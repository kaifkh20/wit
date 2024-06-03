package mod

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ObjectType int32

const (
	COMMIT ObjectType = iota
	TREE
	TAG
	BLOB
)

type GitObject struct {
	data string
}

func (gobj *GitObject) init_object(repo GitRepository, data string) {
	if data != "nil" {
		gobj.deserialize(repo)
	}
}

func (gobj *GitObject) deserialize(repo GitRepository) {
	fmt.Println("Unimplemented yet")
}

func (gobj *GitObject) serialize(repo GitRepository) {
	fmt.Println("Unimplemented")
}

func object_read(repo GitRepository, sha string) ObjectType {
	path, err := repo.repo_file(false, "objects", sha[:2], sha[2:])

	if err != nil {
		log.Fatal("An error occured", err)
	}
	if !filepath.IsLocal(path) {
		log.Fatal("No object dir")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	// fmt.Println(hex.DecodeString(string(raw)))

	// _, err = file.Read(raw)
	// if err != nil {
	// 	log.Fatal("Error reading, ", err)
	// }
	// fmt.Println(raw)

	// stri := hexdump.Dump(raw)
	// fmt.Println(stri)

	// raw = []byte(hex.EncodeToString(raw))

	buf := bytes.NewBuffer(raw)

	r, err := zlib.NewReader(buf)

	if err != nil {
		log.Fatal("Error Reading: ", err)
	}
	r.Close()

	fileContent := new(strings.Builder)

	_, err = io.Copy(fileContent, r)

	r.Close()

	// fmt.Println(hexdump.Dump([]byte(fileContent.String())))

	if err != nil {
		log.Fatal("Error copying", err)
	}

	// Sanitizing the byte

	fileContentString := fileContent.String()

	infoAboutFile := strings.Split(fileContentString, " ")

	objectType := infoAboutFile[0]
	lengthSplit := strings.Split(infoAboutFile[1], "\x00")

	fileContentString = lengthSplit[1]
	length, _ := strconv.Atoi(lengthSplit[0])

	for idx, fi := range infoAboutFile {
		if idx > 1 {
			fileContentString += fi + " "
		}
	}

	if len(fileContentString) != length {
		log.Fatal("Malformed Object")
	}

	switch objectType {
	case "commit":
		return 0
	case "tree":
		return 1
	case "tag":
		return 2
	case "blob":
		return 3
	}
	return -1
}

func Test() {
	gr := GitRepository{}
	gr.init_repo("", false)
	// fmt.Println(gr.gitdir, gr.worktree)
	objectType := object_read(gr, "e5fb83d83deb9adec6e93a4702145101740b84e7")
	fmt.Println(objectType)
}
