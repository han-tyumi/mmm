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

		deps, err := config.Deps()
		if err != nil {
			utils.Error(err)
		}

		for slug, dep := range deps {
			latest, err := get.LatestFileByID(version, dep.ID, dep.Name)
			if err != nil {
				utils.Error(err)
			}

			if latest.Name == dep.File && latest.Uploaded == dep.Uploaded && latest.Size == dep.Size {
				fmt.Printf("%s up to date\n", dep.Name)
				continue
			}

			fmt.Printf("removing %s ...\n", dep.File)
			if err := os.Remove(dep.File); err != nil {
				utils.Error(err)
			}

			fmt.Printf("downloading %s ...\n", latest.Name)
			if err := download.FromURL(latest.Name, latest.URL); err != nil {
				utils.Error(err)
			}

			if err := config.SetDep(slug, dep); err != nil {
				utils.Error(err)
			}
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
