package mod

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"strconv"
	"strings"
)

type GitTreeLeaf struct {
	mode string
	path string
	sha  string
}

type GitTree struct {
	header string
	items  []GitTreeLeaf
}

func (gt *GitTree) serialize() string {
	return tree_serialize(*gt)
}

func (gt *GitTree) deserialize(data string) {
	gt.items = tree_parse(data)
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

func tree_leaf_sort_key(leaf GitTreeLeaf) string {
	if strings.HasPrefix(leaf.mode, "10") {
		return leaf.path
	} else {
		return leaf.path + "/"
	}
}

func tree_serialize(tree GitTree) string {
	for _, item := range tree.items {
		item.path = tree_leaf_sort_key(item)
	}
	result := ""
	for _, item := range tree.items {
		result += item.mode
		result += " "
		result += item.path
		result += "\x00"
		sha, err := strconv.Atoi(item.sha)
		if err != nil {
			log.Fatal(err)
		}
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, sha)
		result += buf.String()

	}
	return result
}
