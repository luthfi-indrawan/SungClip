package types

type TranscriptResult []Segment
type MomentClips []MomentClip
type MetadataClips []MetadataClip

type (
	Segment struct {
		StartMS int64  `json:"start_ms"`
		EndMS   int64  `json:"end_ms"`
		Text    string `json:"text"`
		Words   []Word `json:"words"`
	}

	Word struct {
		Word    string `json:"word"`
		StartMS int64  `json:"start_ms"`
		EndMS   int64  `json:"end_ms"`
	}
)

type (
	FaceTrackerMetadata struct {
		VideoWidth  int64              `json:"videoWidth"`
		VideoHeight int64              `json:"videoHeight"`
		FPS         float64            `json:"fps"`
		TotalFrames int64              `json:"totalFrames"`
		Frames      []FrameFaceTracker `json:"frames"`
	}

	FrameFaceTracker struct {
		Frame   int64   `json:"frame"`
		TimeMS  int64   `json:"timeMs"`
		CenterX float64 `json:"centerX"`
		CenterY float64 `json:"centerY"`
		Width   float64 `json:"width"`
		Height  float64 `json:"height"`
	}
)

type (
	MomentClip struct {
		StartMS        int64    `json:"start_ms"`
		EndMS          int64    `json:"end_ms"`
		Score          int      `json:"score"`
		Title          string   `json:"title"`
		Headline       string   `json:"headline"`
		Caption        string   `json:"caption"`
		WordHighlights []string `json:"word_highlights"`
		Hastags        []string `json:"hastags"`
	}
)

type (
	MetadataClip struct {
		Title             string             `json:"title"`
		Headline          string             `json:"headline"`
		FPS               int                `json:"fps"`
		TargetWidth       int                `json:"target_width"`
		TargetHeight      int                `json:"target_height"`
		OriWidth          int                `json:"ori_width"`
		OriHeight         int                `json:"ori_height"`
		TotalFrames       int                `json:"total_frames"`
		CompositionID     string             `json:"compotion_id"`
		VideoPath         string             `json:"video_path"`
		Caption           string             `json:"caption"`
		Subtitle          []Word             `json:"subtitle"`
		Hastags           []string           `json:"hastags"`
		WordHighlights    []string           `json:"word_highlights"`
		FramesFaceTracker []FrameFaceTracker `json:"frames_face_trackers"`
	}
)