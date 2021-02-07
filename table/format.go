/*
Copyright © 2021 Matthew Champagne <mmchamp95@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package table

import (
	"fmt"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/utils"
)

// DefaultFormat is the default format used for displaying search results.
const DefaultFormat = "{id} {slug} {name} {downloads} {updated}"

var tokenHeader = map[string]string{
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

var tokenFormatter = map[string]func(*mcf.Mod) (value string){
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
// TODO: Allow this function to take in a token map to make this more generic.
func (f *Format) Headers() (headers []string) {
	var token, header string

	for _, r := range *f {
		switch {
		case token != "":
			token += string(r)

			if r != '}' {
				continue
			}

			if h, ok := tokenHeader[token]; ok {
				header += h
			} else {
				header += token
			}

			token = ""
		case r == '{':
			token += string(r)
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
// TODO: Allow this function to take in a token map to make this more generic.
func (f *Format) Values(mod *mcf.Mod) (values []string) {
	var token, value string

	for _, r := range *f {
		switch {
		case token != "":
			token += string(r)

			if r != '}' {
				continue
			}

			if format, ok := tokenFormatter[token]; ok {
				value += format(mod)
			} else {
				value += token
			}

			token = ""
		case r == '{':
			token += string(r)
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
