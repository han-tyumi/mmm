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
	"os"

	"github.com/han-tyumi/mmm/config"
	"github.com/han-tyumi/mmm/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCmd = &cobra.Command{
	Use:   "remove slug...",
	Short: "Deletes and removes a mod from management by its slug",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("at least 1 slug is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			utils.Error("dependency file not found")
		}

		version := viper.GetString("version")
		fmt.Printf("using Minecraft version %s\n", version)

		depMap, err := config.DepMap()
		if err != nil {
			utils.Error(err)
		}

		ch := utils.NewErrCh(len(args))
		for i := range args {
			arg := args[i]

			go ch.Do(func() error {
				dep, ok := depMap.Get(arg)
				if !ok {
					fmt.Printf("slug, %s, not found\n", arg)
					return nil
				}

				fmt.Printf("removing %s ...\n", dep.File)
				if err := os.Remove(dep.File); err != nil {
					return err
				}
				depMap.Delete(arg)

				return nil
			})
		}

		ch.Wait(func(err error) error {
			if err != nil {
				fmt.Println(err)
			}
			return nil
		})

		if err := depMap.Write(); err != nil {
			utils.Error(err)
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
