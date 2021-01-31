package table

import (
	"fmt"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/cmd/utils"
)

// DefaultFormat is the default format used for displaying search results.
const DefaultFormat = "{id} {slug} {name} {downloads} {updated}"

var tagHeader = map[string]string{
	"{id}":         "ID",
	"{slug}":       "Slug",
	"{name}":       "Name",
	"{language}":   "Language",
	"{url}":        "URL",
	"{rank}":       "Rank",
	"{popularity}": "Popularity",
	"{downloads}":  "Downloads",
	"{updated}":    "Updated",
	"{released}":   "Released",
	"{created}":    "Created",
}

var tagFormatter = map[string]func(*mcf.Mod) (value string){
	"{id}":       func(mod *mcf.Mod) string { return fmt.Sprint(mod.ID) },
	"{slug}":     func(mod *mcf.Mod) string { return mod.Slug },
	"{name}":     func(mod *mcf.Mod) string { return mod.Name },
	"{language}": func(mod *mcf.Mod) string { return mod.Language },
	"{url}":      func(mod *mcf.Mod) string { return mod.URL },
	"{rank}":     func(mod *mcf.Mod) string { return fmt.Sprint(mod.Rank) },
	"{popularity}": func(mod *mcf.Mod) string {
		return utils.FormatBigFloat(mod.Popularity)
	},
	"{downloads}": func(mod *mcf.Mod) string {
		return utils.FormatBigFloat(mod.Downloads)
	},
	"{updated}": func(mod *mcf.Mod) string {
		return mod.Updated.Format("Jan 2 15:04 2006")
	},
	"{released}": func(mod *mcf.Mod) string {
		return mod.Released.Format("Jan 2 15:04 2006")
	},
	"{created}": func(mod *mcf.Mod) string {
		return mod.Created.Format("Jan 2 15:04 2006")
	},
}

// Format is used to represent the desired mod table format to use through a string.
type Format string

// Headers returns the table header names for the Format.
func (f *Format) Headers() (headers []string) {
	var tag, header string

	for _, r := range *f {
		switch {
		case tag != "":
			tag += string(r)

			if r != '}' {
				continue
			}

			if h, ok := tagHeader[tag]; ok {
				header += h
			} else {
				header += tag
			}

			tag = ""
		case r == '{':
			tag += string(r)
		case r == ' ':
			if header == "" {
				continue
			}

			headers = append(headers, header)
			header = ""
		default:
			header += string(r)
		}
	}

	if header != "" {
		headers = append(headers, header)
	}

	return
}

// Values returns the table values for a given mod and Format.
func (f *Format) Values(mod *mcf.Mod) (values []string) {
	var tag, value string

	for _, r := range *f {
		switch {
		case tag != "":
			tag += string(r)

			if r != '}' {
				continue
			}

			if format, ok := tagFormatter[tag]; ok {
				value += format(mod)
			} else {
				value += tag
			}

			tag = ""
		case r == '{':
			tag += string(r)
		case r == ' ':
			if value == "" {
				continue
			}

			values = append(values, value)
			value = ""
		default:
			value += string(r)
		}
	}

	if value != "" {
		values = append(values, value)
	}

	return
}
