package services

import (
	"SungClip/internal/types"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (s *services) DownloadVideo(
	ctx context.Context,
	url string,
	outputDir string,
) (videoPath, infoVideoPath string, err error) {

	outputTemplate := filepath.Join(outputDir, "%(title)s.%(ext)s")

	cmd := exec.CommandContext(
		ctx,
		s.utils.GetYTDLP(),
		"-f", "bv*+ba/b",
		"--restrict-filenames",
		"--write-info-json",
		"-o", outputTemplate,
		"--print", "after_move:filepath",
		url,
	)

	var stdout bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("yt-dlp failed: %w", err)
	}

	lines := strings.Split(
		strings.TrimSpace(stdout.String()),
		"\n",
	)

	videoPath = strings.TrimSpace(lines[len(lines)-1])

	if videoPath == "" {
		return "", "", fmt.Errorf("yt-dlp returned empty filepath")
	}

	if _, err := os.Stat(videoPath); err != nil {
		return "", "", fmt.Errorf("video file not found: %w", err)
	}

	infoVideoPath = strings.TrimSuffix(
		videoPath,
		filepath.Ext(videoPath),
	) + ".info.json"

	if _, err := os.Stat(infoVideoPath); err != nil {
		return "", "", fmt.Errorf(
			"metadata file not found: %s",
			infoVideoPath,
		)
	}

	return videoPath, infoVideoPath, nil
}

func (s *services) ExtractAudio(
	ctx context.Context,
	videoPath string,
	outputPath string,
) error {

	cmd := exec.CommandContext(
		ctx,
		s.utils.GetFFMPEG(),
		"-y",
		"-i", videoPath,
		"-ar", "16000",
		"-ac", "1",
		"-c:a", "pcm_s16le",
		outputPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg extract audio failed: %w", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("audio file not found: %w", err)
	}

	return nil
}

func (s *services) Transcribe(
	ctx context.Context,
	audioPath string,
	outputPath string,
) error {

	cmd := exec.CommandContext(
		ctx,
		s.utils.GetPyEXETranscript(),
		s.utils.GetPyTranscript(),
		audioPath,
		outputPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("transcribe failed: %w", err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf(
			"transcript file not found: %w",
			err,
		)
	}

	return nil
}

func (s *services) BuildPrompt(
	metadataVideo types.MetadataVideo,
	clipsCount, minDuration, maxDuration int,
) string {

	var b strings.Builder

	fmt.Fprintf(
		&b,
		`Kamu adalah content strategist profesional untuk TikTok, Reels, dan YouTube Shorts.

Tugasmu adalah mencari %d bagian terbaik dari transcript berikut.

Tujuan:
- menemukan clip yang menarik
- menentukan timestamp yang presisi
- menghindari intro yang terlalu panjang
- menghindari ending yang terlalu lambat

Kriteria:
1. Hook kuat
2. Insight bernilai tinggi
3. Cerita menarik
4. Kontroversi
5. Humor
6. Fakta mengejutkan
7. Potensi viral

PENTING:

- Gunakan timestamp yang diberikan.
- Pilih waktu mulai yang sedekat mungkin dengan awal hook.
- Pilih waktu selesai yang sedekat mungkin dengan akhir poin utama.
- Jangan memilih clip yang terlalu panjang.
- Durasi ideal %d-%d detik.

Untuk setiap clip:

1. Pilih satu rentang transcript.
2. Semua field (title, caption, word_highlights) HARUS dibuat hanya dari isi rentang transcript tersebut.
3. Dilarang menggunakan informasi dari bagian transcript lain.
4. Judul harus dapat dibuktikan oleh isi clip.
5. Caption harus merangkum isi clip.
6. Jika title atau caption tidak sesuai dengan isi clip, clip dianggap gagal.
7. Headline harus berisi satu kalimat yang menarik audience bahkan clipbait yang buat orang penasaran namun sesuai dengan isi clip.
8. Score 0-100

Output HARUS berupa JSON valid.

[
  {
    "start_ms": number,
    "end_ms": number,
    "score": number,
    "title": string,
    "headline": string,
    "caption": string,
    "word_highlights": [
      string
    ],
    "hashtags": [
      string
    ]
  }
]

Informasi video:
- title: %s
- channel: %s

Output HARUS berupa JSON valid.

`,
		clipsCount,
		minDuration,
		maxDuration,
		metadataVideo.Title,
		metadataVideo.Channel,
	)

	for _, segment := range metadataVideo.TranscriptResult {
		fmt.Fprintf(
			&b,
			"\n[%d - %d]\n%s\n",
			segment.StartMS,
			segment.EndMS,
			segment.Text,
		)
	}

	return b.String()
}