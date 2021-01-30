package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a mmm dependency configuration file if it doesn't exist",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("a single Minecraft version argument is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := cfg.SafeWriteConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		cfg.Set("version", args[0])
		if err := cfg.WriteConfig(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
