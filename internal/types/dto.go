package types

type (
	RequestVideoIngestion struct {
		VideoURL          string
		ClipsCount        int
		MinDurationSecond int
		MaxDurationSecond int
	}

	ResponseVideoIngestion struct {
		Title        string
		PromptPath   string
		MetadataPath string
	}
)

type (
	RequestVideoEditing struct {
		Title         string
		CompositionID string
		Resolution    string // enum: hd, fhd, 2k, 4k
	}

	ResponseVideoEditing struct {
		Title           string
		TotalClips      int
		ResultVideoPath string
	}
)