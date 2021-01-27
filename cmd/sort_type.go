package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/han-tyumi/mcf/api"
)

var nameToSortType = map[string]api.SortType{
	"0":        api.Featured,
	"f":        api.Featured,
	"feat":     api.Featured,
	"featured": api.Featured,

	"1":          api.Popularity,
	"p":          api.Popularity,
	"pop":        api.Popularity,
	"popularity": api.Popularity,

	"2":          api.LastUpdate,
	"l":          api.LastUpdate,
	"last":       api.LastUpdate,
	"u":          api.LastUpdate,
	"up":         api.LastUpdate,
	"update":     api.LastUpdate,
	"lastupdate": api.LastUpdate,

	"3":    api.Name,
	"n":    api.Name,
	"name": api.Name,

	"4":      api.Author,
	"a":      api.Author,
	"auth":   api.Author,
	"author": api.Author,

	"5":              api.TotalDownloads,
	"t":              api.TotalDownloads,
	"total":          api.TotalDownloads,
	"d":              api.TotalDownloads,
	"down":           api.TotalDownloads,
	"downloads":      api.TotalDownloads,
	"totaldownloads": api.TotalDownloads,
}

var sortTypeToName = map[api.SortType]string{
	api.Featured:       "featured",
	api.Popularity:     "popularity",
	api.LastUpdate:     "lastupdate",
	api.Name:           "name",
	api.Author:         "author",
	api.TotalDownloads: "totaldownloads",
}

type sortType api.SortType

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
	return sortTypeToName[api.SortType(*t)]
}
