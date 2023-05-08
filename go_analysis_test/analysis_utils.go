package main

import "strings"

func ajustPkgName(name string, goModuleName string) string {
	index1 := strings.Index(name, "\"")
	if index1 > 0 {
		index2 := strings.LastIndex(name, "\"")
		s := name[index1+1 : index2]
		if strings.HasPrefix(s, goModuleName) {
			return s[len(goModuleName)+1:]
		} else {
			return s
		}
	} else {
		v := strings.Split(name, " ")
		return v[1]
	}
}
