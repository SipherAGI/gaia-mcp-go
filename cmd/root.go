package cmd

import (
	"context"
	"gaia-mcp-go/cmd/stdio"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"

	rootCmd = &cobra.Command{
		Use:   "gaia-mcp-server",
		Short: "Gaia MCP Server",
		Long:  `Gaia MCP Server is a MCP server for the ProtoGaia project.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Version: Version,
	}
)

func Execute() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		slog.Info("Received signal to terminate", "signal", <-c)
		cancel()
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error("Failed to execute root command", "error", err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(stdio.StdioCmd)
}
