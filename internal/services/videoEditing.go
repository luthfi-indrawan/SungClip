package services

import (
	"SungClip/internal/types"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func (s *services) ExpandDurationClip(startMS int64, endMS int64) (newStart int64, newEnd int64) {
	padding := int64(5000)

	startMS -= padding + 3000
	endMS += padding

	if startMS < 0 {
		startMS = 0
	}

	return startMS, endMS
}

func (s *services) GenerateMetadataVideo(title string, width int, height int, CompositionID string, videoPath string, caption string, words []types.Word, wordHighlights []string, outputPath string, startMS int64, endMS int64, faceTrackerMetadata types.FaceTrackerMetadata) error {
	var subtitle []types.Word

	for _, word := range words {
		if word.StartMS >= startMS && word.EndMS <= endMS {
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
		Title: title,
		FPS: int(faceTrackerMetadata.FPS),
		TargetWidth: width,
		TargetHeight: height,
		OriWidth: int(faceTrackerMetadata.VideoWidth),
		OriHeight: int(faceTrackerMetadata.VideoHeight),
		TotalFrames: int(faceTrackerMetadata.TotalFrames),
		CompositionID: CompositionID,
		VideoPath: videoPath,
		Caption: caption,
		Subtitle: subtitle,
		WordHighlights: wordHighlights,
		FramesFaceTracker: faceTrackerMetadata.Frames,
	}

	metadataBytes, err := json.MarshalIndent(metadataClip, "", "  ")
	if err != nil {
		return err
	}

	if err := s.utils.WriteFile(outputPath, metadataBytes); err != nil {
		return err
	}

	return  nil
}

func (s *services) CutVideo(inputPath string, outputPath string, startMS int64, endMS int64) error {
	cmd := exec.Command(
		s.utils.GetFFMEPG(),
		"-y",
		"-i", inputPath,
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

func (s *services) FaceTracking(videoPath string, outputPath string) error {
	cmd := exec.Command(
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

func (s *services) ExecuteRemotion(inputPropsPath string, outputClipsPath string) error {
	cmd := exec.Command(
		"npx",
		"remotion",
		"render",
		"clip",
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