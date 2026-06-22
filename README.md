# SungClip

> **AI-powered short video clip generator with word-by-word subtitles**

SungClip adalah tools automation untuk membuat video clip pendek dari video YouTube panjang, lengkap dengan subtitle per kata, face tracking, dan crop ratio otomatis.

---

## Fitur Utama

| Fitur                      | Deskripsi                                                                       |
| -------------------------- | ------------------------------------------------------------------------------- |
| **Auto Download**          | Download video dari URL YouTube otomatis                                        |
| **Audio Extraction**       | Ekstrak audio dari video untuk proses transkripsi                               |
| **AI Transcription**       | Transkripsi otomatis dengan output JSON                                         |
| **AI Content Analysis**    | Analisis potensi konten bagus menggunakan AI                                    |
| **Auto Cutting**           | Potong video otomatis dengan FFmpeg berdasarkan rekomendasi AI                  |
| **Remotion Editing**       | Edit video dengan subtitle per kata, word highlight, crop ratio & face tracking |
| **Word-by-Word Subtitles** | Subtitle yang muncul kata per kata dengan efek highlight                        |
| **Free Cost Target**       | Didesain untuk meminimalkan biaya operasional                                   |

---

## Tech Stack

SungClip dibangun dengan 3 bahasa pemrograman:

| Bahasa         | Penggunaan                                        |
| -------------- | ------------------------------------------------- |
| **Go**         | Backend utama, CLI, orchestration, FFmpeg cutting |
| **TypeScript** | Remotion video templates & rendering              |
| **Python**     | Transkripsi audio & AI processing                 |

---

## Struktur Folder

```
SungClip/
├── bin/                    # Binary dependencies (ffmpeg.exe, yt-dlp.exe)
├── internal/               # Go source code
│   ├── config/            # Konfigurasi aplikasi
│   ├── controllers/       # HTTP/CLI handlers
│   ├── services/          # Business logic
│   ├── types/             # Type definitions
│   └── utils/             # Utility functions
├── remotion/              # Remotion video templates (TypeScript/Node.js)
├── scripts/
│   └── transcript/        # Python transcription scripts
│       └── .venv/         # Python virtual environment
├── storage/               # Temp storage untuk video & hasil
├── go.mod                 # Go module
├── go.sum                 # Go dependencies lock
├── main.go                # Entry point aplikasi
├── setup.ps1              # Setup script (Windows PowerShell)
├── setup.bat              # Setup script (Windows Batch)
├── setup.sh               # Setup script (Linux/Mac)
├── README.md              # Dokumentasi ini
├── INSTALLATION.md        # Panduan instalasi detail
└── USAGE.md               # Panduan penggunaan
```

---

## Quick Start

```powershell
# 1. Clone repository
git clone <repository-url>
cd SungClip

# 2. Jalankan setup (otomatis download binary, install dependencies, build)
.\setup.ps1

# 3. Ingest video dari YouTube
.\sungclip.exe ingest --url "https://youtube.com/watch?v=..."

# 4. Edit & generate clips
.\sungclip.exe edit --title "My Clip" --comp <composition-id>
```

---

## Alur Kerja

```
+-------------+    +--------------+    +-------------+
|  YouTube    |--->|   Download   |--->|   Audio     |
|    URL      |    |   (yt-dlp)   |    |  Extract    |
+-------------+    +--------------+    +------+------+
                                              |
+-------------+    +--------------+    +------v------+
|  Remotion   |<---| Video Cutting|<---| Transcribe  |
|   Render    |    |  (FFmpeg)    |    |  (Python)   |
+-------------+    +--------------+    +------+------+
                                              |
+-------------+    +--------------+    +------v------+
|  Subtitle + |<---|   Render     |<---|  AI Analysis|
| Face Track  |    |   Clips      |    |   (JSON)    |
+-------------+    +--------------+    +-------------+
```

---

## Prerequisites

- **Node.js** (v18+) -- untuk Remotion
- **Go** (v1.21+) -- untuk backend
- **Python** (v3.9+) -- untuk transkripsi
- **FFmpeg** -- akan di-download otomatis oleh setup script
- **yt-dlp** -- akan di-download otomatis oleh setup script

---

## Dokumentasi Lengkap

- **[INSTALLATION.md](INSTALLATION.md)** -- Panduan instalasi detail & troubleshooting
- **[USAGE.md](USAGE.md)** -- Panduan penggunaan lengkap dengan contoh

---

## Roadmap

- [x] Auto download & transcribe
- [x] AI content analysis
- [x] Auto cutting dengan FFmpeg
- [x] Remotion editing (subtitle, crop, face tracking)
- [ ] Campaign tracking otomatis
- [ ] Auto upload ke platform social media

---

## Lisensi

[MIT License](LICENSE)

---

<div align="center">
  <sub>Dibuat dengan untuk content creator clipper</sub>
</div>
