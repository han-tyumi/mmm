package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/han-tyumi/mcf"
)

var nameToSortType = map[string]mcf.SortType{
	"0":        mcf.Featured,
	"f":        mcf.Featured,
	"feat":     mcf.Featured,
	"featured": mcf.Featured,

	"1":          mcf.Popularity,
	"p":          mcf.Popularity,
	"pop":        mcf.Popularity,
	"popularity": mcf.Popularity,

	"2":          mcf.LastUpdate,
	"l":          mcf.LastUpdate,
	"last":       mcf.LastUpdate,
	"u":          mcf.LastUpdate,
	"up":         mcf.LastUpdate,
	"update":     mcf.LastUpdate,
	"lastupdate": mcf.LastUpdate,

	"3":    mcf.Name,
	"n":    mcf.Name,
	"name": mcf.Name,

	"4":      mcf.Author,
	"a":      mcf.Author,
	"auth":   mcf.Author,
	"author": mcf.Author,

	"5":              mcf.TotalDownloads,
	"t":              mcf.TotalDownloads,
	"total":          mcf.TotalDownloads,
	"d":              mcf.TotalDownloads,
	"down":           mcf.TotalDownloads,
	"downloads":      mcf.TotalDownloads,
	"totaldownloads": mcf.TotalDownloads,
}

var sortTypeToName = map[mcf.SortType]string{
	mcf.Featured:       "featured",
	mcf.Popularity:     "popularity",
	mcf.LastUpdate:     "lastupdate",
	mcf.Name:           "name",
	mcf.Author:         "author",
	mcf.TotalDownloads: "totaldownloads",
}

type sortType mcf.SortType

func (t *sortType) Set(s string) error {
	p := regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(s, "")
	p = strings.ToLower(p)

	if sort, ok := nameToSortType[p]; ok {
		*t = sortType(sort)
		return nil
	}

	return fmt.Errorf("%s is not a valid sort type", s)
}

func (t *sortType) Type() string {
	return "sortType"
}

func (t *sortType) String() string {
	return sortTypeToName[mcf.SortType(*t)]
}
