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
	"errors"
	"fmt"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/download"
	"github.com/han-tyumi/mmm/get"
	"github.com/han-tyumi/mmm/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version string

var getCmd = &cobra.Command{
	Use:   "get {id | slug}...",
	Short: "Downloads unmanaged mods to the current working directory by slug or ID",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no arguments specified")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			version = viper.GetString("version")
		}

		if version != "" {
			fmt.Printf("using Minecraft version %s\n", version)
		}

		if err := get.LatestFileForEachArg(args, version, func(_ *mcf.Mod, latest *mcf.ModFile) error {
			fmt.Printf("downloading %s ...\n", latest.Name)
			return download.FromURL(latest.Name, latest.URL)
		}); err != nil {
			utils.Error(err)
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&version, "version", "v", "", "Minecraft version to download latest files for")
}
