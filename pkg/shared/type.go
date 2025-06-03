package shared

type RecipeTaskStatus string
type RecipeType string
type PromptStyle string
type AspectRatio string
type QueueType string
type RecipeId string
type FileAssociatedResource string

const (
	// RecipeTaskStatus
	RecipeTaskStatusQueued    RecipeTaskStatus = "QUEUED"
	RecipeTaskStatusRunning   RecipeTaskStatus = "RUNNING"
	RecipeTaskStatusCompleted RecipeTaskStatus = "COMPLETED"
	RecipeTaskStatusFailed    RecipeTaskStatus = "FAILED"
	RecipeTaskStatusCancelled RecipeTaskStatus = "CANCELLED"
	RecipeTaskStatusCanceled  RecipeTaskStatus = "CANCELED"
	RecipeTaskStatusDraft     RecipeTaskStatus = "DRAFT"

	// RecipeType
	RecipeTypeNormal   RecipeType = "normal"
	RecipeTypeInpaint  RecipeType = "inpaint"
	RecipeTypeChain    RecipeType = "chain"
	RecipeTypeComfyui  RecipeType = "comfyui"
	RecipeTypeDescribe RecipeType = "describe"
	RecipeTypeTurbo    RecipeType = "turbo"
	RecipeTypeOther    RecipeType = "other"

	// PromptStyle
	PromptStyleBase               PromptStyle = "base"
	PromptStyleEnhance            PromptStyle = "enhance"
	PromptStyleAnime              PromptStyle = "anime"
	PromptStylePhotographic       PromptStyle = "photographic"
	PromptStyleCinematic          PromptStyle = "cinematic"
	PromptStyleAnalogFilm         PromptStyle = "analog film"
	PromptStyleDigitalArt         PromptStyle = "digital art"
	PromptStyleFantasyArt         PromptStyle = "fantasy art"
	PromptStyleLineArt            PromptStyle = "line art"
	PromptStylePixelArt           PromptStyle = "pixel art"
	PromptStyleArtstyleWatercolor PromptStyle = "artstyle-watercolor"
	PromptStyleComicBook          PromptStyle = "comic book"
	PromptStyleNeonpunk           PromptStyle = "neonpunk"
	PromptStyle3DModel            PromptStyle = "3d-model"
	PromptStyleMiscFairyTale      PromptStyle = "misc-fairy tale"
	PromptStyleMiscGothic         PromptStyle = "misc-gothic"
	PromptStylePhotoLongExposure  PromptStyle = "photo-long exposure"
	PromptStylePhotoTiltShift     PromptStyle = "photo-tilt-shift"
	PromptStyleLowpoly            PromptStyle = "lowpoly"
	PromptStyleOrigami            PromptStyle = "origami"
	PromptStyleCraftClay          PromptStyle = "craft clay"
	PromptStyleGameMinecraft      PromptStyle = "game-minecraft"

	// AspectRatio
	AspectRatio1_1  AspectRatio = "1:1"
	AspectRatio3_2  AspectRatio = "3:2"
	AspectRatio2_3  AspectRatio = "2:3"
	AspectRatio16_9 AspectRatio = "16:9"
	AspectRatio9_16 AspectRatio = "9:16"

	// QueueType
	QueueTypeDefault   QueueType = "default"
	QueueTypeFast      QueueType = "fast"
	QueueTypeFlux1     QueueType = "flux1"
	QueueTypeDedicated QueueType = "dedicated"
	QueueTypeOther     QueueType = "other"

	// FileAssociatedResource
	FileAssociatedResourceUserAvatar        FileAssociatedResource = "USER_AVATAR"
	FileAssociatedResourceUserCoverImage    FileAssociatedResource = "USER_COVER_IMAGE"
	FileAssociatedResourceWorkspace         FileAssociatedResource = "WORKSPACE"
	FileAssociatedResourceArticleCoverImage FileAssociatedResource = "ARTICLE_COVER_IMAGE"
	FileAssociatedResourceArticleFile       FileAssociatedResource = "ARTICLE_FILE"
	FileAssociatedResourceStyle             FileAssociatedResource = "STYLE"
	FileAssociatedResourceSDWorkflow        FileAssociatedResource = "SD_WORKFLOW"
	FileAssociatedResourceChatRoomThumbnail FileAssociatedResource = "CHAT_ROOM_THUMBNAIL"
	FileAssociatedResourceSDModel           FileAssociatedResource = "SD_MODEL"
	FileAssociatedResourceSDModelTraining   FileAssociatedResource = "SD_MODEL_TRAINING"
	FileAssociatedResourcePromptLibrary     FileAssociatedResource = "PROMPT_LIBRARY"
	FileAssociatedResourceNone              FileAssociatedResource = "NONE"

	// RecipeId
	RecipeIdImageGeneratorSimple RecipeId = "image-generator-simple"
	RecipeIdRemix                RecipeId = "remix"
	RecipeIdFaceEnhancer         RecipeId = "face-enhancer"
	RecipeIdUpscaler             RecipeId = "upscaler"
)

type PromptStyleMap struct {
	promptStyles map[PromptStyle]string
}

// RecipeTaskStatusMap provides a mapping for RecipeTaskStatus types
type RecipeTaskStatusMap struct {
	taskStatuses map[RecipeTaskStatus]string
}

// RecipeTypeMap provides a mapping for RecipeType types
type RecipeTypeMap struct {
	recipeTypes map[RecipeType]string
}

// AspectRatioMap provides a mapping for AspectRatio types
type AspectRatioMap struct {
	aspectRatios map[AspectRatio]string
}

// QueueTypeMap provides a mapping for QueueType types
type QueueTypeMap struct {
	queueTypes map[QueueType]string
}

// FileAssociatedResourceMap provides a mapping for FileAssociatedResource types
type FileAssociatedResourceMap struct {
	resources map[FileAssociatedResource]string
}

// RecipeIdMap provides a mapping for RecipeId types
type RecipeIdMap struct {
	recipeIds map[RecipeId]string
}

func GetPromptStyleMap() *PromptStyleMap {
	return &PromptStyleMap{
		promptStyles: map[PromptStyle]string{
			PromptStyleBase:               "base",
			PromptStyleEnhance:            "enhance",
			PromptStyleAnime:              "anime",
			PromptStylePhotographic:       "photographic",
			PromptStyleCinematic:          "cinematic",
			PromptStyleAnalogFilm:         "analog film",
			PromptStyleDigitalArt:         "digital art",
			PromptStyleFantasyArt:         "fantasy art",
			PromptStyleLineArt:            "line art",
			PromptStylePixelArt:           "pixel art",
			PromptStyleArtstyleWatercolor: "artstyle-watercolor",
			PromptStyleComicBook:          "comic book",
			PromptStyleNeonpunk:           "neonpunk",
			PromptStyle3DModel:            "3d-model",
			PromptStyleMiscFairyTale:      "misc-fairy tale",
			PromptStyleMiscGothic:         "misc-gothic",
			PromptStylePhotoLongExposure:  "photo-long exposure",
			PromptStylePhotoTiltShift:     "photo-tilt-shift",
			PromptStyleLowpoly:            "lowpoly",
			PromptStyleOrigami:            "origami",
			PromptStyleCraftClay:          "craft clay",
			PromptStyleGameMinecraft:      "game-minecraft",
		},
	}
}

// GetRecipeTaskStatusMap creates and returns a new RecipeTaskStatusMap
func GetRecipeTaskStatusMap() *RecipeTaskStatusMap {
	return &RecipeTaskStatusMap{
		taskStatuses: map[RecipeTaskStatus]string{
			RecipeTaskStatusQueued:    "QUEUED",
			RecipeTaskStatusRunning:   "RUNNING",
			RecipeTaskStatusCompleted: "COMPLETED",
			RecipeTaskStatusFailed:    "FAILED",
			RecipeTaskStatusCancelled: "CANCELLED",
			RecipeTaskStatusCanceled:  "CANCELED",
			RecipeTaskStatusDraft:     "DRAFT",
		},
	}
}

// GetRecipeTypeMap creates and returns a new RecipeTypeMap
func GetRecipeTypeMap() *RecipeTypeMap {
	return &RecipeTypeMap{
		recipeTypes: map[RecipeType]string{
			RecipeTypeNormal:   "normal",
			RecipeTypeInpaint:  "inpaint",
			RecipeTypeChain:    "chain",
			RecipeTypeComfyui:  "comfyui",
			RecipeTypeDescribe: "describe",
			RecipeTypeTurbo:    "turbo",
			RecipeTypeOther:    "other",
		},
	}
}

// GetAspectRatioMap creates and returns a new AspectRatioMap
func GetAspectRatioMap() *AspectRatioMap {
	return &AspectRatioMap{
		aspectRatios: map[AspectRatio]string{
			AspectRatio1_1:  "1:1",
			AspectRatio3_2:  "3:2",
			AspectRatio2_3:  "2:3",
			AspectRatio16_9: "16:9",
			AspectRatio9_16: "9:16",
		},
	}
}

// GetQueueTypeMap creates and returns a new QueueTypeMap
func GetQueueTypeMap() *QueueTypeMap {
	return &QueueTypeMap{
		queueTypes: map[QueueType]string{
			QueueTypeDefault:   "default",
			QueueTypeFast:      "fast",
			QueueTypeFlux1:     "flux1",
			QueueTypeDedicated: "dedicated",
			QueueTypeOther:     "other",
		},
	}
}

// GetFileAssociatedResourceMap creates and returns a new FileAssociatedResourceMap
func GetFileAssociatedResourceMap() *FileAssociatedResourceMap {
	return &FileAssociatedResourceMap{
		resources: map[FileAssociatedResource]string{
			FileAssociatedResourceUserAvatar:        "USER_AVATAR",
			FileAssociatedResourceUserCoverImage:    "USER_COVER_IMAGE",
			FileAssociatedResourceWorkspace:         "WORKSPACE",
			FileAssociatedResourceArticleCoverImage: "ARTICLE_COVER_IMAGE",
			FileAssociatedResourceArticleFile:       "ARTICLE_FILE",
			FileAssociatedResourceStyle:             "STYLE",
			FileAssociatedResourceSDWorkflow:        "SD_WORKFLOW",
			FileAssociatedResourceChatRoomThumbnail: "CHAT_ROOM_THUMBNAIL",
			FileAssociatedResourceSDModel:           "SD_MODEL",
			FileAssociatedResourceSDModelTraining:   "SD_MODEL_TRAINING",
			FileAssociatedResourcePromptLibrary:     "PROMPT_LIBRARY",
			FileAssociatedResourceNone:              "NONE",
		},
	}
}

// GetRecipeIdMap creates and returns a new RecipeIdMap
func GetRecipeIdMap() *RecipeIdMap {
	return &RecipeIdMap{
		recipeIds: map[RecipeId]string{
			RecipeIdImageGeneratorSimple: "image-generator-simple",
			RecipeIdRemix:                "remix",
			RecipeIdFaceEnhancer:         "face-enhancer",
			RecipeIdUpscaler:             "upscaler",
		},
	}
}

func (m *PromptStyleMap) Get(promptStyle PromptStyle) string {
	return m.promptStyles[promptStyle]
}

func (m *PromptStyleMap) ToStrings() []string {
	strings := make([]string, len(m.promptStyles))
	for promptStyle := range m.promptStyles {
		strings = append(strings, string(promptStyle))
	}
	return strings
}

// Get retrieves the string value for a given RecipeTaskStatus
func (m *RecipeTaskStatusMap) Get(status RecipeTaskStatus) string {
	return m.taskStatuses[status]
}

// ToStrings converts all RecipeTaskStatus keys to a string slice
func (m *RecipeTaskStatusMap) ToStrings() []string {
	strings := make([]string, 0, len(m.taskStatuses))
	for status := range m.taskStatuses {
		strings = append(strings, string(status))
	}
	return strings
}

// Get retrieves the string value for a given RecipeType
func (m *RecipeTypeMap) Get(recipeType RecipeType) string {
	return m.recipeTypes[recipeType]
}

// ToStrings converts all RecipeType keys to a string slice
func (m *RecipeTypeMap) ToStrings() []string {
	strings := make([]string, 0, len(m.recipeTypes))
	for recipeType := range m.recipeTypes {
		strings = append(strings, string(recipeType))
	}
	return strings
}

// Get retrieves the string value for a given AspectRatio
func (m *AspectRatioMap) Get(ratio AspectRatio) string {
	return m.aspectRatios[ratio]
}

// ToStrings converts all AspectRatio keys to a string slice
func (m *AspectRatioMap) ToStrings() []string {
	strings := make([]string, 0, len(m.aspectRatios))
	for ratio := range m.aspectRatios {
		strings = append(strings, string(ratio))
	}
	return strings
}

// Get retrieves the string value for a given QueueType
func (m *QueueTypeMap) Get(queueType QueueType) string {
	return m.queueTypes[queueType]
}

// ToStrings converts all QueueType keys to a string slice
func (m *QueueTypeMap) ToStrings() []string {
	strings := make([]string, 0, len(m.queueTypes))
	for queueType := range m.queueTypes {
		strings = append(strings, string(queueType))
	}
	return strings
}

// Get retrieves the string value for a given FileAssociatedResource
func (m *FileAssociatedResourceMap) Get(resource FileAssociatedResource) string {
	return m.resources[resource]
}

// ToStrings converts all FileAssociatedResource keys to a string slice
func (m *FileAssociatedResourceMap) ToStrings() []string {
	strings := make([]string, 0, len(m.resources))
	for resource := range m.resources {
		strings = append(strings, string(resource))
	}
	return strings
}

// Get retrieves the string value for a given RecipeId
func (m *RecipeIdMap) Get(recipeId RecipeId) string {
	return m.recipeIds[recipeId]
}

// ToStrings converts all RecipeId keys to a string slice
func (m *RecipeIdMap) ToStrings() []string {
	strings := make([]string, 0, len(m.recipeIds))
	for recipeId := range m.recipeIds {
		strings = append(strings, string(recipeId))
	}
	return strings
}
