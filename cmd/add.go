package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/han-tyumi/mcf"
	"github.com/han-tyumi/mmm/cmd/download"
	"github.com/han-tyumi/mmm/cmd/get"
	"github.com/han-tyumi/mmm/cmd/utils"

	"github.com/mitchellh/mapstructure"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type dependency struct {
	ID       uint      `mapstructure:"id"`
	Name     string    `mapstructure:"name"`
	File     string    `mapstructure:"file"`
	Uploaded time.Time `mapstructure:"uploaded"`
	Size     uint      `mapstructure:"size"`
}

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
			dep := &dependency{
				ID:       mod.ID,
				Name:     mod.Name,
				File:     latest.Name,
				Uploaded: latest.Uploaded,
				Size:     latest.Size,
			}

			key := "mods." + mod.Slug
			if viper.IsSet(key) {
				prev := &dependency{}
				err := viper.UnmarshalKey(key, prev,
					viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
				if err != nil {
					return err
				}

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

			viper.Set(key, dep)
			if err := viper.WriteConfig(); err != nil {
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
