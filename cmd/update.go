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
	"fmt"

	"github.com/han-tyumi/mmm/config"
	"github.com/han-tyumi/mmm/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates all managed mods",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			utils.Error("dependency file not found")
		}

		version := viper.GetString("version")
		fmt.Printf("using Minecraft version %s\n", version)

		depMap, err := config.DepMapSync()
		if err != nil {
			utils.Error(err)
		}

		ch := utils.NewErrCh(len(depMap))
		for slug, dep := range depMap {
			slug := slug
			dep := dep

			go ch.Do(func() error {
				latest, err := dep.LatestFile(version)
				if err != nil {
					return err
				}

				if dep.SameFile(latest) {
					fmt.Printf("%s up to date\n", dep.Name)
					return nil
				}

				fmt.Printf("removing %s ...\n", dep.File)
				if err := dep.RemoveFile(); err != nil {
					return err
				}

				dep.UpdateFile(latest)

				fmt.Printf("downloading %s ...\n", latest.Name)
				if err := dep.Download(); err != nil {
					return err
				}

				return config.SetDep(slug, dep)
			})
		}

		ch.Wait(func(err error) error {
			if err != nil {
				utils.Error(err)
			}
			return err
		})

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
