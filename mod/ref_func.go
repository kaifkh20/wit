package mod

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map"
)

func ref_resolve(repo GitRepository, ref string) (string, error) {
	path, err := repo.repo_file(false, ref)
	if err != nil {
		log.Fatal(err)
	}

	if f, err := os.Stat(path); f.IsDir() {
		if err != nil {
			// log.Fatal(err)
			return "", err
		}
		return "", err
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	var data []byte
	_, err = f.Read(data)
	if err != nil {
		log.Fatal("Error reading file", err)
	}
	data = data[:len(data)-1]

	if strings.HasPrefix(string(data), "ref: ") {
		return ref_resolve(repo, string(data[5:]))
	} else {
		return string(data), nil
	}
}

func ref_list(repo GitRepository, path string) orderedmap.OrderedMap {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = repo.repo_dir(true, []string{"refs"})
	}
	ret := orderedmap.New()

	ldir, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range ldir {
		can := filepath.Join(path, f.Name())
		if fi, _ := os.Stat(can); fi.IsDir() {
			ret.Set(f.Name(), ref_list(repo, can))
		} else {
			p, err := ref_resolve(repo, can)
			if err != nil {
				log.Fatal("Error resolving ref")
			}
			ret.Set(f.Name(), p)
		}
	}
	return *ret
}

func show_ref(repo GitRepository, refs orderedmap.OrderedMap, with_hash bool, prefix string) {
	for pair := refs.Oldest(); pair != nil; pair = pair.Next() {
		if reflect.TypeOf(pair.Value) == reflect.TypeOf("") {
			val := pair.Value.(string)
			result := ""
			if with_hash {
				result += val + " "
			}
			if prefix != "" {
				result += prefix + "/"
			}
			result += pair.Key.(string)
			fmt.Println(result)
		} else {
			if prefix != "" {
				prefix += "/"
			}
			prefix += pair.Key.(string)
			show_ref(repo, pair.Value.(orderedmap.OrderedMap), with_hash, prefix)
		}
	}
}

func Ref_Command() {
	repo := repo_find("", false)
	refs := ref_list(repo, "")
	show_ref(repo, refs, true, "refs")
}
