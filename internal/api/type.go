package api

import "gaia-mcp-go/pkg/shared"

// UploadFile represents a file upload with all associated metadata
type UploadFile struct {
	// Id is the unique identifier for the uploaded file
	Id string `json:"id"`

	// Name is the original filename
	Name string `json:"name"`

	// Size is the file size in bytes
	Size int64 `json:"size,string"`

	// MimeType represents the file's MIME type (e.g., "image/jpeg", "text/plain")
	MimeType string `json:"mimeType"`

	// Metadata is optional additional data stored as key-value pairs
	// Using pointer to indicate it's optional (nil when not provided)
	Metadata *map[string]interface{} `json:"metadata,omitempty"`

	// Url is the download URL for the file, can be null
	// Using pointer to string to represent nullable field
	Url *string `json:"url"`

	// OwnerUid is the unique identifier of the user who owns this file
	OwnerUid string `json:"ownerUid"`

	// CreatedAt is the timestamp when the file was created (stored as string)
	CreatedAt string `json:"createdAt"`

	// Uploaded indicates whether the file has been successfully uploaded
	Uploaded bool `json:"uploaded"`
}

// UploadPart represents a completed upload part
type UploadPart struct {
	ETag       string `json:"eTag"`
	PartNumber int    `json:"partNumber"`
}

// InitUploadResponse represents the response when initializing a file upload
type InitUploadResponse struct {
	// Key is the unique identifier for this upload session
	Key string `json:"key"`

	// Filename is the name of the file being uploaded
	Filename string `json:"filename"`

	// UploadId is the unique identifier for the multipart upload
	UploadId string `json:"uploadId"`

	// UploadUrls is an array of presigned URLs for uploading file chunks
	UploadUrls []string `json:"uploadUrls"`

	// File contains the complete file metadata
	File UploadFile `json:"file"`
}

// SdStyleCreator represents a user who created a style
type SdStyleCreator struct {
	// Uid is the unique identifier for the user
	Uid string `json:"uid"`

	// Name is the user's display name
	Name string `json:"name"`

	// Email is the user's email address
	Email string `json:"email"`

	// Picture is the URL to the user's profile picture
	Picture string `json:"picture"`

	// Username is the user's unique username
	Username string `json:"username"`
}

// SdStyleTag represents a tag that can be applied to styles
type SdStyleTag struct {
	// Id is the unique identifier for the tag
	Id int `json:"id"`

	// Name is the display name of the tag
	Name string `json:"name"`
}

// SdStyleWorkspace represents a workspace containing styles
type SdStyleWorkspace struct {
	// Id is the unique identifier for the workspace
	Id string `json:"id"`

	// Name is the display name of the workspace
	Name string `json:"name"`

	// Picture is the URL to the workspace's image/icon
	Picture string `json:"picture"`
}

// SdStyleImage represents an image used in a style with its weight
type SdStyleImage struct {
	// Url is the URL to the image
	Url string `json:"url"`

	// Weight determines the influence of this image in the style (typically 0.0 to 1.0)
	Weight float64 `json:"weight"`
}

// SdStyleCapabilities represents what actions a user can perform on a style
type SdStyleCapabilities struct {
	// CanView indicates if the user can view the style
	CanView bool `json:"canView"`

	// CanUpdate indicates if the user can modify the style
	CanUpdate bool `json:"canUpdate"`

	// CanMove indicates if the user can move the style between workspaces
	CanMove bool `json:"canMove"`

	// CanDelete indicates if the user can delete the style
	CanDelete bool `json:"canDelete"`

	// CanRecover indicates if the user can recover a deleted style
	CanRecover bool `json:"canRecover"`

	// CanShare indicates if the user can share the style with others
	CanShare bool `json:"canShare"`

	// CanAddToLibrary indicates if the user can add the style to a library
	CanAddToLibrary bool `json:"canAddToLibrary"`

	// CanRemoveFromLibrary indicates if the user can remove the style from a library
	CanRemoveFromLibrary bool `json:"canRemoveFromLibrary"`
}

// SdStyleMetric represents metrics/statistics for a style
type SdStyleMetric struct {
	// Id is the unique identifier for the metric record
	Id int `json:"id"`

	// FavoriteCount is the number of times this style has been favorited
	FavoriteCount int `json:"favoriteCount"`
}

// ThumbnailModerationRating represents the content moderation rating for thumbnails
type ThumbnailModerationRating string

const (
	ThumbnailModerationUnrated   ThumbnailModerationRating = "unrated"
	ThumbnailModerationSafe      ThumbnailModerationRating = "safe"
	ThumbnailModerationSensitive ThumbnailModerationRating = "sensitive"
	ThumbnailModerationUnsafe    ThumbnailModerationRating = "unsafe"
)

// SharingMode represents how a style can be shared
type SharingMode string

const (
	SharingModeRestricted SharingMode = "restricted"
	SharingModePublic     SharingMode = "public"
	SharingModePrivate    SharingMode = "private"
)

// SdStyle represents a complete AI style definition
type SdStyle struct {
	// Id is the unique identifier for the style
	Id string `json:"id"`

	// Name is the display name of the style
	Name string `json:"name"`

	// ThumbnailUrl is the URL to the style's preview image
	ThumbnailUrl string `json:"thumbnailUrl"`

	// ThumbnailWidth is the width of the thumbnail in pixels
	ThumbnailWidth int `json:"thumbnailWidth"`

	// ThumbnailHeight is the height of the thumbnail in pixels
	ThumbnailHeight int `json:"thumbnailHeight"`

	// ThumbnailModerationRating indicates the content rating of the thumbnail
	// Valid values: "unrated", "safe", "sensitive", "unsafe"
	ThumbnailModerationRating ThumbnailModerationRating `json:"thumbnailModerationRating"`

	// IsDraft indicates if this style is still in draft mode
	IsDraft bool `json:"isDraft"`

	// Description provides detailed information about the style
	Description string `json:"description"`

	// DiscoverableAt is the timestamp when the style became publicly discoverable
	// Can be null if the style is not discoverable
	DiscoverableAt *string `json:"discoverableAt"`

	// DeletedAt is the timestamp when the style was deleted
	// Can be null if the style is not deleted
	DeletedAt *string `json:"deletedAt"`

	// SharingMode determines how the style can be shared
	// Valid values: "restricted", "public", "private"
	SharingMode SharingMode `json:"sharingMode"`

	// Creator contains information about who created this style
	Creator SdStyleCreator `json:"creator"`

	// Tags is an array of tags associated with this style
	Tags []SdStyleTag `json:"tags"`

	// Workspace contains information about the workspace this style belongs to
	Workspace SdStyleWorkspace `json:"workspace"`

	// WorkspaceId is the unique identifier of the workspace
	WorkspaceId string `json:"workspaceId"`

	// Images is an array of reference images used to create this style
	Images []SdStyleImage `json:"images"`

	// Pinned indicates if this style is pinned for the user
	Pinned bool `json:"pinned"`

	// Capabilities defines what actions the current user can perform on this style
	Capabilities SdStyleCapabilities `json:"capabilities"`

	// FavoritedByUser indicates if the current user has favorited this style
	FavoritedByUser bool `json:"favoritedByUser"`

	// Metric contains usage statistics and metrics for this style
	Metric SdStyleMetric `json:"metric"`

	// CreatedAt is the timestamp when the style was created
	CreatedAt string `json:"createdAt"`
}

// RecipeTaskRequest represents a request to execute a recipe with parameters
type RecipeTaskRequest struct {
	// RecipeId is the unique identifier of the recipe to execute
	RecipeId string `json:"recipeId"`

	// Params is a map of parameter names to their values
	// Values can be of any type (string, number, boolean, object, etc.)
	Params map[string]interface{} `json:"params"`
}

// Image represents an AI-generated image with complete metadata
type Image struct {
	// CreatedAt is the timestamp when the image was created
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp when the image was last updated
	UpdatedAt string `json:"updatedAt"`

	// Id is the unique identifier for the image
	Id string `json:"id"`

	// S3Key is the key used to store the image in S3
	S3Key string `json:"s3Key"`

	// Name is the display name of the image
	Name string `json:"name"`

	// Sampler is the sampling method used for generation (e.g., "DPM++ 2M Karras")
	Sampler string `json:"sampler"`

	// Steps is the number of inference steps used in generation
	Steps int `json:"steps"`

	// Width is the image width in pixels
	Width int `json:"width"`

	// Height is the image height in pixels
	Height int `json:"height"`

	// CfgScale is the classifier-free guidance scale value
	CfgScale string `json:"cfgScale"`

	// Seed is the random seed used for generation
	Seed string `json:"seed"`

	// ModelHash is the hash identifier of the AI model used
	ModelHash string `json:"modelHash"`

	// ModelName is the name of the AI model used for generation
	ModelName string `json:"modelName"`

	// Prompt is the text prompt used to generate the image
	Prompt string `json:"prompt"`

	// NegativePrompt is the negative prompt used to avoid certain features
	NegativePrompt string `json:"negativePrompt"`

	// DenoisingStrength controls the amount of denoising applied
	DenoisingStrength string `json:"denoisingStrength"`

	// Upscaler is the upscaling method used (if any)
	Upscaler string `json:"upscaler"`

	// ScaleFactor is the scaling factor applied during upscaling
	ScaleFactor string `json:"scaleFactor"`

	// BlurHash is a compact representation of the image for placeholders
	BlurHash string `json:"blurHash"`

	// Note is a user-provided note or description for the image
	Note string `json:"note"`

	// Tags is optional metadata tags associated with the image
	// Can be null if no tags are assigned
	Tags *map[string]interface{} `json:"tags,omitempty"`

	// Owner contains optional information about the image owner
	// Can be null if owner information is not available
	Owner *map[string]interface{} `json:"owner,omitempty"`

	// OwnerUid is the unique identifier of the image owner
	OwnerUid string `json:"ownerUid"`

	// Folder contains optional information about the containing folder
	// Can be null if the image is not in a specific folder
	Folder *map[string]interface{} `json:"folder,omitempty"`

	// FolderId is the unique identifier of the containing folder
	FolderId string `json:"folderId"`

	// DeletedAt is the timestamp when the image was deleted
	// Can be null if the image is not deleted
	DeletedAt *string `json:"deletedAt,omitempty"`

	// ExpireAt is the timestamp when the image will expire
	// Can be null if the image doesn't have an expiration
	ExpireAt *string `json:"expireAt,omitempty"`

	// OriginalFolder contains optional information about the original folder
	// Can be null if not applicable
	OriginalFolder *map[string]interface{} `json:"originalFolder,omitempty"`

	// RecipeTask contains optional information about the generation task
	// Can be null if not generated from a recipe task
	RecipeTask *map[string]interface{} `json:"recipeTask,omitempty"`

	// RecipeTaskId is the unique identifier of the generation task
	RecipeTaskId string `json:"recipeTaskId"`

	// Recipe contains optional information about the generation recipe
	// Can be null if not generated from a recipe
	Recipe *map[string]interface{} `json:"recipe,omitempty"`

	// RecipeId is the unique identifier of the generation recipe
	RecipeId string `json:"recipeId"`

	// ImagePermissions contains optional permission settings for the image
	// Can be null if using default permissions
	ImagePermissions *map[string]interface{} `json:"imagePermissions,omitempty"`

	// Size is the file size of the image in bytes
	Size int `json:"size"`

	// FullMetadata contains the complete generation metadata as a JSON string
	FullMetadata string `json:"fullMetadata"`

	// Url is the public URL to access the image
	Url string `json:"url"`
}

// RecipeType represents the type of recipe being executed
type RecipeType string

const (
	// TODO: Replace these with actual values from gaiaRecipeType
	RecipeTypeNormal RecipeType = "normal"
	// Add other recipe types here: RecipeTypeAdvanced, RecipeTypeCustom, etc.
)

// RecipeTaskStatus represents the current status of a recipe task
type RecipeTaskStatus string

const (
	// TODO: Replace these with actual values from gaiaRecipeTaskStatus
	RecipeTaskStatusPending   RecipeTaskStatus = "pending"
	RecipeTaskStatusRunning   RecipeTaskStatus = "running"
	RecipeTaskStatusCompleted RecipeTaskStatus = "completed"
	RecipeTaskStatusFailed    RecipeTaskStatus = "failed"
	// Add other statuses here: RecipeTaskStatusCancelled, etc.
)

// QueueType represents the type of processing queue
type QueueType string

const (
	// TODO: Replace these with actual values from gaiaQueueType
	QueueTypeDefault  QueueType = "default"
	QueueTypePriority QueueType = "priority"
	// Add other queue types here: QueueTypeBatch, QueueTypeExpress, etc.
)

// RecipeTaskCreator represents the creator of a recipe task
type RecipeTaskCreator struct {
	// Uid is the unique identifier of the user who created the task
	Uid string `json:"uid"`
}

// RecipeTask represents a recipe execution task with complete metadata
type RecipeTask struct {
	// CreatedAt is the timestamp when the task was created
	CreatedAt string `json:"createdAt"`

	// UpdatedAt is the timestamp when the task was last updated
	UpdatedAt string `json:"updatedAt"`

	// Id is the unique identifier for the task
	Id string `json:"id"`

	// RecipeId is the unique identifier of the recipe being executed
	RecipeId string `json:"recipeId"`

	// RecipeType indicates the type of recipe (defaults to "normal")
	RecipeType RecipeType `json:"recipeType"`

	// Params contains optional parameters for the recipe execution
	// Can be null if no custom parameters are provided
	Params *map[string]interface{} `json:"params,omitempty"`

	// FolderId is the unique identifier of the folder containing results
	FolderId string `json:"folderId"`

	// Creator contains information about who created this task
	Creator RecipeTaskCreator `json:"creator"`

	// Status indicates the current execution status of the task
	Status RecipeTaskStatus `json:"status"`

	// Priority determines the execution priority (higher numbers = higher priority)
	Priority int `json:"priority"`

	// StartedAt is the timestamp when task execution began
	// Can be null if the task hasn't started yet
	StartedAt *string `json:"startedAt,omitempty"`

	// CompletedAt is the timestamp when task execution completed
	// Can be null if the task hasn't completed yet
	CompletedAt *string `json:"completedAt,omitempty"`

	// IsDeleted indicates whether this task has been deleted
	// Can be null (defaults to false)
	IsDeleted *bool `json:"isDeleted,omitempty"`

	// DeletedAt is the timestamp when the task was deleted
	// Can be null if the task is not deleted
	DeletedAt *string `json:"deletedAt,omitempty"`

	// Images contains all images generated by this task
	Images []Image `json:"images"`

	// Name is the display name of the task
	Name string `json:"name"`

	// Prompt is the text prompt used for generation
	Prompt string `json:"prompt"`

	// Seed is the random seed used for generation reproducibility
	Seed int `json:"seed"`

	// RunnerId identifies which runner/worker executed this task
	RunnerId string `json:"runnerId"`

	// Error contains error information if the task failed
	// Can be null if no error occurred
	Error *string `json:"error,omitempty"`

	// ResultImages contains URLs or identifiers of the generated images
	ResultImages []string `json:"resultImages"`

	// ExecutionDuration is the time taken to complete the task in milliseconds
	// Can be null if the task hasn't completed or duration wasn't measured
	ExecutionDuration *int `json:"executionDuration,omitempty"`

	// QueueType indicates which processing queue this task uses (defaults to "default")
	QueueType QueueType `json:"queueType"`
}

type ImageGeneratedResponse struct {
	Success bool     `json:"success"`
	Images  []string `json:"images"`
	Error   *string  `json:"error,omitempty"`
}

type GenerateImagesRequest struct {
	RecipeId shared.RecipeId        `json:"recipeId"`
	Params   map[string]interface{} `json:"params"`
}
