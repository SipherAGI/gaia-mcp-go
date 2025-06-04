package stdio

import (
	"gaia-mcp-go/internal/api"
	"gaia-mcp-go/internal/tools"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var (
	StdioCmd = &cobra.Command{
		Use:   "stdio",
		Short: "Run the Gaia MCP server in stdio mode",
		Long:  `Run the Gaia MCP server in stdio mode.`,
		Run:   runStdio,
	}

	ServerName = "gaia-mcp-server"
)

func init() {
	StdioCmd.Flags().StringP("api-key", "k", "", "The API key to use for the Gaia MCP server")
}

func runStdio(cmd *cobra.Command, args []string) {
	// Get the API key from the args
	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		slog.Error("Failed to get API key", "error", err)
		os.Exit(1)
	}

	// Create the API client
	apiClient := api.NewGaiaApi(api.GaiaApiConfig{
		BaseUrl: "https://api.protogaia.com",
		ApiKey:  apiKey,
	})

	// Create the tools
	generateImageTool := tools.NewGenerateImageTool(apiClient)
	faceEnhancerTool := tools.NewFaceEnhancerTool(apiClient)
	remixTool := tools.NewRemixTool(apiClient)
	upscalerTool := tools.NewUpscalerTool(apiClient)
	uploadImageTool := tools.NewUploadImageTool(apiClient)

	// Create the server
	s := server.NewMCPServer(
		ServerName,
		cmd.Version,
		server.WithToolCapabilities(false),
	)

	// Add the tools to the server
	s.AddTool(generateImageTool.MCPTool(), generateImageTool.Handler)
	s.AddTool(faceEnhancerTool.MCPTool(), faceEnhancerTool.Handler)
	s.AddTool(remixTool.MCPTool(), remixTool.Handler)
	s.AddTool(upscalerTool.MCPTool(), upscalerTool.Handler)
	s.AddTool(uploadImageTool.MCPTool(), uploadImageTool.Handler)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		slog.Error("Failed to serve stdio", "error", err)
	}
}
