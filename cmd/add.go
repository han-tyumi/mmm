package cmd

import (
	"errors"
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

			dep := &dependency{
				ID:       mod.ID,
				Name:     mod.Name,
				File:     modFile.Name,
				Uploaded: modFile.Uploaded,
				Size:     modFile.Size,
			}

			key := "mods." + mod.Slug
			if viper.IsSet(key) {
				prev := &dependency{}
				err := viper.UnmarshalKey(key, prev,
					viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
				if err != nil {
					utils.Error(err)
				}

				if prev.File != dep.File || prev.Uploaded != dep.Uploaded || prev.Size != dep.Size {
					fmt.Printf("removing %s ...\n", prev.File)
					if err := os.Remove(prev.File); err != nil {
						utils.Error(err)
					}
				}

				if info, err := os.Stat(dep.File); err == nil && info.Size() == int64(dep.Size) {
					fmt.Printf("skipping %s\n", mod.Name)
					continue
				}
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

			viper.Set(key, dep)
			if err := viper.WriteConfig(); err != nil {
				utils.Error(err)
			}
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
