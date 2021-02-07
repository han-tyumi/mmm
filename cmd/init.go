package cmd

import (
	"errors"

	"github.com/han-tyumi/mmm/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init version",
	Short: "Initializes a mod dependency file using a Minecraft version",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("a Minecraft version argument is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.SafeWriteConfig(); err != nil {
			utils.Error(err)
		}

		viper.Set("version", args[0])
		if err := viper.WriteConfig(); err != nil {
			utils.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
