package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	mapCount := make(map[string]int)
	for _, v := range strings.Fields(s) {
		_, ok := mapCount[v]
		if !ok {
			mapCount[v] = 1
		} else {
			mapCount[v] += 1
		}
	}
	return mapCount
}

func main() {
	wc.Test(WordCount)
}

