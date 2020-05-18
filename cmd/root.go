package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dot",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		fmt.Println("Current location of your dotfiles:", viper.Get("dotfiles.path"))
		dotfilePath, err := homedir.Expand(viper.Get("dotfiles.path").(string))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		backupPath, err := homedir.Expand(viper.Get("dotfiles.backup").(string))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var backupFiles []string

		fmt.Println("Current location of your dotfiles:", dotfilePath)

		fmt.Println("We're currently backing up:")

		filecount := 0
		err = filepath.Walk(backupPath, visit(&backupFiles))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, file := range backupFiles {
			if filecount == 0 {
				filecount++
				continue
			}
			fmt.Println("\t", filecount, ":", file)
			filecount++
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dot.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dot.toml")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		*files = append(*files, path)
		return nil
	}
}
