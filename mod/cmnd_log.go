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
	log_graphviz(repo, object_find(repo, "", commit, false), map[string]bool{})
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
	mess := commit.kvlm["msg"][0]

	if strings.Contains(mess, "\n") {
		mess = mess[:strings.Index(mess, "\n")]
	}
	fmt.Printf(" c_%s [label= \\ %s: %s \\]", sha, sha[0:7], mess)

	if _, ok := commit.kvlm["parent"]; !ok {
		return
	}

	parents := commit.kvlm["parent"]

	for _, p := range parents {
		fmt.Printf("c_%s -> c_%s;", sha, p)
		log_graphviz(repo, p, seen)
	}

}
