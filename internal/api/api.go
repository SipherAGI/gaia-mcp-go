package api

import (
	"context"
	"fmt"
	"gaia-mcp-go/pkg/httpclient"
)

type GaiaApi interface {
	CreateStyle(ctx context.Context, imageUrls []string, name string, description *string) (SdStyle, error)
	GenerateImages(ctx context.Context, req GenerateImagesRequest) (ImageGeneratedResponse, error)
}

type GaiaApiConfig struct {
	BaseUrl string
	ApiKey  string
}

type gaiaApi struct {
	client *httpclient.Client
}

func NewGaiaApi(cfg GaiaApiConfig) GaiaApi {
	client := httpclient.New(httpclient.Config{
		BaseURL: cfg.BaseUrl,
		DefaultHeaders: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cfg.ApiKey),
		},
	})
	return &gaiaApi{client: client}
}

func (a *gaiaApi) CreateStyle(ctx context.Context, imageUrls []string, name string, description *string) (SdStyle, error) {
	// Formatting imageUrls to be an array of images
	images := make([]map[string]interface{}, len(imageUrls))
	for i, imageUrl := range imageUrls {
		images[i] = map[string]interface{}{
			"url":    imageUrl,
			"weight": 0.5,
		}
	}

	payload := map[string]interface{}{
		"images": images,
		"name":   name,
	}

	if description != nil {
		payload["description"] = *description
	}

	rb := a.client.PostJSON(ctx, "/api/sd-styles", payload, map[string]string{})

	var sdStyle SdStyle
	err := rb.Into(&sdStyle)
	if err != nil {
		return SdStyle{}, err
	}

	return sdStyle, nil
}

func (a *gaiaApi) GenerateImages(ctx context.Context, req GenerateImagesRequest) (ImageGeneratedResponse, error) {
	rb := a.client.PostJSON(ctx, "/api/recipe/agi-tasks/create-task", req, map[string]string{})

	var imageGeneratedResponse ImageGeneratedResponse
	err := rb.Into(&imageGeneratedResponse)
	if err != nil {
		return ImageGeneratedResponse{}, err
	}

	return imageGeneratedResponse, nil
}
