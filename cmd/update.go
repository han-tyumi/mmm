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

		batch := false
		viperVersion := viper.GetString("version")

		if version != "" && version != viperVersion {
			batch = true
			fmt.Printf("updating mods from %s to %s ...\n", viperVersion, version)
		} else {
			version = viperVersion
			fmt.Printf("updating mods to use latest %s files ...\n", version)
		}

		depMap, err := config.DepMap()
		if err != nil {
			utils.Error(err)
		}
		clone := depMap.Clone()

		ch := utils.NewErrCh(depMap.Len())
		depMap.Each(func(slug string, dep *config.Dependency) {
			go ch.Do(func() error {
				latest, err := dep.LatestFile(version)
				if err != nil {
					return fmt.Errorf("%s: %s", slug, err)
				}

				if dep.SameFile(latest) {
					fmt.Printf("%s up to date\n", dep.Name)
					return nil
				}

				fmt.Printf("removing %s ...\n", dep.File)
				if err := dep.RemoveFile(); err != nil {
					return fmt.Errorf("%s: %s", dep.File, err)
				}

				dep.UpdateFile(latest)

				fmt.Printf("downloading %s ...\n", latest.Name)
				if err := dep.Download(); err != nil {
					return fmt.Errorf("%s: %s", latest.Name, err)
				}

				if batch {
					depMap.Set(slug, dep)
					return nil
				}

				if err := config.SetDep(slug, dep); err != nil {
					return fmt.Errorf("%s: %s", slug, err)
				}
				return nil
			})
		})

		revertUpdates := func(err error) {
			fmt.Fprintln(os.Stderr, err)
			fmt.Println("reverting updates; you may need to reinstall ...")
			if err := clone.Write(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(1)
		}

		ch.Wait(func(err error) error {
			if err != nil {
				revertUpdates(err)
			}
			return err
		})

		if batch {
			viper.Set("version", version)
			if err := depMap.Write(); err != nil {
				viper.Set("version", viperVersion)
				revertUpdates(err)
			}
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&version, "version", "v", "", "Minecraft version to update mods to")
}
