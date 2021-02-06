package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/config"
	"github.com/han-tyumi/mmm/download"
	"github.com/han-tyumi/mmm/get"
	"github.com/han-tyumi/mmm/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add {id | slug}...",
	Short: "Gets Minecraft CurseForge mods by ID or slug and adds them to your dependency file",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no arguments specified")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			utils.Error("dependency file not found")
		}

		version := viper.GetString("version")
		fmt.Printf("using Minecraft version %s\n", version)

		if err := get.LatestFileForEachArg(args, version, func(mod *mcf.Mod, latest *mcf.ModFile) error {
			dep := &config.Dependency{
				ID:       mod.ID,
				Name:     mod.Name,
				File:     latest.Name,
				Uploaded: latest.Uploaded,
				Size:     latest.Size,
			}

			if prev, err := config.Dep(mod.Slug); err == nil {
				if prev.File != dep.File || prev.Uploaded != dep.Uploaded || prev.Size != dep.Size {
					fmt.Printf("removing %s ...\n", prev.File)
					if err := os.Remove(prev.File); err != nil {
						return err
					}
				}

				if info, err := os.Stat(dep.File); err == nil && info.Size() == int64(dep.Size) {
					fmt.Printf("skipping %s\n", mod.Name)
					return nil
				}
			}

			fmt.Printf("downloading %s ...\n", latest.Name)
			if err := download.FromURL(latest.Name, latest.URL); err != nil {
				return err
			}

			if err := config.SetDep(mod.Slug, dep); err != nil {
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
	rootCmd.AddCommand(addCmd)
}
