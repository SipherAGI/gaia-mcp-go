package interfaces

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type GaiaTool interface {
	ToolName() string
	MCPTool() mcp.Tool
	Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
