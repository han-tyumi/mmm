package search

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/han-tyumi/mcf"
)

// SortType is a wrapper of mcf.SortType that implements the pflag.Value interface.
type SortType mcf.SortType

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

// Set sets the value of the SortType for a given string argument.
func (t *SortType) Set(s string) error {
	p := regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(s, "")
	p = strings.ToLower(p)

	if sort, ok := nameToSortType[p]; ok {
		*t = SortType(sort)
		return nil
	}

	return fmt.Errorf("%s is not a valid sort type", s)
}

func (t *SortType) String() string {
	return sortTypeToName[mcf.SortType(*t)]
}

// Type returns the type name for SortType.
func (t *SortType) Type() string {
	return "sortType"
}
