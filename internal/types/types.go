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
	MomentClip struct {
		StartMS        int64    `json:"start_ms"`
		EndMS          int64    `json:"end_ms"`
		Score          int      `json:"score"`
		Title          string   `json:"title"`
		Caption        string   `json:"caption"`
		WordHighlights []string `json:"word_highlights"`
	}
)

type (
	MetadataClip struct {
		Title          string   `json:"title"`
		Width          int      `json:"width"`
		Height         int      `json:"height"`
		Duration       int      `json:"duration"`
		CompositionID  string   `json:"compotion_id"`
		VideoPath      string   `json:"video_path"`
		Caption        string   `json:"caption"`
		Subtitle       []Word   `json:"subtitle"`
		WordHighlights []string `json:"word_highlights"`
	}
)