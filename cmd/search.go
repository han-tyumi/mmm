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

package cmd

import (
	"strings"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/cmd/search"
	"github.com/han-tyumi/mmm/table"
	"github.com/han-tyumi/mmm/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sort = search.SortType(mcf.Featured)
var limit uint
var format string

var searchCmd = &cobra.Command{
	Use:   "search [terms]...",
	Short: "Displays search results for Minecraft CurseForge mods",
	Long: strings.ReplaceAll(`#### Sort Types
- ^featured, feat, f, 0^
- ^popularity, pop, p, 1^
- ^lastupdate, update, up, u, last, l, 2^
- ^name, n, 3^
- ^author, auth, a, 4^
- ^totaldownloads, downloads, down, d, total, t, 5^

#### Table Format Tokens
- ^{id}^
- ^{slug}^
- ^{name}^
- ^{language}^
- ^{url}^
- ^{rank}^
- ^{popularity}^
- ^{downloads}^
- ^{updated}^
- ^{released}^
- ^{created}^`, "^", "`"),
	Run: func(cmd *cobra.Command, args []string) {
		version := viper.GetString("version")

		mods, err := mcf.Search(&mcf.SearchParams{
			Search:   strings.Join(args, " "),
			Sort:     mcf.SortType(sort),
			PageSize: limit,
			Version:  version,
		})
		if err != nil {
			utils.Error(err)
		}

		if len(mods) == 0 {
			return
		}

		table := table.SimpleTable(table.Format(format), mods)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("version", "v", "", "Minecraft version to filter by")
	searchCmd.Flags().VarP(&sort, "sort", "s", "how to sort mod results")
	searchCmd.Flags().UintVarP(&limit, "limit", "l", 5, "how many results to return")
	searchCmd.Flags().StringVarP(&format, "format", "f", table.DefaultFormat, "table format to use")

	viper.BindPFlag("version", searchCmd.Flags().Lookup("version"))
}
