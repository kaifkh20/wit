package mod

import (
	"log"
	"strings"

	orderedmap "github.com/wk8/go-ordered-map"
)

type GitCommit struct {
	header string
	kvlm   orderedmap.OrderedMap
}

func (gcm *GitCommit) serialize() string {
	return kvlm_serialize(gcm.kvlm)
}

func (gcm *GitCommit) deserialize(data string) {
	var dct orderedmap.OrderedMap
	gcm.kvlm = kvlm_parse(data, 0, dct)
}

func kvlm_parse(raw string, start int, dct orderedmap.OrderedMap) orderedmap.OrderedMap {
	if dct.Len() == 0 {
		dct = *orderedmap.New()
	}
	spc := strings.Index(raw, " ")
	nl := strings.Index(raw, "\n")

	if spc < 0 || (nl < spc) {
		if nl == start {
			log.Fatal("nl not equals start")
		}
		dct.Set("msg", []string{raw[start-1:]})
		// dct["msg"] = []string{raw[start-1:]}
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

	if val, ok := dct.Get(key); ok {
		if val != nil {
			strArr, _ := val.([]string)
			dct.Set(key, append(strArr, value))
		}
		// dct[key] = append(val, value)
	} else {
		dct.Set(key, []string{value})
	}

	return kvlm_parse(raw, end+1, dct)
}

func kvlm_serialize(dct orderedmap.OrderedMap) string {
	result := ""
	for kv := dct.Oldest(); kv != nil; kv = kv.Next() {
		if kv.Key == "msg" {
			continue
		}
		value := kv.Value
		valueA := value.([]string)
		for _, v := range valueA {
			key := kv.Key.(string)
			result += key + " " + strings.ReplaceAll(v, "\n", "\n ") + "\n"
		}

	}
	value, _ := dct.Get("msg")
	msgV := value.([]string)
	result += "\n" + msgV[0] + "\n"
	return result
}
