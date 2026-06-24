package controllers

import (
	"SungClip/internal/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

func (c *controllers) VideoEditing(
	ctx context.Context,
	request *types.RequestVideoEditing,
) (response *types.ResponseVideoEditing, err error) {
	startTime := time.Now()

	log.Printf("VideoEditing running with request: %+v", request)

	// validasi request
	log.Println("Validating request...")

	if request.Title == "" {
		return nil, errors.New("invalid request: missing title")
	}

	if request.CompositionID == "" {
		return nil, errors.New("invalid request: missing compostion id")
	}

	if request.Width <= 0 {
		return nil, errors.New("invalid request: min width 1")
	}

	if request.Height <= 0 {
		return nil, errors.New("invalid request: min height 1")
	}

	log.Println("Request validated")

	// build path or dir
	transcriptPath := filepath.Join(c.utils.BuildTranscriptDir(request.Title), request.Title+".json")
	// videoPath := filepath.Join(c.utils.BuildDownloadsDir(), request.Title+".webm")
	momentsPath := filepath.Join(c.utils.BuildMomentDir(request.Title), request.Title+".json")

	remotionPublicDir := c.utils.BuildRemotionPublicDir(request.Title)
	if err := c.utils.MkdirAll(remotionPublicDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	metadataDir := c.utils.BuildMetadataDir(request.Title)
	if err := c.utils.MkdirAll(metadataDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	faceTrackerDir := c.utils.BuildFaceTrackerDir(request.Title)
	if err := c.utils.MkdirAll(faceTrackerDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	// load transcript
	transcriptBytes, err := os.ReadFile(transcriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed read transcript file: %w", err)
	}

	var transcript types.TranscriptResult
	if err := json.Unmarshal(transcriptBytes, &transcript); err != nil {
		return nil, fmt.Errorf("failed unmarshal bytes transcript: %w", err)
	}

	// extract words
	var words []types.Word
	for _, segment := range transcript {
		words = append(words, segment.Words...)
	}

	// load moment
	momentBytes, err := os.ReadFile(momentsPath)
	if err != nil {
		return nil, fmt.Errorf("failed read moment file: %w", err)
	}

	var momentClips types.MomentClips
	if err := json.Unmarshal(momentBytes, &momentClips); err != nil {
		return nil, fmt.Errorf("failed unmarshal bytes transcript: %w", err)
	}

	var (
		metadataPathClips []string
		outputPathClips   []string
	)

	// looping moment clips
	for i, momentClip := range momentClips {
		newStart, newEnd := c.services.ExpandDurationClip(momentClip.StartMS, momentClip.EndMS)

		log.Printf("Processing clip %d/%d | title=%s start=%d end=%d", i+1, len(momentClips), momentClip.Title, newStart, newEnd)

		outputMetdataPath := filepath.Join(metadataDir, fmt.Sprintf("%02d_%s.json", i+1, c.SafeFilename(momentClip.Title)))
		// outputVideoPath := filepath.Join(remotionPublicDir, fmt.Sprintf("%02d_%s.mp4", i+1, c.SafeFilename(momentClip.Title)))
		outputVideoPublic := filepath.Join(request.Title, fmt.Sprintf("%02d_%s.mp4", i+1, c.SafeFilename(momentClip.Title)))
		outputFaceTrackerPath := filepath.Join(faceTrackerDir, fmt.Sprintf("%02d_%s.json", i+1, c.SafeFilename(momentClip.Title))) 

		log.Printf("Cutting video | from: %d to: %d", newStart, newEnd)

		// if err := c.services.CutVideo(videoPath, outputVideoPath, newStart, newEnd); err != nil {
		// 	return nil, fmt.Errorf("failed cutting video: %w", err)
		// }

		log.Printf("Success cut video | title: %s", momentClip.Title)

		log.Printf("Tracking face on video | from: %d to: %d", newStart, newEnd)

		// if err := c.services.FaceTracking(outputVideoPath, outputFaceTrackerPath); err != nil {
		// 	return nil, fmt.Errorf("failed tracking face on video: %w", err)
		// }

		log.Printf("Face tracked | title: %s", momentClip.Title)

		// load face tracker data
		momentBytes, err := os.ReadFile(outputFaceTrackerPath)
		if err != nil {
			return nil, fmt.Errorf("failed read face tracker file: %w", err)
		}

		var faceTrackerMetadata types.FaceTrackerMetadata
		if err := json.Unmarshal(momentBytes, &faceTrackerMetadata); err != nil {
			return nil, fmt.Errorf("failed unmarshal bytes: %w", err)
		}

		log.Printf("Generating metadata | title: %s...", momentClip.Title)

		if err := c.services.GenerateMetadataVideo(
			momentClip.Title, request.Width, request.Height, request.CompositionID, outputVideoPublic, momentClip.Caption,
			words, momentClip.WordHighlights, outputMetdataPath, newStart, newEnd, faceTrackerMetadata,
		); err != nil {
			return nil, fmt.Errorf("failed generate metadata: %w", err)
		}

		log.Printf("Metadata generated | title: %s", momentClip.Title)

		metadataPathClips = append(metadataPathClips, outputMetdataPath)

		outputPath := filepath.Join(c.utils.BuildResultDir(request.Title), fmt.Sprintf("%02d_%s.mp4", i+1, c.SafeFilename(momentClip.Title)))

		outputPathClips = append(outputPathClips, outputPath)
	}

	// render remotion
	log.Println("Starting remotion rendering...")

	for i := 0; i < len(metadataPathClips); i++ {
		metadataPath := metadataPathClips[i]
		outputPath := outputPathClips[i]

		log.Printf("Remotion render start | metadata=%s output=%s", metadataPath, outputPath)

		if err := c.services.ExecuteRemotion(metadataPath, outputPath); err != nil {
			return nil, fmt.Errorf("failed render clip %d: %w", i+1, err)
		}

		log.Printf("Remotion render completed | clip=%d output=%s", i+1, outputPath)
	}

	// response
	log.Printf("VideoEditing completed successfully, total_duration:%s", time.Since(startTime))

	return &types.ResponseVideoEditing{
		Title: request.Title,
		TotalClips: len(metadataPathClips),
		ResultVideoPath: strings.Join(outputPathClips, " | "),
	}, nil
}

func (c *controllers) SafeFilename(name string) string {
	return slug.MakeLang(name, "id")
}