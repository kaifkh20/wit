package mod

import (
	"log"
	"strings"
)

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

	return kvlm_parse(raw,end+1,dct)
}


func kvlm_serialize(dct map[string][]string){
	
}