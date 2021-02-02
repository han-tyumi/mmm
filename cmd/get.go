package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/cmd/get"
	"github.com/han-tyumi/mmm/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var useSlug, useSearch bool

var getCmd = &cobra.Command{
	Use:   "get [-s] ...{id | slug}",
	Short: "Downloads the specified mods by ID",
	Run: func(cmd *cobra.Command, args []string) {
		version := viper.GetString("version")

		var mods []mcf.Mod
		var err error

		if useSearch {
			mods, err = get.ModsBySearch(args, version)
		} else if useSlug {
			mods, err = get.ModsBySlug(args, version)
		} else {
			if ids, err := utils.StringsToUints(args); err == nil {
				mods, err = get.ModsByID(ids)
			}
		}

		if err != nil {
			utils.Error(err)
		}

		if len(mods) == 0 {
			utils.Error("no mods found")
		}

		for i := range mods {
			mod := mods[i]

			modFile, err := get.LatestFileByMod(version, &mod)
			if err != nil {
				utils.Error(err)
			}

			name := path.Base(modFile.URL)

			fmt.Printf("downloading %s ...\n", name)
			res, err := http.Get(modFile.URL)
			if err != nil {
				utils.Error(err)
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				utils.Error(res.Status)
			}

			file, err := os.Create(name)
			if err != nil {
				utils.Error(err)
			}
			defer file.Close()

			if _, err := io.Copy(file, res.Body); err != nil {
				utils.Error(err)
			}
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&useSlug, "slug", "s", false, "Add mods based on its slug")
	getCmd.Flags().BoolVarP(&useSearch, "search", "S", false, "Add mods based on search terms")
	getCmd.Flags().StringP("version", "v", "", "Download the latest for a Minecraft version")

	viper.BindPFlag("version", getCmd.Flags().Lookup("version"))
}
