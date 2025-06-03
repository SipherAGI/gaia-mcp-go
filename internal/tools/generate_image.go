package tools

import (
	"context"
	"gaia-mcp-go/internal/api"
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
				mcp.DefaultString("1:1"),
				mcp.Enum(
					"1:1",
					"3:2",
					"2:3",
					"16:9",
					"9:16",
				),
			),
			mcp.WithString(
				"promptStyle",
				mcp.Description("Style to apply to the generated image. Choose from predefined styles. It's not style id and style name."),
				mcp.DefaultString("base"),
				mcp.Enum(
					"base",
					"enhance",
					"anime",
					"photographic",
					"cinematic",
					"analog film",
					"digital art",
					"fantasy art",
					"line art",
					"pixel art",
					"artstyle-watercolor",
					"comic book",
					"neonpunk",
					"3d-model",
					"misc-fairy tale",
					"misc-gothic",
					"photo-long exposure",
					"photo-tilt-shift",
					"lowpoly",
					"origami",
					"craft clay",
					"game-minecraft",
				),
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

	// TODO: The response is an url, but we need to return a base64 encoded image.
	// We need to download the image from the url and then encode it to base64 and resize it to maximum 1024x1024

	return mcp.NewToolResultText(res.Images[0]), nil
}
