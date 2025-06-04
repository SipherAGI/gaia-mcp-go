package tools

import (
	"context"
	"fmt"
	"gaia-mcp-go/internal/api"
	"gaia-mcp-go/pkg/shared"

	"github.com/mark3labs/mcp-go/mcp"
)

type UploadImageTool struct {
	api  api.GaiaApi
	tool mcp.Tool
}

func NewUploadImageTool(api api.GaiaApi) *UploadImageTool {
	return &UploadImageTool{
		api: api,
		tool: mcp.NewTool(
			"upload_image",
			mcp.WithDescription("Upload an image to GAIA"),
			mcp.WithArray(
				"image_urls",
				mcp.Items(map[string]any{"type": "string"}),
				mcp.Required(),
				mcp.Description("The URLs of the images to upload"),
			),
		),
	}
}

func (t *UploadImageTool) ToolName() string {
	return "upload_image"
}

func (t *UploadImageTool) MCPTool() mcp.Tool {
	return t.tool
}

func (t *UploadImageTool) Handler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	// First, get the raw value and check if it exists
	rawImageUrls, exists := args["image_urls"]
	if !exists {
		return mcp.NewToolResultError("image_urls parameter is required"), nil
	}

	// Convert from []interface{} to []string safely
	imageUrlsInterface, ok := rawImageUrls.([]interface{})
	if !ok {
		return mcp.NewToolResultError("image_urls must be an array"), nil
	}

	// Convert each interface{} element to string
	imageUrls := make([]string, len(imageUrlsInterface))
	for i, url := range imageUrlsInterface {
		urlStr, ok := url.(string)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("image_urls[%d] must be a string", i)), nil
		}
		imageUrls[i] = urlStr
	}

	// Now we can safely use imageUrls as []string
	uploadedFiles, err := t.api.UploadImages(ctx, imageUrls, shared.FileAssociatedResourceStyle)

	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	resultMsg := fmt.Sprintf("Uploaded %d images successfully", len(uploadedFiles))
	resultMsg += "\nFile urls:\n"
	for _, file := range uploadedFiles {
		if file.Url != nil {
			resultMsg += fmt.Sprintf("- %s\n", *file.Url)
		}
	}

	return mcp.NewToolResultText(resultMsg), nil
}
