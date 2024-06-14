package mod

/*
	ctime : last time file's meta data changed
	mtime : last time file's data changed
	dev : ID of device containing this file
	ino : file's inode number
	mode_type : Object type {regular,symlink,gitlink}
	mode_perms : Permission
	uid : User Id of owner
	gid : Group id
	fsize : Object size
	sha : Object's sha
	name : Name of object
*/

type GitIndexEntry struct {
	ctime             string
	mtime             string
	dev               string
	ino               string
	mode_type         string
	mode_perms        int
	uid               string
	gid               string
	fsize             string
	sha               string
	flag_assume_valid string
	flag_stage        string
	name              string
}

type GitIndex struct {
	version  int
	enteries []string
}

func (gi *GitIndex) init(version int, enteries []string) {
	gi.version = version
	gi.enteries = enteries
}
