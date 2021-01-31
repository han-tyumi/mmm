package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates all managed mods",
	Run: func(cmd *cobra.Command, args []string) {
		if viper.ConfigFileUsed() == "" {
			fmt.Fprintln(os.Stderr, "dependency file not found")
			os.Exit(1)
		}

		version = viper.GetString("version")
		fmt.Printf("using Minecraft version %s\n", version)

		modList := map[string]*dependency{}
		err := viper.UnmarshalKey("mods", &modList,
			viper.DecodeHook(mapstructure.StringToTimeHookFunc(time.RFC3339)))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(modList) == 0 {
			fmt.Fprintln(os.Stderr, "no mods being managed")
			os.Exit(1)
		}

		for slug, dep := range modList {
			modFile, err := findLatestByID(dep.ID, dep.Name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if modFile.Name == dep.File && modFile.Uploaded == dep.Uploaded && modFile.Size == dep.Size {
				fmt.Printf("%s up to date\n", dep.Name)
				continue
			}

			fmt.Printf("removing %s ...\n", dep.File)
			if err := os.Remove(dep.File); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
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

			viper.Set("mods."+slug, dep)
			if err := viper.WriteConfig(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
