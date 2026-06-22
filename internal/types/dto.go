package types

type (
	RequestVideoIngestion struct {
		VideoURL string
	}

	ResponseVideoIngestion struct {
		Title      string
		PromptPath string
	}
)

type (
	RequestVideoEditing struct {
		Title         string
		Width         int
		Height        int
		CompositionID string
	}

	ResponseVideoEditing struct {
		Title           string
		TotalClips      int
		ResultVideoPath string
	}
)