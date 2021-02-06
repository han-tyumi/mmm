package cmd

import (
	"fmt"
	"os"

	"github.com/han-tyumi/mmm/config"
	"github.com/han-tyumi/mmm/download"
	"github.com/han-tyumi/mmm/get"
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
				latest, err := get.LatestFileByID(version, dep.ID, dep.Name)
				if err != nil {
					return err
				}

				if latest.Name == dep.File && latest.Uploaded == dep.Uploaded && latest.Size == dep.Size {
					fmt.Printf("%s up to date\n", dep.Name)
					return nil
				}

				fmt.Printf("removing %s ...\n", dep.File)
				if err := os.Remove(dep.File); err != nil {
					return err
				}

				fmt.Printf("downloading %s ...\n", latest.Name)
				if err := download.FromURL(latest.Name, latest.URL); err != nil {
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
