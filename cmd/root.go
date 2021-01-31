package cmd

import (
	"fmt"
	"os"

	"github.com/han-tyumi/mmm/cmd/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cwd string

var rootCmd = &cobra.Command{
	Use:   "mmm",
	Short: "Manages Minecraft CurseForge mods",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Error(err)
	}
}

func init() {
	cobra.OnInitialize(cobraInit)

	rootCmd.PersistentFlags().StringVarP(&cwd, "cwd", "C", "", "change the working directory")
}

func cobraInit() {
	if cwd != "" {
		if err := os.Chdir(cwd); err != nil {
			utils.Error(err)
		}
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("mmm")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("using config file:", viper.ConfigFileUsed())
	}
}
