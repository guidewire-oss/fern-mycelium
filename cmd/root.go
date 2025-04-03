package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command for the `mycel` CLI.
var rootCmd = &cobra.Command{
	Use:   "mycel",
	Short: "mycel is the CLI for fern-mycelium, the test intelligence layer",
	Long: `mycel is the command-line interface for fern-mycelium, 
an intelligent context engine that enhances your test ecosystem 
with insights, analytics, and AI agents using the Model Context Protocol (MCP).

Use this CLI to run the mycelium API server, interact with test data,
query agents like the Test Coach or Postmortem Generator, and more.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to custom config file (optional)")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable verbose debug logging")
}
