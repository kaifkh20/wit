package mod

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type GitRepository struct {
	worktree string
	gitdir   string
	conf     *ini.File
}

func init_repo(gr GitRepository, path string, force bool) GitRepository {

	if path == "" {
		gr.worktree, _ = os.Getwd()
	} else {
		gr.worktree = path
	}

	// fmt.Println(path)
	gr.gitdir = filepath.Join(path, ".git")

	// fmt.Println(gr.gitdir, filepath.IsLocal(gr.gitdir))

	if !force || !filepath.IsLocal(gr.gitdir) {
		log.Fatal("Not Git Repository")
	}

	mkdir := false
	if force {
		mkdir = true
	}

	cf, err := repo_file(gr, mkdir, "config")

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(cf)
	if err != nil {
		log.Fatal("Error opening config file", err)
	}

	fmt.Println("This is cf", cf)
	// cf:=re
	if cf != "" && err == nil {
		if _, err := os.Stat(cf); err == nil {
			conf, err := ini.Load(cf)
			if err != nil {
				log.Fatal("This error is in Config", err)
			}
			gr.conf = conf
			fmt.Println("Parsed config parser")
		}
	} else {
		if !force {
			log.Fatal("Config file missings")
		}
	}
	if !force {
		// if gr.conf.HasSection("core")
		vers := gr.conf.Section("core").Key("repositoryformatversion").String()
		if err != nil {
			log.Fatal(err)
		}
		if ver, _ := strconv.Atoi(vers); ver != 0 {
			log.Fatal("Unsupported Repo Version")
		}
	}
	file.Close()
	return gr
}

func repo_path(repo GitRepository, path []string) (string, error) {
	var pathString strings.Builder
	for _, p := range path {
		pathString.Write([]byte(p))
	}
	return filepath.Join(repo.gitdir, strings.TrimSpace(pathString.String())), nil
}

func repo_file(repo GitRepository, mkdir bool, path ...string) (string, error) {

	if p := repo_dir(repo, mkdir, path[:len(path)-1]); p != "" {
		os.Create(filepath.Join(p, path[len(path)-1]))
		return repo_path(repo, path)
	}
	return "", nil
}

func repo_dir(repo GitRepository, mkdir bool, path []string) string {

	pathString, _ := repo_path(repo, path)
	fmt.Println(pathString)
	if _, err := os.Stat(pathString); err != nil {
		// mkdir = true
		fmt.Println("No Git Directory", err)
		// return pathString
	} else {
		// fmt.Println("Returning path string", pathString)
		return pathString
	}
	if mkdir {
		err := os.MkdirAll(pathString, 0755)
		if err != nil {
			log.Fatal("Error in creating directory", err)
		}
		fmt.Println(".git created")
		return pathString
	}
	return ""
}

func repo_create(path string) GitRepository {
	// fmt.Println(path)
	repo := GitRepository{}
	repo = init_repo(repo, path, true)
	// fmt.Println(repo, "Repo Initalized")
	// fmt.Println()
	if fi, err := os.Stat(repo.worktree); err != nil || !fi.IsDir() {
		log.Fatalf("%s is not a directory", path)
	}
	// fi, err := os.Open(repo.gitdir)
	// if err == nil {
	// 	// fmt.Println(fi.Name())
	// 	names, err := fi.Readdirnames(2)
	// 	if err != nil || len(names) >= 0 {
	// 		log.Fatalf("%s not an empty directory", path)
	// 	}
	// } else {
	// 	os.Mkdir(repo.worktree, 1)
	// }

	repo_dir(repo, true, []string{"branches"})
	repo_dir(repo, true, []string{"objects"})
	repo_dir(repo, true, []string{"refs", "tags"})
	repo_dir(repo, true, []string{"refs", "heads"})

	fmt.Println("Created all refs and HEADS")

	pS, err := repo_file(repo, false, "description")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pS)
	file, err := os.OpenFile(pS, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write([]byte("Unnamed repository; edit this file 'description' to name the repository.\n"))
	if err != nil {
		log.Fatal("Error on writing description file", err)
	}
	file.Close()

	pS, err = repo_file(repo, false, "HEAD")
	if err != nil {
		log.Fatal(err)
	}
	file, err = os.OpenFile(pS, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write([]byte("ref: refs/heads/master\n"))
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	pS, err = repo_file(repo, false, "config")
	if err != nil {
		log.Fatal(err)
	}
	file, err = os.OpenFile(pS, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	// config, err := ini.Load(file)

	if err != nil {
		log.Fatal(err)
	}
	// co
	config := repo_default_config()
	config.WriteTo(file)

	file.Close()
	return repo
}

func repo_default_config() *ini.File {
	ret := ini.Empty()

	core, err := ret.NewSection("core")
	if err != nil {
		log.Fatal("Error in sec ore", err)
	}
	core.NewKey("repositoryformatversion", "0")
	core.NewKey("filemode", "false")
	core.NewKey("bare", "false")
	// fmt.Println(ret)

	return ret
}

func Cmnd_Init(path string) {
	fmt.Println("Path given", path)
	repo_create(path)
	fmt.Println("Initalized Succesfully")
}
