package utils

import (
	"SungClip/internal/config"
	"os"
	"path/filepath"
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
	downloadsDir = "downloads"
	audioDir = "audio"
	transcriptDir = "transcript"
	momentDir = "moment"
	metadataDir = "metadata"
	resultDir = "result"
	faceTrackerDir = "face-tracker"
)

// dir builder
func (u *Utils) BuildDownloadsDir() string {
	return filepath.Join(u.cfg.StoragePath, downloadsDir)
}

func (u *Utils) BuildAudioDir(name string) string {
	return filepath.Join(u.cfg.StoragePath, name, audioDir)
}

func (u *Utils) BuildTranscriptDir(name string) string {
	return filepath.Join(u.cfg.StoragePath, name, transcriptDir)
}

func (u *Utils) BuildFaceTrackerDir(name string) string {
	return filepath.Join(u.cfg.StoragePath, name, faceTrackerDir)
}

func (u *Utils) BuildPromptDir(name string) string {
	return filepath.Join(u.cfg.StoragePath, name)
}

func (u *Utils) BuildMomentDir(name string) string {
	return filepath.Join(u.cfg.StoragePath, name, momentDir)
}

func (u *Utils) BuildRemotionPublicDir(name string) string {
	return filepath.Join(u.cfg.RemotionPublicPath, name)
}

func (u *Utils) BuildMetadataDir(name string) string {
	return filepath.Join(u.cfg.StoragePath, name, metadataDir)
}

func (u *Utils) BuildResultDir(name string) string {
	return filepath.Join(u.cfg.StoragePath, name, resultDir)
}

// gether scripts
func (u *Utils) GetYTDLP() string {
	return u.cfg.YTDLP
}

func (u *Utils) GetFFMEPG() string {
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

// os utils
func (u *Utils) MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

func (u *Utils) WriteFile(path string, fileBytes []byte) error {
	return os.WriteFile(path, fileBytes, 0644)
}