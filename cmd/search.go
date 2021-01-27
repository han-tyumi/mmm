package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/han-tyumi/mcf/api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var version string
var sort sortType
var limit uint

var searchCmd = &cobra.Command{
	Use:   "search [-s sortType] [-l limit] [-v version] term...",
	Short: "Filter for mods by search terms",
	Run: func(cmd *cobra.Command, args []string) {
		mods, err := api.Many(&api.ManyParams{
			Search:   strings.Join(args, " "),
			Sort:     api.SortType(sort),
			PageSize: limit,
			Version:  version,
		})
		if err != nil {
			panic(err)
		}

		if len(mods) == 0 {
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Summary", "Downloads", "Updated"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator(" ")
		table.SetBorder(false)
		table.SetRowLine(true)

		for i := range mods {
			mod := mods[i]
			table.Append(modRow(&mod))
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&version, "version", "v", "", "Minecraft version to filter by")
	searchCmd.Flags().VarP(&sort, "sort", "s", "How to sort mod results")
	searchCmd.Flags().UintVarP(&limit, "limit", "l", 5, "How many results to return")
}

func modRow(mod *api.Mod) []string {
	return []string{
		fmt.Sprint(mod.ID),
		mod.Name,
		mod.Summary,
		fmt.Sprint(mod.Downloads),
		mod.Updated,
	}
}
