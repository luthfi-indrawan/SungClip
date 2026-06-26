package services

import (
	"SungClip/internal/types"
	"context"
)

type IServices interface {
	DownloadVideo(ctx context.Context, url string, outputDir string) (videoPath, infoVideoPath string, err error)
	ExtractAudio(ctx context.Context, videoPath string, outputPath string) error
	Transcribe(ctx context.Context, audioPath string, outputPath string) error
	BuildPrompt(metadataVideo types.MetadataVideo, clipsCount, minDuration, maxDuration int) string

	ParseResolution(width int, height int, resolution string) (newWidth int, newHeight int, err error)
	ParseCompositionResolution(compositionID string, resolution string) (width int, height int, err error)

	GenerateHookAudio(ctx context.Context, text string, outputPath  string) (*types.HookTTSResult, error)
	CutVideo(ctx context.Context, inputPath string, outputPath string, width int, height int, startMS int64, endMS int64) error
	FaceTracking(ctx context.Context, videoPath string, outputPath string) error
	GenerateMetadataVideo(width int, height int, CompositionID string, words []types.Word, moment types.MomentClip, faceTrackerMetadata types.FaceTrackerMetadata, hook types.HookTTSResult, videoPathClip string, audioPathClip string, outputMetadataPath string) error
	ExecuteRemotion(ctx context.Context, inputPropsPath string, outputClipsPath string, compostionID string) error
}