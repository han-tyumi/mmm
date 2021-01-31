package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCmd = &cobra.Command{
	Use:   "remove slug...",
	Short: "Deletes and removes a mod from management by its slug",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("at least 1 slug is required")
		}
		return nil
	},
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

		for i := range args {
			arg := args[i]

			dep, ok := modList[arg]
			if !ok {
				fmt.Printf("slug, %s, not found\n", arg)
				continue
			}

			fmt.Printf("removing %s ...\n", dep.File)
			if err := os.Remove(dep.File); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			delete(modList, arg)

			viper.Set("mods", &modList)
			if err := viper.WriteConfig(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
