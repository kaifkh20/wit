package mod

import (
	"log"
	"strings"
)

type GitCommit struct {
	header string
	kvlm   map[string][]string
}

func (gcm *GitCommit) serialize() string {
	return kvlm_serialize(gcm.kvlm)
}

func (gcm *GitCommit) deserialize(data string) {
	var dct map[string][]string
	gcm.kvlm = kvlm_parse(data, 0, dct)
}

func kvlm_parse(raw string, start int, dct map[string][]string) map[string][]string {
	if len(dct) == 0 {
		dct = make(map[string][]string)
	}
	spc := strings.Index(raw, " ")
	nl := strings.Index(raw, "\n")

	if spc < 0 || (nl < spc) {
		if nl == start {
			log.Fatal("nl not equals start")
		}
		dct["msg"] = []string{raw[start-1:]}
		return dct
	}

	key := raw[start:spc]

	end := start
	for {
		temp := raw[end+1:]
		end = strings.Index(temp, "\n")
		if string(temp[end+1]) != " " {
			break
		}
	}

	value := strings.ReplaceAll(raw[spc+1:end], "\n ", "\n")

	if val, ok := dct[key]; ok {
		dct[key] = append(val, value)
	} else {
		dct[key] = []string{value}
	}

	return kvlm_parse(raw, end+1, dct)
}

func kvlm_serialize(dct map[string][]string) string {
	result := ""
	for k := range dct {
		if k == "msg" {
			continue
		}
		value := dct[k]
		for _, v := range value {
			result += k + " " + strings.ReplaceAll(v, "\n", "\n ") + "\n"
		}
		result += "\n" + dct["msg"][0] + "\n"
	}
	return result
}
