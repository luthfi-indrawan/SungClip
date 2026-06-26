package types

type (
	MetadataVideo struct {
		ID               string       `json:"id"`
		Title            string       `json:"title"`
		Channel          string       `json:"channel"`
		Language         string       `json:"language"`
		Width            int          `json:"width"`
		Height           int          `json:"height"`
		TranscriptResult []Segment    `json:"transcript_result"`
		MomentsForClip   []MomentClip `json:"moments_for_clip"`
	}

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

	MomentClip struct {
		StartMS        int64    `json:"start_ms"`
		EndMS          int64    `json:"end_ms"`
		Score          int      `json:"score"`
		Title          string   `json:"title"`
		Headline       string   `json:"headline"`
		Caption        string   `json:"caption"`
		WordHighlights []string `json:"word_highlights"`
		Hashtags       []string `json:"hashtags"`
	}
)

// types for clip
type (
	FaceTrackerMetadata struct {
		VideoWidth  int64              `json:"videoWidth"`
		VideoHeight int64              `json:"videoHeight"`
		FPS         float64            `json:"fps"`
		TotalFrames int64              `json:"totalFrames"`
		Frames      []FrameFaceTracker `json:"frames"`
	}

	FrameFaceTracker struct {
		Frame  int64   `json:"frame"`
		TimeMS int64   `json:"timeMs"`
		Tracks []Track `json:"tracks"`
	}

	Track struct {
		TrackID    int     `json:"trackId"`
		CenterX    float64 `json:"centerX"`
		CenterY    float64 `json:"centerY"`
		Width      float64 `json:"width"`
		Height     float64 `json:"height"`
		Confidence float64 `json:"confidence"`
	}
)

type (
	HookTTSResult struct {
		AudioPath  string `json:"audio_path"`
		DurationMs int64  `json:"duration_ms"`
	}

	HookMetadata struct {
		Text       string `json:"text"`
		AudioPath  string `json:"audio_path"`
		DurationMs int64  `json:"duration_ms"`
	}
)

type (
	MetadataClip struct {
		Title             string             `json:"title"`
		Headline          string             `json:"headline"`
		Hook              HookMetadata       `json:"hook"`
		FPS               float64            `json:"fps"`
		TargetWidth       int                `json:"target_width"`
		TargetHeight      int                `json:"target_height"`
		OriWidth          int                `json:"ori_width"`
		OriHeight         int                `json:"ori_height"`
		TotalFrames       int                `json:"total_frames"`
		CompositionID     string             `json:"composition_id"`
		VideoPath         string             `json:"video_path"`
		Caption           string             `json:"caption"`
		Subtitle          []Word             `json:"subtitle"`
		Hashtags          []string           `json:"hashtags"`
		WordHighlights    []string           `json:"word_highlights"`
		FramesFaceTracker []FrameFaceTracker `json:"frames_face_trackers"`
	}
)