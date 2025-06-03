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

var (
	RecipeTaskStatuses = []RecipeTaskStatus{
		RecipeTaskStatusQueued,
		RecipeTaskStatusRunning,
		RecipeTaskStatusCompleted,
		RecipeTaskStatusFailed,
		RecipeTaskStatusCancelled,
		RecipeTaskStatusCanceled,
	}
	RecipeTypes = []RecipeType{
		RecipeTypeNormal,
		RecipeTypeInpaint,
		RecipeTypeChain,
		RecipeTypeComfyui,
		RecipeTypeDescribe,
		RecipeTypeTurbo,
		RecipeTypeOther,
	}
	PromptStyles = []PromptStyle{
		PromptStyleBase,
		PromptStyleEnhance,
		PromptStyleAnime,
		PromptStylePhotographic,
		PromptStyleCinematic,
		PromptStyleAnalogFilm,
		PromptStyleDigitalArt,
		PromptStyleFantasyArt,
		PromptStyleLineArt,
		PromptStylePixelArt,
		PromptStyleArtstyleWatercolor,
		PromptStyleComicBook,
		PromptStyleNeonpunk,
		PromptStyle3DModel,
		PromptStyleMiscFairyTale,
		PromptStyleMiscGothic,
		PromptStylePhotoLongExposure,
		PromptStylePhotoTiltShift,
		PromptStyleLowpoly,
		PromptStyleOrigami,
		PromptStyleCraftClay,
		PromptStyleGameMinecraft,
	}
	AspectRatios = []AspectRatio{
		AspectRatio1_1,
		AspectRatio3_2,
		AspectRatio2_3,
		AspectRatio16_9,
		AspectRatio9_16,
	}
	QueueTypes = []QueueType{
		QueueTypeDefault,
		QueueTypeFast,
		QueueTypeFlux1,
		QueueTypeDedicated,
		QueueTypeOther,
	}
	FileAssociatedResources = []FileAssociatedResource{
		FileAssociatedResourceUserAvatar,
		FileAssociatedResourceUserCoverImage,
		FileAssociatedResourceWorkspace,
		FileAssociatedResourceArticleCoverImage,
		FileAssociatedResourceArticleFile,
		FileAssociatedResourceStyle,
		FileAssociatedResourceSDWorkflow,
		FileAssociatedResourceChatRoomThumbnail,
		FileAssociatedResourceSDModel,
		FileAssociatedResourceSDModelTraining,
		FileAssociatedResourcePromptLibrary,
		FileAssociatedResourceNone,
	}
	RecipeIds = []RecipeId{
		RecipeIdImageGeneratorSimple,
		RecipeIdRemix,
		RecipeIdFaceEnhancer,
		RecipeIdUpscaler,
	}
)
