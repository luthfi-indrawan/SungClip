package utils

import (
	"SungClip/internal/config"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Utils struct {
	cfg *config.Config
}

func NewUtils(cfg *config.Config) *Utils {
	return &Utils{
		cfg: cfg,
	}
}

const (
	downloadsDir   = "downloads"
	audioDir       = "audio"
	transcriptDir  = "transcript"
	momentDir      = "moment"
	metadataDir    = "metadata"
	resultDir      = "result"
	faceTrackerDir = "face-tracker"

	DefaultDirPerm  = 0755
	DefaultFilePerm = 0644
)

// DIR BUILDER

func (u *Utils) BuildDownloadsDir() string {
	return filepath.Join(u.cfg.StoragePath, downloadsDir)
}

func (u *Utils) BuildAudioDir(videoTitle string) string {
	return filepath.Join(u.cfg.StoragePath, videoTitle, audioDir)
}

func (u *Utils) BuildTranscriptDir(videoTitle string) string {
	return filepath.Join(u.cfg.StoragePath, videoTitle, transcriptDir)
}

func (u *Utils) BuildFaceTrackerDir(videoTitle string) string {
	return filepath.Join(u.cfg.StoragePath, videoTitle, faceTrackerDir)
}

func (u *Utils) BuildPromptDir(videoTitle string) string {
	return filepath.Join(u.cfg.StoragePath, videoTitle)
}

func (u *Utils) BuildMomentDir(videoTitle string) string {
	return filepath.Join(u.cfg.StoragePath, videoTitle, momentDir)
}

func (u *Utils) BuildRemotionPublicDir(videoTitle string) string {
	return filepath.Join(u.cfg.RemotionPublicPath, videoTitle)
}

func (u *Utils) BuildMetadataDir(videoTitle string) string {
	return filepath.Join(u.cfg.StoragePath, videoTitle, metadataDir)
}

func (u *Utils) BuildResultDir(videoTitle string) string {
	return filepath.Join(u.cfg.StoragePath, videoTitle, resultDir)
}

// BINARIES

func (u *Utils) GetYTDLP() string {
	return u.cfg.YTDLP
}

func (u *Utils) GetFFMPEG() string {
	return u.cfg.FFMPEG
}

func (u *Utils) GetPyEXETranscript() string {
	return u.cfg.PYEXETranscript
}

func (u *Utils) GetPyTranscript() string {
	return u.cfg.PYTranscript
}

func (u *Utils) GetPyEXEFaceTracker() string {
	return u.cfg.PYEXEFaceTracker
}

func (u *Utils) GetPyFaceTracker() string {
	return u.cfg.PYFaceTracker
}

func (u *Utils) GetPyEXEGenerateHook() string {
	return u.cfg.PYEXEGenerateHook
}

func (u *Utils) GetPyGenerateHook() string {
	return u.cfg.PYGenerateHook
}

// FILE SYSTEM

func (u *Utils) MkdirAll(path string) error {
	return os.MkdirAll(path, DefaultDirPerm)
}

func (u *Utils) WriteFile(path string, fileBytes []byte) error {
	return os.WriteFile(path, fileBytes, DefaultFilePerm)
}

func (u *Utils) ReadAndParse(path string, v any) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(fileBytes, v)
}

// HELPERS

func (u *Utils) NormalizeTitle(title string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	cleanTitle := reg.ReplaceAllString(title, "")

	cleanTitle = strings.TrimSpace(cleanTitle)

	return strings.ReplaceAll(cleanTitle, " ", "_")
}