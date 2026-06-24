package config

import (
	"log"
	"os"
	"path/filepath"
)

const (
	storageDir = "storage"

	remotionDir = "remotion"
	publicDir = "public"

	binDir = "bin"
	ytdlp      = "yt-dlp.exe"
	ffmpeg     = "ffmpeg.exe"

	scriptsDir = "scripts"
	transcriptDir = "transcript"
	faceTrackerDir = "face-tracker"
)

type (
	Config struct {
		// paths
		StoragePath        string
		RemotionPublicPath string

		// scripts
		YTDLP        string
		FFMPEG       string
		PYEXETranscript        string
		PYTranscript string
		PYEXEFaceTracker        string
		PYFaceTracker string
	}
)

func NewConfig() *Config {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed get working directory: %v", err)
	}

	return &Config{
		StoragePath: filepath.Join(wd, storageDir),
		RemotionPublicPath: filepath.Join(wd, remotionDir, publicDir),

		YTDLP: filepath.Join(wd, binDir, ytdlp),
		FFMPEG: filepath.Join(wd, binDir, ffmpeg),
		PYEXETranscript: filepath.Join(wd, scriptsDir, transcriptDir, ".venv", "Scripts", "python.exe"),
		PYTranscript: filepath.Join(wd, scriptsDir, transcriptDir, "main.py"),
		PYEXEFaceTracker: filepath.Join(wd, scriptsDir, faceTrackerDir, "venv", "Scripts", "python.exe"),
		PYFaceTracker: filepath.Join(wd, scriptsDir, faceTrackerDir, "main.py"),
	}
}