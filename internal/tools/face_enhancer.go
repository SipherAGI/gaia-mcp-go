package tools

import (
	"context"
	"fmt"
	"gaia-mcp-go/internal/api"
	"gaia-mcp-go/pkg/imageutil"
	"gaia-mcp-go/pkg/shared"

	"github.com/mark3labs/mcp-go/mcp"
)

type FaceEnhancerTool struct {
	api            api.GaiaApi
	tool           mcp.Tool
	imageProcessor *imageutil.Processor
}

func NewFaceEnhancerTool(
	api api.GaiaApi,
	imageProcessor *imageutil.Processor,
) *FaceEnhancerTool {
	return &FaceEnhancerTool{
		api:            api,
		imageProcessor: imageProcessor,
		tool: mcp.NewTool(
			"face_enhancer",
			mcp.WithDescription("Enhance face's details in an existing image"),
			mcp.WithString(
				"image_url",
				mcp.Required(),
				mcp.Description("The image URL to enhance. It must be GAIA's image url: starts with `https://cdn.protogaia.com/`"),
			),
			mcp.WithString(
				"prompt",
				mcp.Description("The prompt to tell AI what to enhance."),
			),
		),
	}
}

func (t *FaceEnhancerTool) ToolName() string {
	return "face_enhancer"
}

func (t *FaceEnhancerTool) MCPTool() mcp.Tool {
	return t.tool
}

func (t *FaceEnhancerTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	imageUrl := args["image_url"]
	prompt := args["prompt"]

	res, err := t.api.GenerateImages(ctx, api.GenerateImagesRequest{
		RecipeId: shared.RecipeIdFaceEnhancer,
		Params: map[string]interface{}{
			"imageUrl": imageUrl,
			"prompt":   prompt,
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

	base64Data, mimeType, err := t.imageProcessor.ProcessImageFromURLForMCP(ctx, res.Images[0])
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to process image: %v", err)), nil
	}

	msg := fmt.Sprintf("Face enhanced successfully. Image url: %s", res.Images[0])

	return mcp.NewToolResultImage(msg, base64Data, mimeType), nil
}
