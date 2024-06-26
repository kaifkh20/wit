package mod

import (
	"bytes"
	"compress/zlib"
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map"
)

type ObjectTypes struct {
	header string
	GitCommit
	GitBlob
	GitTree
	GitTag
}

type GitObject struct {
	header string
	data   string
	ObjectTypes
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

func object_read(repo GitRepository, sha string) ObjectTypes {
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
		return ObjectTypes{header: objectType, GitCommit: GitCommit{header: "commit", kvlm: *orderedmap.New()}}
	case "tree":
		return ObjectTypes{header: objectType, GitTree: GitTree{header: "tree", items: []GitTreeLeaf{}}}
	case "tag":
		return ObjectTypes{header: objectType, GitTag: GitTag{GitCommit: GitCommit{header: "commit", kvlm: *orderedmap.New()}}}
	case "blob":
		return ObjectTypes{header: objectType, GitBlob: GitBlob{header: "blob", blobData: fileContentString}}
	}
	return ObjectTypes{header: "error"}
}

func (gobj *GitObject) obj_write(repo *GitRepository) string {

	gobj.serialize(*repo)
	result := gobj.header + " " + string(len(gobj.data)) + "\x00" + gobj.data
	h := crypto.SHA1.New()
	io.WriteString(h, result)
	hBy := h.Sum(nil)
	encS := hex.EncodeToString(hBy)

	path, err := repo.repo_file(true, "objects", encS[:2], encS[2:])
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal("Error writing object", err)
		}
		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		w.Write([]byte(result))
		w.Close()
		file.Write(b.Bytes())
		file.Close()
	}

	return encS
}

func obj_resolve(repo GitRepository, name string) ([]string, error) {
	candidates := []string{}
	hashRE, err := regexp.Compile("^[0-9A-Fa-f]{4,40}$")
	if err != nil {
		log.Fatal(err)
	}
	if strings.TrimSpace(name) == "" {
		return []string{}, fmt.Errorf("Empty name")
	}
	if name == "HEAD" {
		ref, err := ref_resolve(repo, "HEAD")
		if err != nil {
			log.Fatal("Error occured in resolving HEAD")
		}
		return []string{ref}, nil
	}
	if hashRE.Match([]byte(name)) {
		name = strings.ToLower(name)
		prefix := name[0:2]
		path := repo.repo_dir(false, []string{"objects"})
		rem := name[2:]
		enteries, err := os.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range enteries {
			if strings.HasPrefix(f.Name(), rem) {
				candidates = append(candidates, prefix+f.Name())
			}
		}
	}
	as_tag, err := ref_resolve(repo, "refs/tags/"+name)
	if err != nil {
		log.Fatal("Error here object resolve")
	}
	candidates = append(candidates, as_tag)
	as_branch, err := ref_resolve(repo, "refs/heads/"+name)
	if err != nil {
		log.Fatal("Error here in object resolve heads")
	}
	candidates = append(candidates, as_branch)
	return candidates, nil

}

func object_find(repo GitRepository, name string, header string, follow bool) (string, error) {
	sha, err := obj_resolve(repo, name)
	if err != nil {
		log.Fatal("No such reference", name, err)
	}

	if len(sha) > 1 {
		log.Fatalf("Ambigious reference %s: Candidates are:\n - %s", name, strings.Join(sha, "\n -"))
	}

	sha_s := sha[0]

	if header == "" {
		return sha_s, nil
	}

	for {
		obj := object_read(repo, sha_s)
		if obj.header == header {
			return sha_s, nil
		}
		if !follow {
			return "", nil
		}
		if obj.header == "tag" {
			val, _ := obj.kvlm.Get("object")
			sha_s = val.(string)
		} else if obj.header == "commit" && header == "tree" {
			val, _ := obj.kvlm.Get("tree")
			sha_s = val.(string)
		} else {
			return "", nil
		}
	}

	// return name
}

func Test() {
	gr := GitRepository{}
	gr.init_repo("", false)

	// fmt.Println(gr.gitdir, gr.worktree)
	// objectType := object_read(gr, "e5fb83d83deb9adec6e93a4702145101740b84e7")
	// fmt.Println(objectType)
}
