package mod

import (
	"fmt"
	"log"
	"strings"
)

func cmnd_log(commit string, args ...string) {
	repo := repo_find(".", true)
	fmt.Print("digraph witlog{")
	fmt.Print("\t node [shape=rect]")
	obj_f, err := object_find(repo, "", commit, false)
	if err != nil {
		log.Fatal(err)
	}
	log_graphviz(repo, obj_f, map[string]bool{})
}

func log_graphviz(repo GitRepository, sha string, seen map[string]bool) {
	if _, ok := seen[sha]; !ok {
		return
	}
	seen[sha] = true

	objt := object_read(repo, sha)
	if objt.header == "error" {
		log.Fatal("Error object type")
	}

	commit := objt.GitCommit
	// short_hash := sha[:8]
	mv, _ := commit.kvlm.Get("mess")
	mess := mv.([]string)

	if strings.Contains(mess[0], "\n") {
		mess = mess[:strings.Index(mess[0], "\n")]
	}
	fmt.Printf(" c_%s [label= \\ %s: %s \\]", sha, sha[0:7], mess)

	if _, ok := commit.kvlm.Get("parent"); !ok {
		return
	}

	pv, _ := commit.kvlm.Get("parent")
	parents := pv.([]string)

	for _, p := range parents {
		fmt.Printf("c_%s -> c_%s;", sha, p)
		log_graphviz(repo, p, seen)
	}

}
