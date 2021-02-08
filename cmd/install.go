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
	"os"

	"github.com/han-tyumi/mmm/config"
	"github.com/han-tyumi/mmm/download"
	"github.com/han-tyumi/mmm/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var force bool

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs all mods being managed within a configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			utils.Error("dependency file not found")
		}

		depMap, err := config.DepMapSync()
		if err != nil {
			utils.Error(err)
		}

		ch := utils.NewErrCh(len(depMap))
		for _, dep := range depMap {
			dep := dep

			go ch.Do(func() error {
				if !force {
					info, err := os.Stat(dep.File)
					if err == nil && info.Size() == int64(dep.Size) {
						fmt.Printf("%s already installed\n", dep.Name)
						return nil
					}
				}

				fmt.Printf("downloading %s ...\n", dep.Name)
				return download.FromURL(dep.File, dep.URL)
			})
		}

		if err := ch.Wait(func(err error) error {
			return err
		}); err != nil {
			utils.Error(err)
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing mods")
}
