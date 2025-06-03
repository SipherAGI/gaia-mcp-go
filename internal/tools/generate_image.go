package tools

import (
	"context"
	"fmt"
	"gaia-mcp-go/internal/api"
	"gaia-mcp-go/pkg/imageutil"
	"gaia-mcp-go/pkg/shared"

	"github.com/mark3labs/mcp-go/mcp"
)

// GenerateImageTool implements the GaiaTool interface
type GenerateImageTool struct {
	api  api.GaiaApi
	tool mcp.Tool
}

func NewGenerateImageTool(api api.GaiaApi) *GenerateImageTool {
	return &GenerateImageTool{
		api: api,
		tool: mcp.NewTool(
			"generate_image",
			mcp.WithDescription("Generate images with Protogaia"),
			mcp.WithString(
				"prompt",
				mcp.Required(),
				mcp.Description("The prompt to generate an image with"),
			),
			mcp.WithString(
				"aspectRatio",
				mcp.Description("Aspect ratio of the image. One of the following: '1:1', '3:2', '2:3', '16:9', '9:16'"),
				mcp.DefaultString(string(shared.AspectRatio1_1)),
				mcp.Enum(shared.GetAspectRatioMap().ToStrings()...),
			),
			mcp.WithString(
				"promptStyle",
				mcp.Description("Style to apply to the generated image. Choose from predefined styles. It's not style id and style name."),
				mcp.DefaultString(string(shared.PromptStyleBase)),
				mcp.Enum(shared.GetPromptStyleMap().ToStrings()...),
			),
			mcp.WithString(
				"styleId",
				mcp.Description("The style ID to use. It must be styleId created by create_style_tool from Gaia"),
			),
		),
	}
}

func (t *GenerateImageTool) ToolName() string {
	return "generate_image"
}

func (t *GenerateImageTool) MCPTool() mcp.Tool {
	return t.tool
}

func (t *GenerateImageTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	// Get the arguments from tool call request
	prompt := args["prompt"]
	aspectRatio := args["aspectRatio"]
	promptStyle := args["promptStyle"]
	styleId := args["styleId"]

	res, err := t.api.GenerateImages(ctx, api.GenerateImagesRequest{
		RecipeId: shared.RecipeIdImageGeneratorSimple,
		Params: map[string]interface{}{
			"prompt":         prompt,
			"aspectRatio":    aspectRatio,
			"promptStyle":    promptStyle,
			"styleId":        styleId,
			"numberOfImages": 1, // Always generate 1 image
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

	// Check if we actually received any images
	if len(res.Images) == 0 {
		return mcp.NewToolResultError("No images were generated. Please try again."), nil
	}

	// Process the image using the imageutil package for MCP
	base64Data, mimeType, err := imageutil.ProcessImageQuickForMCP(ctx, res.Images[0])
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process image: %v", err)), nil
	}

	msg := fmt.Sprintf("Image generated successfully. Image url: %s", res.Images[0])

	return mcp.NewToolResultImage(msg, base64Data, mimeType), nil
}
