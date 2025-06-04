package tools

import (
	"context"
	"fmt"
	"gaia-mcp-go/internal/api"
	"gaia-mcp-go/pkg/imageutil"
	"gaia-mcp-go/pkg/shared"

	"github.com/mark3labs/mcp-go/mcp"
)

type UpscalerTool struct {
	api  api.GaiaApi
	tool mcp.Tool
}

func NewUpscalerTool(api api.GaiaApi) *UpscalerTool {
	return &UpscalerTool{
		api: api,
		tool: mcp.NewTool(
			"upscaler",
			mcp.WithDescription("Enhance the resolution quality of images"),
			mcp.WithString(
				"image_url",
				mcp.Required(),
				mcp.Description("The image URL to upscale. It must be GAIA's image url: starts with `https://cdn.protogaia.com/`"),
			),
			mcp.WithNumber(
				"ratio",
				mcp.DefaultNumber(2),
				mcp.Min(1),
				mcp.Max(4),
				mcp.Description("The ratio to upscale the image. It must be a number between 1 and 4"),
			),
		),
	}
}

func (t *UpscalerTool) ToolName() string {
	return "upscaler"
}

func (t *UpscalerTool) MCPTool() mcp.Tool {
	return t.tool
}

func (t *UpscalerTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	imageUrl := args["image_url"]
	ratio := args["ratio"]

	res, err := t.api.GenerateImages(ctx, api.GenerateImagesRequest{
		RecipeId: shared.RecipeIdUpscaler,
		Params: map[string]interface{}{
			"image":         imageUrl,
			"upscale_mode":  "4x-Ultrasharp.pt",
			"upscale_ratio": ratio,
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

	msg := fmt.Sprintf("Upscaled successfully. Image url: %s", res.Images[0])

	return mcp.NewToolResultImage(msg, base64Data, mimeType), nil
}
