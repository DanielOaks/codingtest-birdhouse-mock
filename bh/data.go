// Copyright (c) 2023 Daniel Oaks <daniel@danieloaks.net>
// released under the CC0 license

package bh

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed moby-word-lists/names-f.txt
var namesF string

//go:embed moby-word-lists/names-m.txt
var namesM string

func getNameList() (names []string) {
	names = strings.Split(namesF, "\n")
	names = append(names, strings.Split(namesM, "\n")...)

	return names
}

type Data struct {
}

func GenerateData(config BirdhousesConfig) (data Data) {
	names := getNameList()
	fmt.Println(names)

	return data
}
