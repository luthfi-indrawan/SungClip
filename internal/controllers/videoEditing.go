package controllers

import (
	"SungClip/internal/types"
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"
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

	if request.Resolution == "" {
		return nil, errors.New("invalid request: missing resolution")
	}

	log.Println("Request validated")

	// load metadata video
	baseName := c.utils.NormalizeTitle(request.Title)
	metadataPath := filepath.Join(c.utils.BuildPromptDir(baseName), "metadata.json")

	var metadataVideo types.MetadataVideo
	if err := c.utils.ReadAndParse(metadataPath, &metadataVideo); err != nil {
		return nil, fmt.Errorf("failed read and parse: %w", err)
	}

	newWidth, newHeight, err := c.services.ParseResolution(metadataVideo.Width, metadataVideo.Height, request.Resolution)
	if err != nil {
		return nil, fmt.Errorf("failed parser resolution: %w", err)
	}

	metadataVideo.Width = newWidth
	metadataVideo.Height = newHeight

	// extract words
	var words []types.Word
	for _, segment := range metadataVideo.TranscriptResult {
		words = append(words, segment.Words...)
	}

	// prepare dir
	metadataClipsDir := c.utils.BuildMetadataDir(baseName)
	if err := c.utils.MkdirAll(metadataClipsDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	remotionPublicDir := c.utils.BuildRemotionPublicDir(baseName)
	if err := c.utils.MkdirAll(remotionPublicDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	faceTrackerDir := c.utils.BuildFaceTrackerDir(baseName)
	if err := c.utils.MkdirAll(faceTrackerDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	resultDir := c.utils.BuildResultDir(baseName)
	if err := c.utils.MkdirAll(resultDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	// videoPath := filepath.Join(c.utils.BuildDownloadsDir(), baseName+".webm")

	var metadataClipsPath []string
	var outputPathClips []string

	targetWidth, targetHeight, err := c.services.ParseCompositionResolution(request.CompositionID, request.Resolution)
	if err != nil {
		return nil, fmt.Errorf("failed parse compotion resolution: %w", err)
	}

	for i, moment := range metadataVideo.MomentsForClip {
		log.Printf("Processing clip %d/%d", i+1, len(metadataVideo.MomentsForClip))

		titleClip := c.utils.NormalizeTitle(moment.Title)

		outputMetadataClipPath := filepath.Join(metadataClipsDir, titleClip+".json")
		// outputVideoClipPath := filepath.Join(remotionPublicDir, titleClip+".mp4")
		outputAudioPath := filepath.Join(remotionPublicDir, titleClip+".mp3")
		videoClipPathForMetadata := filepath.Join(baseName, titleClip+".mp4")
		audioClipPathForMetadata := filepath.Join(baseName, titleClip+".mp3")
		outputFaceTrackerPath := filepath.Join(faceTrackerDir, titleClip+".json")

		log.Println("Cutting video")

		// if err := c.services.CutVideo(ctx, videoPath, outputVideoClipPath, newWidth, newHeight, moment.StartMS, moment.EndMS); err != nil {
		// 	return nil, fmt.Errorf("failed cutting video: %w", err)
		// }

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		log.Println("Success cut video")

		log.Println("Tracking face on video")

		// if err := c.services.FaceTracking(ctx, outputVideoClipPath, outputFaceTrackerPath); err != nil {
		// 	return nil, fmt.Errorf("failed tracking face on video: %w", err)
		// }

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var faceTracker types.FaceTrackerMetadata
		if err := c.utils.ReadAndParse(outputFaceTrackerPath, &faceTracker); err != nil {
			return nil, fmt.Errorf("failed read and parse: %w", err)
		}

		log.Println("Face tracked")

		log.Println("Generate hook")

		hookRes, err := c.services.GenerateHookAudio(ctx, moment.Headline, outputAudioPath)
		if err != nil {
			return nil, fmt.Errorf("failed generate hook: %w", err)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		log.Println("Hook generated")

		log.Println("Generate metadata")

		if err := c.services.GenerateMetadataVideo(targetWidth, targetHeight, request.CompositionID, words, moment, faceTracker, *hookRes, videoClipPathForMetadata, audioClipPathForMetadata, outputMetadataClipPath); err != nil {
			return nil, fmt.Errorf("failed generate metadata: %w", err)
		}

		log.Println("Metadata generated")

		outputPathClip := filepath.Join(resultDir, c.utils.NormalizeTitle(moment.Title)+".mp4")

		metadataClipsPath = append(metadataClipsPath, outputMetadataClipPath)
		outputPathClips = append(outputPathClips, outputPathClip)
	}

	// render remotion
	log.Println("Starting remotion rendering...")

	for i := 0; i < len(metadataClipsPath); i++ {
		metadataPath := metadataClipsPath[i]
		outputPath := outputPathClips[i]

		log.Printf("Remotion render start | metadata=%s output=%s", metadataPath, outputPath)

		if err := c.services.ExecuteRemotion(ctx, metadataPath, outputPath, request.CompositionID); err != nil {
			return nil, fmt.Errorf("failed render clip %d: %w", i+1, err)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		log.Printf("Remotion render completed | clip=%d output=%s", i+1, outputPath)
	}

	// response
	log.Printf("VideoEditing completed successfully, total_duration:%s", time.Since(startTime))

	return &types.ResponseVideoEditing{
		Title: request.Title,
		TotalClips: len(metadataClipsPath),
		ResultVideoPath: strings.Join(outputPathClips, " | "),
	}, nil
}