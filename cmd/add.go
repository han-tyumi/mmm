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
	"github.com/han-tyumi/mmm/config"
	"github.com/han-tyumi/mmm/get"
	"github.com/han-tyumi/mmm/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add {id | slug}...",
	Short: "Downloads and adds mods to your dependency file by slug or ID",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no arguments specified")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			utils.Error("dependency file not found")
		}

		version := viper.GetString("version")
		fmt.Printf("using Minecraft version %s\n", version)

		if err := get.LatestFileForEachArg(args, version, func(mod *mcf.Mod, latest *mcf.ModFile) error {
			dep := &config.Dependency{
				ID:       mod.ID,
				Name:     mod.Name,
				URL:      latest.URL,
				File:     latest.Name,
				Uploaded: latest.Uploaded,
				Size:     latest.Size,
			}

			if prev, err := config.Dep(mod.Slug); err == nil {
				// remove older/previous files
				if !dep.SameDepFile(prev) {
					fmt.Printf("removing %s ...\n", prev.File)
					if err := prev.RemoveFile(); err != nil {
						return err
					}
				}

				// skip already downloaded files
				if downloaded, _ := dep.Downloaded(); downloaded {
					fmt.Printf("%s already added\n", dep.Name)
					return nil
				}
			}

			fmt.Printf("downloading %s ...\n", dep.File)
			if err := dep.Download(); err != nil {
				return err
			}

			return config.SetDep(mod.Slug, dep)
		}); err != nil {
			utils.Error(err)
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
