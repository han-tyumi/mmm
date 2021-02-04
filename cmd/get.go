package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/han-tyumi/mmm/cmd/get"
	"github.com/han-tyumi/mmm/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:   "get [-v version] ...{id | slug}",
	Short: "Downloads the specified mods by ID or slug",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no arguments specified")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		version := viper.GetString("version")
		fmt.Printf("using Minecraft version %s\n", version)

		mods, err := get.ModsByArgs(args, version)
		if err != nil {
			utils.Error(err)
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

	getCmd.Flags().StringP("version", "v", "", "Download the latest for a Minecraft version")

	viper.BindPFlag("version", getCmd.Flags().Lookup("version"))
}
