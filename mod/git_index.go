package mod

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
)

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
	ctime             []int
	mtime             []int
	dev               int
	ino               uint32
	mode_type         uint32
	mode_perms        uint32
	uid               uint32
	gid               uint32
	fsize             uint32
	sha               string
	flag_assume_valid bool
	flag_stage        uint32
	name              string
}

type GitIndex struct {
	version  int
	enteries []GitIndexEntry
}

func (gi *GitIndex) init(version int, enteries []GitIndexEntry) {
	gi.version = version
	gi.enteries = enteries
}

func index_read(repo GitRepository) GitIndex {
	index_file, err := repo.repo_file(false, "index")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(index_file); os.IsNotExist(err) {
		return GitIndex{}
	}

	fi, err := os.Open(index_file)

	if err != nil {
		log.Fatal("Error opening file", err)
	}
	var raw []byte
	_, err = fi.Read(raw)

	if err != nil {
		log.Fatal("Error reading file", err)
	}

	header := raw[:12]
	signature := header[:4]

	if string(signature) != "DIRC" {
		log.Fatal("Signature is not DIRC")
	}
	ve := header[4:8]
	version := binary.BigEndian.Uint32(ve)

	if version != 2 {
		log.Fatal("Wit only supports index file version 2")
	}
	count := binary.BigEndian.Uint32(header[8:12])

	enteries := []GitIndexEntry{}

	content := raw[12:]

	idx := 0

	for i := 0; i < int(count); i += 1 {
		ctime_s := binary.BigEndian.Uint32(content[idx : idx+4])
		ctime_ns := binary.BigEndian.Uint32(content[idx+4 : idx+8])
		mtime_s := binary.BigEndian.Uint32(content[idx+8 : idx+12])
		mtime_ns := binary.BigEndian.Uint32(content[idx+12 : idx+16])
		dev := binary.BigEndian.Uint32(content[idx+16 : idx+20])
		ino := binary.BigEndian.Uint32(content[idx+20 : idx+24])
		unused_ig := binary.BigEndian.Uint32(content[idx+24 : idx+26])
		if unused_ig == 0 {
			log.Fatal("Error in reading index ")
		}
		mode := binary.BigEndian.Uint32(content[idx+26 : idx+28])
		mode_type := mode >> 12
		if !slices.Contains([]uint32{0b1000, 0b1010, 0b1110}, mode_type) {
			log.Fatal("Invalid mode_type")
		}
		mode_perms := mode & 0b0000000111111111
		uid := binary.BigEndian.Uint32(content[idx+28 : idx+32])
		gid := binary.BigEndian.Uint32(content[idx+32 : idx+36])
		fsize := binary.BigEndian.Uint32(content[idx+36 : idx+40])

		sha := fmt.Sprintf("%040x", binary.BigEndian.Uint32(content[idx+40:idx+60]))
		flags := binary.BigEndian.Uint32(content[idx+60 : idx+62])
		flag_assume_valid := ((flags & 0b1000000000000000) != 0)
		flag_extended := ((flags & 0b0100000000000000) != 0)
		if flag_extended {
			log.Fatal("Problem in flag_extended")
		}
		flag_stage := flags & 0b0011000000000000
		name_length := flags & 0b0000111111111111
		idx += 62

		var raw_name []byte
		if name_length < 0xFFF {
			if content[idx+int(name_length)] != 0x00 {
				log.Fatal("Content 0x00")
			}
			raw_name = content[idx : idx+int(name_length)]
			idx += int(name_length) + 1
		} else {
			fmt.Printf("Notice: Name is 0x%d bytes long", name_length)
			// null_idx := byte(content[idx+0xFFF:], "\x00")
			null_idx := bytes.IndexByte(content[idx+0xFFF:], 0x00)
			raw_name = content[idx:null_idx]
			idx = null_idx + 1
		}
		name := string(raw_name)
		idx = int(8 * math.Ceil(float64(idx)/8.0))
		enteries = append(enteries, GitIndexEntry{
			ctime:             []int{int(ctime_s), int(ctime_ns)},
			mtime:             []int{int(mtime_s), int(mtime_ns)},
			dev:               int(dev),
			ino:               ino,
			mode_type:         mode_type,
			mode_perms:        mode_perms,
			uid:               uid,
			gid:               gid,
			fsize:             fsize,
			sha:               sha,
			flag_assume_valid: flag_assume_valid,
			flag_stage:        flag_stage,
			name:              name,
		})
	}
	// raw:=fi.Read()
	return GitIndex{version: int(version), enteries: enteries}
}
