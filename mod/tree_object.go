package mod

import (
	"encoding/hex"
	"log"
	"strings"
)

type GitTreeLeaf struct {
	mode string
	path string
	sha  string
}

func parse_tree_one(raw string, start int) (int, GitTreeLeaf) {
	x := strings.Index(raw, " ")

	if x == -1 {
		log.Fatal("Error in tree")
	}

	if x-start != 5 || x-start != 6 {
		log.Fatal("Error in reading mode")
	}

	mode := raw[start:x]
	if len(mode) == 5 {
		mode = " " + mode
	}
	y := strings.Index(raw, "\x00")
	path := raw[x+1 : y]
	sha := hex.EncodeToString([]byte(raw[y+1 : y+21]))

	return y + 21, GitTreeLeaf{mode, path, sha}
}

func tree_parse(raw string) []GitTreeLeaf {
	posT := 0
	max := len(raw)
	ret := []GitTreeLeaf{}

	for posT < max {
		pos, data := parse_tree_one(raw, posT)
		posT = pos
		ret = append(ret, data)
	}
	return ret
}
