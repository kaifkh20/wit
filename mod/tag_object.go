package mod

import (
	"fmt"
	"log"
	"os"

	orderedmap "github.com/wk8/go-ordered-map"
)

type GitTag struct {
	GitCommit
}

func cmnd_tag(name string, object string, create_tag_object bool) {
	repo := repo_find("", true)

	var typeO string
	if create_tag_object {
		typeO = "object"
	} else {
		typeO = "ref"
	}

	fmt.Println(typeO)

	if name != "" {
		tag_create(repo, name, object, create_tag_object)
	} else {
		refs := ref_list(repo, "")
		ta, _ := refs.Get("tags")
		tags := ta.(orderedmap.OrderedMap)
		// show_ref(repo, refs.Get("tags").(orderedmap.OrderedMap), false, "")
		show_ref(repo, tags, false, "")
	}
}

func tag_create(repo GitRepository, name string, object string, create_tag_object bool) {
	sha := object_find(repo, "", object, false)

	if create_tag_object {
		tag := GitTag{GitCommit: GitCommit{header: "tag"}}
		tag.kvlm = *orderedmap.New()
		tag.kvlm.Set("object", sha)
		tag.kvlm.Set("type", "commit")
		tag.kvlm.Set("tag", name)
		tag.kvlm.Set("tagger", "Wit <wit@wit.com>")
		tag.kvlm.Set("msg", "Just a message")
		obj := GitObject{header: "tag", ObjectTypes: ObjectTypes{GitTag: tag}}
		tag_sha := obj.obj_write(&repo)
		ref_create(repo, "tags/"+name, tag_sha)

	} else {
		ref_create(repo, "tags/"+name, sha)
	}
}

func ref_create(repo GitRepository, ref_name string, sha string) {
	path_s, _ := repo.repo_file(false, "refs/"+ref_name)
	f, err := os.Open(path_s)
	if err != nil {
		log.Fatal(err)
	}
	f.Write([]byte(sha + "\n"))
	defer f.Close()
}
