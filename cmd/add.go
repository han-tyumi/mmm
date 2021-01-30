package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/han-tyumi/mcf"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type dependency struct {
	Name string    `mapstructure:"name"`
	Date time.Time `mapstructure:"date"`
	Size uint      `mapstructure:"size"`
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Gets Minecraft CurseForge mods by ID and adds them to your dependency configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg.ConfigFileUsed() == "" {
			fmt.Fprintln(os.Stderr, "configuration file not found")
			os.Exit(1)
		}

		version = cfg.GetString("version")
		fmt.Printf("using Minecraft version %s\n", version)

		var mods []mcf.Mod
		var err error

		if useSearch {
			mods, err = modsBySearch(args)
		} else {
			mods, err = modsByID(args)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(mods) == 0 {
			fmt.Fprintln(os.Stderr, "no mods found")
			os.Exit(1)
		}

		for i := range mods {
			mod := mods[i]

			modFile, err := findFile(&mod)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			dep := &dependency{
				Name: modFile.Name,
				Date: modFile.Uploaded,
				Size: modFile.Size,
			}

			key := "mods>" + mod.Name
			if cfg.IsSet(key) {
				prev := &dependency{}
				err := cfg.UnmarshalKey(key, prev,
					viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				if prev.Name != dep.Name || prev.Date != dep.Date || prev.Size != dep.Size {
					fmt.Printf("removing %s ...\n", prev.Name)
					if err := os.Remove(prev.Name); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				} else {
					fmt.Printf("skipping %s\n", mod.Name)
					continue
				}
			}

			fmt.Printf("downloading %s ...\n", modFile.Name)
			res, err := http.Get(modFile.URL)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				fmt.Fprintln(os.Stderr, res.Status)
				os.Exit(1)
			}

			file, err := os.Create(modFile.Name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer file.Close()

			if _, err := io.Copy(file, res.Body); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			cfg.Set(key, dep)
			if err := cfg.WriteConfig(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolVarP(&useSearch, "search", "s", false, "Add mods based on search terms")
}
