package mod

import (
	"log"
	"os"
	"path/filepath"
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
		path = repo.repo_dir(true, []string{path})
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
