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
	"strings"
	"time"
)

func (c *controllers) VideoIngestion(
	ctx context.Context,
	request *types.RequestVideoIngestion,
) (response *types.ResponseVideoIngestion, err error) {
	startTime := time.Now()

	log.Printf("VideoIngestion running with request: %+v", request)

	// validasi request
	log.Println("Validating request...")

	if _, err := url.ParseRequestURI(request.VideoURL); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	log.Println("Request validated")
	
	// download video
	log.Println("Downloading video...")

	downloadsDir := c.utils.BuildDownloadsDir()

	if err :=  c.utils.MkdirAll(downloadsDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	videoPath, err := c.services.DownloadVideo(request.VideoURL, downloadsDir)
	if err != nil {
		return nil, fmt.Errorf("failed download video: %w", err)
	}

	log.Println("Video Downloaded: ", videoPath)

	// extract name
	baseName := strings.TrimSuffix(
		filepath.Base(videoPath),
		filepath.Ext(videoPath),
	)

	log.Printf("Video basename extracted | name=%s", baseName)

	// extract audio
	log.Println("Extract audio...")

	audioDir := c.utils.BuildAudioDir(baseName)

	if err := c.utils.MkdirAll(audioDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	audioPath := filepath.Join(audioDir, baseName+".wav")

	if err := c.services.ExtractAudio(videoPath, audioPath); err != nil {
		return nil, fmt.Errorf("failed extract audio: %w", err)
	}

	log.Println("Audio extracted")

	// transcribe
	log.Println("Transcribe...")

	transcriptDir := c.utils.BuildTranscriptDir(baseName)

	if err := c.utils.MkdirAll(transcriptDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	transcriptPath := filepath.Join(transcriptDir, baseName+".json")

	if err := c.services.Transcribe(audioPath, transcriptPath); err != nil {
		return nil, fmt.Errorf("failed transcribe: %w", err)
	}

	log.Println("Transcribe successfully")

	// load transcript
	transcriptBytes, err := os.ReadFile(transcriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed read transcript file: %w", err)
	}

	var transcript types.TranscriptResult
	if err := json.Unmarshal(transcriptBytes, &transcript); err != nil {
		return nil, fmt.Errorf("failed unmarshal bytes transcript: %w", err)
	}

	// build prompt
	log.Println("Build prompt...")

	promptString := c.services.BuildPrompt(transcript)

	promptDir := c.utils.BuildPromptDir(baseName)

	if err := c.utils.MkdirAll(promptDir); err != nil {
		return nil, fmt.Errorf("failed mkdir: %w", err)
	}

	promptPath := filepath.Join(promptDir, "prompt.txt")

	if err := c.utils.WriteFile(promptPath, []byte(promptString)); err != nil {
		return nil, fmt.Errorf("failed write file: %w", err)
	}

	log.Println("Build prompt successfully")

	// response
	log.Printf("VideoIngestion completed successfully, total_duration:%s", time.Since(startTime))
	return &types.ResponseVideoIngestion{
		Title: baseName,
		PromptPath: promptPath,
	}, nil
}