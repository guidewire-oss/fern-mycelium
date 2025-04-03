package cmd

import (
	"fmt"
	"github.com/guidewire-oss/fern-mycelium/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the fern-mycelium MCP API server",
	Long:  "Launches the fern-mycelium server exposing GraphQL and REST APIs for test context and agent interaction.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌱 Starting Mycelium MCP API server...")
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
