package services

import (
	"SungClip/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

func (s *services) ParseResolution(
	width int,
	height int,
	resolution string,
) (newWidth int, newHeight int, err error) {
	if width <= 0 || height <= 0 {
		return 0, 0, fmt.Errorf("invalid width or height")
	}

	var targetLongSide int

	switch strings.ToLower(resolution) {
	case "hd":
		targetLongSide = 1280

	case "fhd":
		targetLongSide = 1920

	case "2k":
		targetLongSide = 2560

	case "4k":
		targetLongSide = 3840

	default:
		return 0, 0, fmt.Errorf("unknown resolution: %s", resolution)
	}

	aspectRatio := float64(width) / float64(height)

	// Landscape
	if width >= height {
		newWidth = targetLongSide
		newHeight = int(math.Round(float64(newWidth) / aspectRatio))
	} else {
		// Portrait
		newHeight = targetLongSide
		newWidth = int(math.Round(float64(newHeight) * aspectRatio))
	}

	// Biar aman untuk encoder video (genap)
	newWidth -= newWidth % 2
	newHeight -= newHeight % 2

	return
}

func (s *services) ParseCompositionResolution(
	compositionID string,
	resolution string,
) (width int, height int, err error) {
	switch strings.ToLower(compositionID) {

	case "single-podcast":
		switch strings.ToLower(resolution) {
		case "hd":
			return 1080, 1350, nil
		case "fhd":
			return 1440, 1800, nil
		case "2k":
			return 1920, 2400, nil
		case "4k":
			return 3840, 4800, nil
		}
	}

	return 0, 0, fmt.Errorf(
		"unsupported composition=%s resolution=%s",
		compositionID,
		resolution,
	)
}

func (s *services) GenerateMetadataVideo(width int, height int, CompositionID string, words []types.Word, moment types.MomentClip, faceTrackerMetadata types.FaceTrackerMetadata, hook types.HookTTSResult, videoPathClip string, audioPathClip string, outputMetadataPath string) error {
	var subtitle []types.Word

	for _, word := range words {
		if word.EndMS >= moment.StartMS && word.StartMS <= moment.EndMS {
			subtitle = append(subtitle, word)
		}
	}

	if len(subtitle) > 0 {
		baseTime := subtitle[0].StartMS

		for i := range subtitle {
			subtitle[i].StartMS -= baseTime
			subtitle[i].EndMS -= baseTime
		}
	}

	metadataClip := types.MetadataClip{
		Title: moment.Title,
		Headline: moment.Headline,
		Hook: types.HookMetadata{
			Text: moment.Headline,
			AudioPath: audioPathClip,
			DurationMs: hook.DurationMs,
		},
		FPS: faceTrackerMetadata.FPS,
		TargetWidth: width,
		TargetHeight: height,
		OriWidth: int(faceTrackerMetadata.VideoWidth),
		OriHeight: int(faceTrackerMetadata.VideoHeight),
		TotalFrames: int(faceTrackerMetadata.TotalFrames),
		CompositionID: CompositionID,
		VideoPath: videoPathClip,
		Caption: moment.Caption,
		Subtitle: subtitle,
		Hashtags: moment.Hashtags,
		WordHighlights: moment.WordHighlights,
		FramesFaceTracker: faceTrackerMetadata.Frames,
	}

	metadataBytes, err := json.MarshalIndent(metadataClip, "", "  ")
	if err != nil {
		return err
	}

	if err := s.utils.WriteFile(outputMetadataPath, metadataBytes); err != nil {
		return err
	}

	return  nil
}

func (s *services) GenerateHookAudio(ctx context.Context, text string, outputPath  string) (*types.HookTTSResult, error) {
	cmd := exec.CommandContext(
		ctx,
		s.utils.GetPyEXEGenerateHook(),
		s.utils.GetPyGenerateHook(),

		"--text",
		text,

		"--output",
		outputPath,

		"--voice",
		"id-ID-ArdiNeural",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"generate hook audio failed: %w\n%s",
			err,
			string(output),
		)
	}

	var result types.HookTTSResult

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf(
			"unmarshal hook result failed: %w\nraw=%s",
			err,
			string(output),
		)
	}

	return &result, nil
}

func (s *services) CutVideo(ctx context.Context, inputPath string, outputPath string, width int, height int, startMS int64, endMS int64) error {
	cmd := exec.CommandContext(
		ctx,
		s.utils.GetFFMPEG(),
		"-y",
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=%d:%d", width, height),
		"-ss", fmt.Sprintf("%.3f", float64(startMS)/1000),
		"-to", fmt.Sprintf("%.3f", float64(endMS)/1000),
		"-c:v", "libx264",
		"-c:a", "aac",
		outputPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (s *services) FaceTracking(ctx context.Context, videoPath string, outputPath string) error {
	cmd := exec.CommandContext(
		ctx,
		s.utils.GetPyEXEFaceTracker(),
		s.utils.GetPyFaceTracker(),
		videoPath,
		outputPath,
	)

	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return err
    }

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("face tracker file not found: %w", err)
	}

    return nil
}

func (s *services) ExecuteRemotion(ctx context.Context, inputPropsPath string, outputClipsPath string, compostionID string) error {
	cmd := exec.CommandContext(
		ctx,
		"npx",
		"remotion",
		"render",
		compostionID,
		outputClipsPath,
		"--crf=14",
		fmt.Sprintf("--props=%s", inputPropsPath),
	)

	// pindah ke project remotion
	cmd.Dir = "remotion"

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}