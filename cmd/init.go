package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a .mmm.yml dependency file with a Minecraft version if it does not exist",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("a Minecraft version argument is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.SafeWriteConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		viper.Set("version", args[0])
		if err := viper.WriteConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
