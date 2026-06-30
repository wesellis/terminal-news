package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wesellis/terminal-news/cli/internal/config"
	"github.com/wesellis/terminal-news/cli/internal/ui"
)

var (
	cfgFile string
	offline bool
)

var rootCmd = &cobra.Command{
	Use:   "terminal-news",
	Short: "AM Radio for the Information Age",
	Long: `Terminal News - A terminal-native news aggregator with community curation,
local classifieds, and real-time weather.

Navigate with keyboard shortcuts, vote on articles, post classifieds,
and stay informed - all from your terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		runApp()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Terminal News v0.1.0")
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.terminal-news/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&offline, "offline", "o", false, "run in offline mode (uses cached data)")

	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error finding home directory:", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home + "/.terminal-news")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// Create default config if it doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := config.CreateDefault(); err != nil {
				log.Printf("Warning: Could not create default config: %v\n", err)
			}
		}
	}
}

func runApp() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Override offline mode if flag is set
	if offline {
		cfg.Offline = true
	}

	// Initialize the app
	app, err := ui.NewApp(cfg)
	if err != nil {
		fmt.Printf("Error initializing app: %v\n", err)
		os.Exit(1)
	}

	// Create the Bubble Tea program
	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running Terminal News: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
