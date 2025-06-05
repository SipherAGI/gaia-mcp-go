package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRecipeTaskStatus tests the RecipeTaskStatus constants
func TestRecipeTaskStatus(t *testing.T) {
	t.Run("Verify all status constants", func(t *testing.T) {
		// Test that constants have expected values
		assert.Equal(t, RecipeTaskStatus("QUEUED"), RecipeTaskStatusQueued)
		assert.Equal(t, RecipeTaskStatus("RUNNING"), RecipeTaskStatusRunning)
		assert.Equal(t, RecipeTaskStatus("COMPLETED"), RecipeTaskStatusCompleted)
		assert.Equal(t, RecipeTaskStatus("FAILED"), RecipeTaskStatusFailed)
		assert.Equal(t, RecipeTaskStatus("CANCELLED"), RecipeTaskStatusCancelled)
		assert.Equal(t, RecipeTaskStatus("CANCELED"), RecipeTaskStatusCanceled)
		assert.Equal(t, RecipeTaskStatus("DRAFT"), RecipeTaskStatusDraft)
	})
}

// TestRecipeType tests the RecipeType constants
func TestRecipeType(t *testing.T) {
	t.Run("Verify all recipe type constants", func(t *testing.T) {
		assert.Equal(t, RecipeType("normal"), RecipeTypeNormal)
		assert.Equal(t, RecipeType("inpaint"), RecipeTypeInpaint)
		assert.Equal(t, RecipeType("chain"), RecipeTypeChain)
		assert.Equal(t, RecipeType("comfyui"), RecipeTypeComfyui)
		assert.Equal(t, RecipeType("describe"), RecipeTypeDescribe)
		assert.Equal(t, RecipeType("turbo"), RecipeTypeTurbo)
		assert.Equal(t, RecipeType("other"), RecipeTypeOther)
	})
}

// TestPromptStyle tests the PromptStyle constants
func TestPromptStyle(t *testing.T) {
	t.Run("Verify critical prompt styles", func(t *testing.T) {
		assert.Equal(t, PromptStyle("base"), PromptStyleBase)
		assert.Equal(t, PromptStyle("enhance"), PromptStyleEnhance)
		assert.Equal(t, PromptStyle("anime"), PromptStyleAnime)
		assert.Equal(t, PromptStyle("photographic"), PromptStylePhotographic)
		assert.Equal(t, PromptStyle("cinematic"), PromptStyleCinematic)
	})
}

// TestAspectRatio tests the AspectRatio constants
func TestAspectRatio(t *testing.T) {
	t.Run("Verify aspect ratio constants", func(t *testing.T) {
		assert.Equal(t, AspectRatio("1:1"), AspectRatio1_1)
		assert.Equal(t, AspectRatio("3:2"), AspectRatio3_2)
		assert.Equal(t, AspectRatio("2:3"), AspectRatio2_3)
		assert.Equal(t, AspectRatio("16:9"), AspectRatio16_9)
		assert.Equal(t, AspectRatio("9:16"), AspectRatio9_16)
	})
}

// TestRecipeId tests the RecipeId constants
func TestRecipeId(t *testing.T) {
	t.Run("Verify recipe ID constants", func(t *testing.T) {
		assert.Equal(t, RecipeId("image-generator-simple"), RecipeIdImageGeneratorSimple)
		assert.Equal(t, RecipeId("face-enhancer"), RecipeIdFaceEnhancer)
		assert.Equal(t, RecipeId("remix"), RecipeIdRemix)
		assert.Equal(t, RecipeId("upscaler"), RecipeIdUpscaler)
	})
}

// TestFileAssociatedResource tests the FileAssociatedResource constants
func TestFileAssociatedResource(t *testing.T) {
	t.Run("Verify file resource constants", func(t *testing.T) {
		assert.Equal(t, FileAssociatedResource("STYLE"), FileAssociatedResourceStyle)
		assert.Equal(t, FileAssociatedResource("NONE"), FileAssociatedResourceNone)
		assert.Equal(t, FileAssociatedResource("USER_AVATAR"), FileAssociatedResourceUserAvatar)
	})
}

// TestPromptStyleMap tests the PromptStyleMap functionality
func TestPromptStyleMap(t *testing.T) {
	t.Run("Test GetPromptStyleMap", func(t *testing.T) {
		styleMap := GetPromptStyleMap()
		require.NotNil(t, styleMap, "PromptStyleMap should not be nil")

		// Test specific mappings
		assert.Equal(t, "base", styleMap.Get(PromptStyleBase))
		assert.Equal(t, "enhance", styleMap.Get(PromptStyleEnhance))
		assert.Equal(t, "anime", styleMap.Get(PromptStyleAnime))
		assert.Equal(t, "photographic", styleMap.Get(PromptStylePhotographic))
	})

	t.Run("Test ToStrings method", func(t *testing.T) {
		styleMap := GetPromptStyleMap()
		strings := styleMap.ToStrings()

		// Check that we got some strings back
		assert.Greater(t, len(strings), 0, "ToStrings should return non-empty slice")

		// Check that some expected values are present
		assert.Contains(t, strings, "base")
		assert.Contains(t, strings, "enhance")
		assert.Contains(t, strings, "anime")
	})
}

// TestRecipeTaskStatusMap tests the RecipeTaskStatusMap functionality
func TestRecipeTaskStatusMap(t *testing.T) {
	t.Run("Test GetRecipeTaskStatusMap", func(t *testing.T) {
		statusMap := GetRecipeTaskStatusMap()
		require.NotNil(t, statusMap, "RecipeTaskStatusMap should not be nil")

		// Test specific mappings
		assert.Equal(t, "QUEUED", statusMap.Get(RecipeTaskStatusQueued))
		assert.Equal(t, "RUNNING", statusMap.Get(RecipeTaskStatusRunning))
		assert.Equal(t, "COMPLETED", statusMap.Get(RecipeTaskStatusCompleted))
		assert.Equal(t, "FAILED", statusMap.Get(RecipeTaskStatusFailed))
	})

	t.Run("Test ToStrings method", func(t *testing.T) {
		statusMap := GetRecipeTaskStatusMap()
		strings := statusMap.ToStrings()

		assert.Greater(t, len(strings), 0, "ToStrings should return non-empty slice")
		assert.Contains(t, strings, "QUEUED")
		assert.Contains(t, strings, "RUNNING")
		assert.Contains(t, strings, "COMPLETED")
	})
}

// TestRecipeTypeMap tests the RecipeTypeMap functionality
func TestRecipeTypeMap(t *testing.T) {
	t.Run("Test GetRecipeTypeMap", func(t *testing.T) {
		typeMap := GetRecipeTypeMap()
		require.NotNil(t, typeMap, "RecipeTypeMap should not be nil")

		// Test specific mappings
		assert.Equal(t, "normal", typeMap.Get(RecipeTypeNormal))
		assert.Equal(t, "inpaint", typeMap.Get(RecipeTypeInpaint))
		assert.Equal(t, "chain", typeMap.Get(RecipeTypeChain))
	})

	t.Run("Test ToStrings method", func(t *testing.T) {
		typeMap := GetRecipeTypeMap()
		strings := typeMap.ToStrings()

		assert.Greater(t, len(strings), 0, "ToStrings should return non-empty slice")
		assert.Contains(t, strings, "normal")
		assert.Contains(t, strings, "inpaint")
	})
}

// TestAspectRatioMap tests the AspectRatioMap functionality
func TestAspectRatioMap(t *testing.T) {
	t.Run("Test GetAspectRatioMap", func(t *testing.T) {
		ratioMap := GetAspectRatioMap()
		require.NotNil(t, ratioMap, "AspectRatioMap should not be nil")

		// Test specific mappings
		assert.Equal(t, "1:1", ratioMap.Get(AspectRatio1_1))
		assert.Equal(t, "16:9", ratioMap.Get(AspectRatio16_9))
		assert.Equal(t, "9:16", ratioMap.Get(AspectRatio9_16))
	})

	t.Run("Test ToStrings method", func(t *testing.T) {
		ratioMap := GetAspectRatioMap()
		strings := ratioMap.ToStrings()

		assert.Greater(t, len(strings), 0, "ToStrings should return non-empty slice")
		assert.Contains(t, strings, "1:1")
		assert.Contains(t, strings, "16:9")
	})
}

// TestQueueTypeMap tests the QueueTypeMap functionality
func TestQueueTypeMap(t *testing.T) {
	t.Run("Test GetQueueTypeMap", func(t *testing.T) {
		queueMap := GetQueueTypeMap()
		require.NotNil(t, queueMap, "QueueTypeMap should not be nil")

		// Test specific mappings - using actual constants from the code
		assert.Equal(t, "default", queueMap.Get(QueueTypeDefault))
		assert.Equal(t, "fast", queueMap.Get(QueueTypeFast))
		assert.Equal(t, "flux1", queueMap.Get(QueueTypeFlux1))
		assert.Equal(t, "dedicated", queueMap.Get(QueueTypeDedicated))
	})
}

// TestFileAssociatedResourceMap tests the FileAssociatedResourceMap functionality
func TestFileAssociatedResourceMap(t *testing.T) {
	t.Run("Test GetFileAssociatedResourceMap", func(t *testing.T) {
		resourceMap := GetFileAssociatedResourceMap()
		require.NotNil(t, resourceMap, "FileAssociatedResourceMap should not be nil")

		// Test specific mappings using actual constants
		assert.Equal(t, "STYLE", resourceMap.Get(FileAssociatedResourceStyle))
		assert.Equal(t, "NONE", resourceMap.Get(FileAssociatedResourceNone))
		assert.Equal(t, "USER_AVATAR", resourceMap.Get(FileAssociatedResourceUserAvatar))
	})
}

// TestRecipeIdMap tests the RecipeIdMap functionality
func TestRecipeIdMap(t *testing.T) {
	t.Run("Test GetRecipeIdMap", func(t *testing.T) {
		recipeMap := GetRecipeIdMap()
		require.NotNil(t, recipeMap, "RecipeIdMap should not be nil")

		// Test specific mappings
		assert.Equal(t, "image-generator-simple", recipeMap.Get(RecipeIdImageGeneratorSimple))
		assert.Equal(t, "face-enhancer", recipeMap.Get(RecipeIdFaceEnhancer))
		assert.Equal(t, "remix", recipeMap.Get(RecipeIdRemix))
		assert.Equal(t, "upscaler", recipeMap.Get(RecipeIdUpscaler))
	})
}

// TestConstants tests the exported constants
func TestConstants(t *testing.T) {
	t.Run("Test exported constants", func(t *testing.T) {
		assert.Equal(t, "https://protogaia.com", HOMEPAGE_URL)
		assert.Equal(t, "https://api.protogaia.com", BASE_API_URL)
		assert.Equal(t, 1024*1024*10, UPLOAD_CHUNK_SIZE, "Upload chunk size should be 10MB")
	})
}

// Example of testing edge cases and error conditions
func TestMapEdgeCases(t *testing.T) {
	t.Run("Test with invalid enum values", func(t *testing.T) {
		styleMap := GetPromptStyleMap()

		// Test with invalid PromptStyle (this should return empty string)
		invalidStyle := PromptStyle("invalid-style")
		result := styleMap.Get(invalidStyle)
		assert.Equal(t, "", result, "Invalid enum should return empty string")
	})
}

// Benchmark tests to ensure performance is acceptable
func BenchmarkPromptStyleMapGet(b *testing.B) {
	styleMap := GetPromptStyleMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = styleMap.Get(PromptStyleBase)
	}
}

func BenchmarkPromptStyleMapToStrings(b *testing.B) {
	styleMap := GetPromptStyleMap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = styleMap.ToStrings()
	}
}
