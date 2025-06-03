package tools

import (
	"context"
	"fmt"
	"gaia-mcp-go/internal/api"
	"gaia-mcp-go/pkg/imageutil"
	"gaia-mcp-go/pkg/shared"

	"github.com/mark3labs/mcp-go/mcp"
)

type RemixTool struct {
	api  api.GaiaApi
	tool mcp.Tool
}

func NewRemixTool(
	api api.GaiaApi,
) *RemixTool {
	return &RemixTool{
		api: api,
		tool: mcp.NewTool(
			"remix",
			mcp.WithDescription("Remix an image with a prompt"),
			mcp.WithString(
				"inputImage",
				mcp.Required(),
				mcp.Description("The image URL to remix. It must be GAIA's image url: starts with `https://cdn.protogaia.com/`"),
			),
			mcp.WithString(
				"variationControl",
				mcp.Description("The variation control of the remix. One of the following: 'subtle', 'medium', 'strong'"),
				mcp.DefaultString("subtle"),
				mcp.Enum("subtle", "medium", "strong"),
			),
		),
	}
}

func (t *RemixTool) ToolName() string {
	return "remix"
}

func (t *RemixTool) MCPTool() mcp.Tool {
	return t.tool
}

func (t *RemixTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	inputImage := args["inputImage"]
	variationControl := args["variationControl"]

	res, err := t.api.GenerateImages(ctx, api.GenerateImagesRequest{
		RecipeId: shared.RecipeIdRemix,
		Params: map[string]interface{}{
			"inputImage":       inputImage,
			"variationControl": variationControl,
			"numberOfImages":   1, // Always generate 1 image
		},
	})

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if !res.Success {
		return mcp.NewToolResultError(*res.Error), nil
	}

	if res.Error != nil {
		return mcp.NewToolResultError(*res.Error), nil
	}

	if len(res.Images) == 0 {
		return mcp.NewToolResultError("No images were generated. Please try again."), nil
	}

	base64Data, mimeType, err := imageutil.ProcessImageQuickForMCP(ctx, res.Images[0])
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process image: %v", err)), nil
	}

	msg := fmt.Sprintf("Remix generated successfully. Image url: %s", res.Images[0])

	return mcp.NewToolResultImage(msg, base64Data, mimeType), nil
}
