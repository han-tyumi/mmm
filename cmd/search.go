package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/han-tyumi/mcf"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version string
var sort sortType
var limit uint

var searchCmd = &cobra.Command{
	Use:   "search [-s sortType] [-l limit] [-v version] term...",
	Short: "Filter for mods by search terms",
	Run: func(cmd *cobra.Command, args []string) {
		version = viper.GetString("version")

		mods, err := mcf.Search(&mcf.SearchParams{
			Search:   strings.Join(args, " "),
			Sort:     mcf.SortType(sort),
			PageSize: limit,
			Version:  version,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(mods) == 0 {
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Slug", "Name", "Summary", "Downloads", "Updated"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator(" ")
		table.SetBorder(false)
		table.SetRowLine(true)
		table.SetColWidth(24)

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

	viper.BindPFlag("version", searchCmd.Flags().Lookup("version"))
}

func modRow(mod *mcf.Mod) []string {
	return []string{
		fmt.Sprint(mod.ID),
		mod.Slug,
		mod.Name,
		mod.Summary,
		downloadsToString(mod.Downloads),
		mod.Updated.Format("Jan 2 15:04 2006"),
	}
}

func downloadsToString(downloads float64) string {
	switch {
	case downloads >= 1_000_000_000:
		return fmt.Sprintf("%.1f B", downloads/1_000_000_000)
	case downloads >= 1_000_000:
		return fmt.Sprintf("%.1f M", downloads/1_000_000)
	case downloads >= 1_000:
		return fmt.Sprintf("%.1f K", downloads/1_000)
	}
	return fmt.Sprint(downloads)
}
