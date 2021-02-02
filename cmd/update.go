package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/han-tyumi/mmm/cmd/get"
	"github.com/han-tyumi/mmm/cmd/utils"

	"github.com/mitchellh/mapstructure"

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

		modList := map[string]*dependency{}
		err := viper.UnmarshalKey("mods", &modList,
			viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
		if err != nil {
			utils.Error(err)
		}

		if len(modList) == 0 {
			utils.Error("no mods being managed")
		}

		for slug, dep := range modList {
			modFile, err := get.LatestFileByID(version, dep.ID, dep.Name)
			if err != nil {
				utils.Error(err)
			}

			if modFile.Name == dep.File && modFile.Uploaded == dep.Uploaded && modFile.Size == dep.Size {
				fmt.Printf("%s up to date\n", dep.Name)
				continue
			}

			fmt.Printf("removing %s ...\n", dep.File)
			if err := os.Remove(dep.File); err != nil {
				utils.Error(err)
			}

			fmt.Printf("downloading %s ...\n", modFile.Name)
			res, err := http.Get(modFile.URL)
			if err != nil {
				utils.Error(err)
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				utils.Error(res.Status)
			}

			file, err := os.Create(modFile.Name)
			if err != nil {
				utils.Error(err)
			}
			defer file.Close()

			if _, err := io.Copy(file, res.Body); err != nil {
				utils.Error(err)
			}

			viper.Set("mods."+slug, dep)
			if err := viper.WriteConfig(); err != nil {
				utils.Error(err)
			}
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
