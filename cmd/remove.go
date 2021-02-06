package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/han-tyumi/mmm/config"
	"github.com/han-tyumi/mmm/utils"

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

		depMap, err := config.DepMap()
		if err != nil {
			utils.Error(err)
		}

		ch := utils.NewErrCh(len(args))
		for i := range args {
			arg := args[i]

			go ch.Do(func() error {
				dep, ok := depMap.Get(arg)
				if !ok {
					fmt.Printf("slug, %s, not found\n", arg)
					return nil
				}

				fmt.Printf("removing %s ...\n", dep.File)
				if err := os.Remove(dep.File); err != nil {
					return err
				}
				depMap.Delete(arg)

				return nil
			})
		}

		ch.Wait(func(err error) error {
			if err != nil {
				fmt.Println(err)
			}
			return nil
		})

		if err := depMap.Write(); err != nil {
			utils.Error(err)
		}

		fmt.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
