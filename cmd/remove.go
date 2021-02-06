package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/han-tyumi/mmm/cmd/config"
	"github.com/han-tyumi/mmm/cmd/utils"

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
			utils.Error("dependency file not found")
		}

		version := viper.GetString("version")
		fmt.Printf("using Minecraft version %s\n", version)

		deps, err := config.Deps()
		if err != nil {
			utils.Error(err)
		}

		for i := range args {
			arg := args[i]

			dep, ok := deps[arg]
			if !ok {
				fmt.Printf("slug, %s, not found\n", arg)
				continue
			}

			fmt.Printf("removing %s ...\n", dep.File)
			if err := os.Remove(dep.File); err != nil {
				utils.Error(err)
			}

			delete(deps, arg)

			if err := config.SetDeps(deps); err != nil {
				utils.Error(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
