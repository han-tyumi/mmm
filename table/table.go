package table

import (
	"os"

	"github.com/han-tyumi/mcf"

	"github.com/olekukonko/tablewriter"
)

// Table returns a tablewriter.Table using the specified Format and mod data.
func Table(format Format, mods []mcf.Mod) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(format.Headers())

	for i := range mods {
		table.Append(format.Values(&mods[i]))
	}

	return table
}

// SimpleTable returns a preformatted tablewriter.Table with minimal formatting.
func SimpleTable(format Format, mods []mcf.Mod) *tablewriter.Table {
	table := Table(format, mods)

	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator(" ")
	table.SetBorder(false)
	table.SetRowLine(true)

	return table
}
