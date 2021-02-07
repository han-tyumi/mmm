/*
Package table provides functionality for displaying formatted tables based on tokens.

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
