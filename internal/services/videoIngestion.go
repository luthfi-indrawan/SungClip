package services

import (
	"SungClip/internal/types"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (s *services) DownloadVideo(url string, outputDir string) (videoPath string, err error) {
	outputTemplate := filepath.Join(outputDir, "%(title)s.%(ext)s")

	cmd := exec.Command(
		s.utils.GetYTDLP(),
		"-f", "bv*+ba/b",
		"--restrict-filenames",
		"-o", outputTemplate,
		"--print", "after_move:filepath",
		url,
	)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// ambil videoPath nya
	videoPath = strings.TrimSpace(stdout.String())

	// validasi output
	if videoPath == "" {
		return "", fmt.Errorf("yt-dlp did not return output file path")
	}
	
	if _, err := os.Stat(videoPath); err != nil {
		return "", err
	}

	return videoPath, nil
}

func (s *services) ExtractAudio(videoPath string, outputDir string) error {
	cmd := exec.Command(
		s.utils.GetFFMEPG(),
		"-y",
		"-i", videoPath,
		"-ar", "16000",
		"-ac", "1",
		"-c:a", "pcm_s16le",
		outputDir,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	if _, err := os.Stat(outputDir); err != nil {
		return err
	}

	return nil
}

func (s *services) Transcribe(audioPath string, outputPath string) error {
	cmd := exec.Command(
		s.utils.GetPyEXETranscript(),
		s.utils.GetPyTranscript(),
		audioPath,
		outputPath,
	)

	cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return err
    }

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("transcript file not found: %w", err)
	}

    return nil
}

func (s *services) BuildPrompt(transcript types.TranscriptResult) string {
	var b strings.Builder

	b.WriteString(`Kamu adalah content strategist profesional untuk TikTok, Reels, dan YouTube Shorts.

Tugasmu adalah mencari 5 bagian terbaik dari transcript berikut.

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
- Durasi ideal 60-120 detik.

Untuk setiap clip:

1. Pilih satu rentang transcript.
2. Semua field (title, caption, word_highlights) HARUS dibuat hanya dari isi rentang transcript tersebut.
3. Dilarang menggunakan informasi dari bagian transcript lain.
4. Judul harus dapat dibuktikan oleh isi clip.
5. Caption harus merangkum isi clip.
6. Jika title atau caption tidak sesuai dengan isi clip, clip dianggap gagal.
7. Headline harus berisi satu kalimat yang menarik audience, dan sesuai dengan isi clip.
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
    "hastags": [
      string
    ]
  }
]

Transcript:
`)

	for _, segment := range transcript {
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

func (s *services) formatDuration(ms int64) string {
	totalSeconds := ms / 1000

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	milliseconds := ms % 1000

	return fmt.Sprintf(
		"%02d:%02d:%02d.%03d",
		hours,
		minutes,
		seconds,
		milliseconds,
	)
}