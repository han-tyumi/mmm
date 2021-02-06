package cmd

import (
	"errors"
	"fmt"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/download"
	"github.com/han-tyumi/mmm/get"
	"github.com/han-tyumi/mmm/utils"

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

		if err := get.LatestFileForEachArg(args, version, func(_ *mcf.Mod, latest *mcf.ModFile) error {
			fmt.Printf("downloading %s ...\n", latest.Name)
			if err := download.FromURL(latest.Name, latest.URL); err != nil {
				return err
			}
			return nil
		}); err != nil {
			utils.Error(err)
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("version", "v", "", "Download the latest for a Minecraft version")

	viper.BindPFlag("version", getCmd.Flags().Lookup("version"))
}
