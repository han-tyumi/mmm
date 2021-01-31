package cmd

import (
	"strings"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/cmd/search"
	"github.com/han-tyumi/mmm/cmd/table"
	"github.com/han-tyumi/mmm/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sort search.SortType
var limit uint
var format string

var searchCmd = &cobra.Command{
	Use:   "search [-s sortType] [-l limit] [-v version] [-f tableFormat] term...",
	Short: "Filter for mods by search terms",
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
	searchCmd.Flags().VarP(&sort, "sort", "s", "How to sort mod results")
	searchCmd.Flags().UintVarP(&limit, "limit", "l", 5, "How many results to return")
	searchCmd.Flags().StringVarP(&format, "format", "f", table.DefaultFormat, "Table format to use")

	viper.BindPFlag("version", searchCmd.Flags().Lookup("version"))
}
