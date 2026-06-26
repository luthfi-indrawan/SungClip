package controllers

import (
	"SungClip/internal/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func (c *controllers) VideoIngestion(
	ctx context.Context,
	request *types.RequestVideoIngestion,
) (response *types.ResponseVideoIngestion, err error) {

	startTime := time.Now()

	log.Printf(
		"VideoIngestion started, video_url=%s",
		request.VideoURL,
	)

	// VALIDATE REQUEST

	log.Println("Validating request...")

	u, err := url.ParseRequestURI(request.VideoURL)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid video url: %w",
			err,
		)
	}

	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid video url")
	}

	if request.ClipsCount <= 0 {
		return nil, fmt.Errorf("invalid clips count")
	}

	if request.MinDurationSecond <= 0 {
		return nil, fmt.Errorf("invalid min duration")
	}

	if request.MaxDurationSecond <= 0 {
		return nil, fmt.Errorf("invalid max duration")
	}

	log.Println("Request validated")

	// DOWNLOAD VIDEO

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	log.Println("Downloading video...")

	downloadsDir := c.utils.BuildDownloadsDir()

	if err := c.utils.MkdirAll(downloadsDir); err != nil {
		return nil, fmt.Errorf(
			"failed create downloads dir: %w",
			err,
		)
	}

	videoPath, infoVideoPath, err := c.services.DownloadVideo(
		ctx,
		request.VideoURL,
		downloadsDir,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed download video: %w",
			err,
		)
	}

	// LOAD VIDEO METADATA

	var metadataVideo types.MetadataVideo

	if err := c.utils.ReadAndParse(
		infoVideoPath,
		&metadataVideo,
	); err != nil {
		return nil, fmt.Errorf(
			"failed read metadata: %w",
			err,
		)
	}

	log.Printf(
		"Video downloaded successfully, title=%s",
		metadataVideo.Title,
	)

	baseName := c.utils.NormalizeTitle(
		metadataVideo.Title,
	)

	// EXTRACT AUDIO

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	log.Println("Extracting audio...")

	audioDir := c.utils.BuildAudioDir(baseName)

	if err := c.utils.MkdirAll(audioDir); err != nil {
		return nil, fmt.Errorf(
			"failed create audio dir: %w",
			err,
		)
	}

	audioPath := filepath.Join(
		audioDir,
		baseName+".wav",
	)

	if err := c.services.ExtractAudio(
		ctx,
		videoPath,
		audioPath,
	); err != nil {
		return nil, fmt.Errorf(
			"failed extract audio: %w",
			err,
		)
	}

	log.Println("Audio extracted successfully")

	// TRANSCRIBE

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	log.Println("Transcribing audio...")

	transcriptDir := c.utils.BuildTranscriptDir(
		baseName,
	)

	if err := c.utils.MkdirAll(transcriptDir); err != nil {
		return nil, fmt.Errorf(
			"failed create transcript dir: %w",
			err,
		)
	}

	transcriptPath := filepath.Join(
		transcriptDir,
		baseName+".json",
	)

	if err := c.services.Transcribe(
		ctx,
		audioPath,
		transcriptPath,
	); err != nil {
		return nil, fmt.Errorf(
			"failed transcribe audio: %w",
			err,
		)
	}

	log.Println("Transcription completed")

	// LOAD TRANSCRIPT

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	log.Println("Loading transcript...")

	transcriptBytes, err := os.ReadFile(
		transcriptPath,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed read transcript file: %w",
			err,
		)
	}

	var segments []types.Segment

	if err := json.Unmarshal(
		transcriptBytes,
		&segments,
	); err != nil {
		return nil, fmt.Errorf(
			"failed unmarshal transcript: %w",
			err,
		)
	}

	metadataVideo.TranscriptResult = segments

	log.Printf(
		"Transcript loaded, segments=%d",
		len(segments),
	)

	// BUILD PROMPT

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	log.Println("Building prompt...")

	promptString := c.services.BuildPrompt(
		metadataVideo, request.ClipsCount, request.MinDurationSecond, request.MaxDurationSecond,
	)

	promptDir := c.utils.BuildPromptDir(
		baseName,
	)

	if err := c.utils.MkdirAll(promptDir); err != nil {
		return nil, fmt.Errorf(
			"failed create prompt dir: %w",
			err,
		)
	}

	promptPath := filepath.Join(
		promptDir,
		"prompt.txt",
	)

	if err := c.utils.WriteFile(
		promptPath,
		[]byte(promptString),
	); err != nil {
		return nil, fmt.Errorf(
			"failed write prompt file: %w",
			err,
		)
	}

	log.Println("Prompt created successfully")

	// SAVE METADATA

	metadataBytes, err := json.MarshalIndent(
		metadataVideo,
		"",
		"  ",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed marshal metadata: %w",
			err,
		)
	}

	metadataPath := filepath.Join(
		promptDir,
		"metadata.json",
	)

	if err := c.utils.WriteFile(
		metadataPath,
		metadataBytes,
	); err != nil {
		return nil, fmt.Errorf(
			"failed write metadata file: %w",
			err,
		)
	}

	log.Printf(
		"VideoIngestion completed successfully, duration=%s",
		time.Since(startTime),
	)

	return &types.ResponseVideoIngestion{
		Title:        baseName,
		PromptPath:   promptPath,
		MetadataPath: metadataPath,
	}, nil
}