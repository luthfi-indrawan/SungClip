package services

import "SungClip/internal/types"

type IServices interface {
	DownloadVideo(url string, outputDir string) (videoPath string, err error)
	ExtractAudio(videoPath string, outputDir string) error
	Transcribe(audioPath string, outputPath string) error
	BuildPrompt(transcript types.TranscriptResult) string

	ExpandDurationClip(startMS int64, endMS int64) (newStart int64, newEnd int64)
	GenerateMetadataVideo(title string, width int, height int, CompositionID string, videoPath string, caption string, words []types.Word, wordHighlights []string, outputPath string, startMS int64, endMS int64) error
	CutVideo(inputPath string, outputPath string, startMS int64, endMS int64) error
	ExecuteRemotion(inputPropsPath string, outputClipsPath string) error
}