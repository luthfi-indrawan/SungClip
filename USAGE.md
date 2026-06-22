# Panduan Penggunaan SungClip

Panduan lengkap cara menggunakan SungClip untuk membuat video clip dari URL YouTube.

---

## Daftar Isi

- [Perintah Dasar](#perintah-dasar)
- [Alur Kerja Lengkap](#alur-kerja-lengkap)
- [Perintah `ingest`](#perintah-ingest)
- [Perintah `edit`](#perintah-edit)
- [Contoh Penggunaan](#contoh-penggunaan)
- [Konfigurasi](#konfigurasi)
- [Tips & Best Practices](#tips--best-practices)

---

## Perintah Dasar

SungClip memiliki 2 perintah utama:

| Perintah | Fungsi                                                 | Status         |
| -------- | ------------------------------------------------------ | -------------- |
| `ingest` | Download video, ekstrak audio, transcribe, analisis AI | Wajib pertama  |
| `edit`   | Generate video clip dengan subtitle & efek             | Setelah ingest |

### Bantuan Umum

```powershell
.\sungclip.exe --help
```

### Bantuan Per Perintah

```powershell
.\sungclip.exe ingest --help
.\sungclip.exe edit --help
```

---

## Alur Kerja Lengkap

```
+-------------------------------------------------------------+
|  LANGKAH 1: INGEST                                          |
|  .\sungclip.exe ingest --url "<youtube-url>"               |
+-------------------------------------------------------------+
|  1. Download video dari YouTube (yt-dlp)                    |
|  2. Ekstrak audio dari video (FFmpeg)                       |
|  3. Transcribe audio ke JSON (Python)                       |
|  4. Generate prompt untuk AI                                |
|  5. AI analisis & kembalikan JSON dengan timestamp clip      |
+-------------------------------------------------------------+
|  Output: Composition ID, Title, Prompt Path                 |
+-------------------------------------------------------------+
                              |
                              v
+-------------------------------------------------------------+
|  LANGKAH 2: EDIT                                            |
|  .\sungclip.exe edit --title "Judul" --comp <composition-id> |
+-------------------------------------------------------------+
|  1. Baca JSON hasil analisis AI                             |
|  2. Potong video sesuai timestamp (FFmpeg)                  |
|  3. Render dengan Remotion:                                 |
|     * Subtitle per kata (word-by-word)                      |
|     * Word highlight effect                                 |
|     * Crop ratio (9:16, 1:1, dll)                           |
|     * Face tracking                                         |
+-------------------------------------------------------------+
|  Output: Video clip siap upload!                            |
+-------------------------------------------------------------+
```

---

## Perintah `ingest`

Mengunduh dan memproses video dari URL YouTube.

### Sintaks

```powershell
.\sungclip.exe ingest --url "<youtube-url>"
```

### Flags

| Flag    | Singkat | Wajib | Default | Deskripsi                            |
| ------- | ------- | ----- | ------- | ------------------------------------ |
| `--url` | `-u`    | Ya    | --      | URL video YouTube yang akan diingest |

### Contoh

```powershell
# Ingest video YouTube
.\sungclip.exe ingest --url "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

# Atau dengan flag singkat
.\sungclip.exe ingest -u "https://youtu.be/dQw4w9WgXcQ"
```

### Output yang Diharapkan

```
Ingestion successful!
   Title:      How to Build a YouTube Channel
   PromptPath: storage/ingest/abc123/prompt.txt
```

### File yang Dihasilkan

```
storage/
└── ingest/
    └── {composition-id}/
        ├── video.mp4          # Video asli (downloaded)
        ├── audio.wav          # Audio hasil ekstrak
        ├── transcript.json    # Hasil transkripsi
        ├── prompt.txt         # Prompt untuk AI
        └── ai-response.json # Hasil analisis AI
```

---

## Perintah `edit`

Membuat video clip dari hasil ingest.

### Sintaks

```powershell
.\sungclip.exe edit --title "<judul>" --comp "<composition-id>" [flags]
```

### Flags

| Flag       | Singkat | Wajib | Default | Deskripsi                              |
| ---------- | ------- | ----- | ------- | -------------------------------------- |
| `--title`  | `-t`    | Ya    | --      | Judul untuk video clip yang dihasilkan |
| `--comp`   | `-c`    | Ya    | --      | Composition ID dari hasil ingest       |
| `--width`  | `-W`    | Tidak | `1080`  | Lebar video dalam pixel                |
| `--height` | `-H`    | Tidak | `1920`  | Tinggi video dalam pixel               |

### Contoh

```powershell
# Edit dengan setting default (1080x1920 - portrait/shorts)
.\sungclip.exe edit --title "Viral Clip #1" --comp abc123

# Edit dengan custom resolution (landscape)
.\sungclip.exe edit -t "Best Moments" -c abc123 -W 1920 -H 1080

# Edit dengan ratio 1:1 (square)
.\sungclip.exe edit -t "Instagram Clip" -c abc123 -W 1080 -H 1080
```

### Output yang Diharapkan

```
Editing successful!
   Title:           Viral Clip #1
   Total Clips:     5
   Result Path:     storage/edit/abc123/result/
```

### File yang Dihasilkan

```
storage/
└── edit/
    └── {composition-id}/
        └── result/
            ├── clip_01.mp4    # Clip 1 dengan subtitle & efek
            ├── clip_02.mp4    # Clip 2 dengan subtitle & efek
            ├── ...
            └── clip_05.mp4    # Clip terakhir
```

---

## Contoh Penggunaan Lengkap

### Skenario 1: Basic Workflow

```powershell
# Step 1: Ingest video
.\sungclip.exe ingest -u "https://www.youtube.com/watch?v=VIDEO_ID"

# Output: Composition ID = "comp_20240622_abc123"

# Step 2: Edit & generate clips
.\sungclip.exe edit -t "Best Moments" -c comp_20240622_abc123
```

### Skenario 2: Multiple Resolutions

```powershell
# Ingest sekali
.\sungclip.exe ingest -u "https://youtube.com/watch?v=VIDEO_ID"
# ID: comp_xyz789

# Generate untuk Shorts (9:16)
.\sungclip.exe edit -t "Shorts Version" -c comp_xyz789 -W 1080 -H 1920

# Generate untuk Reels (9:16)
.\sungclip.exe edit -t "Reels Version" -c comp_xyz789 -W 1080 -H 1920

# Generate untuk Feed (1:1)
.\sungclip.exe edit -t "Feed Version" -c comp_xyz789 -W 1080 -H 1080
```

### Skenario 3: Batch Processing (PowerShell)

```powershell
# Daftar URL
$urls = @(
    "https://youtube.com/watch?v=VIDEO1",
    "https://youtube.com/watch?v=VIDEO2",
    "https://youtube.com/watch?v=VIDEO3"
)

# Ingest semua
foreach ($url in $urls) {
    .\sungclip.exe ingest -u $url
}

# Edit semua (sesuaikan composition ID dari output di atas)
.\sungclip.exe edit -t "Clip Batch 1" -c <comp-id-1>
.\sungclip.exe edit -t "Clip Batch 2" -c <comp-id-2>
.\sungclip.exe edit -t "Clip Batch 3" -c <comp-id-3>
```

---

## Konfigurasi

### Environment Variables (Opsional)

Buat file `.env` di root folder (jika didukung oleh aplikasi):

```env
# AI API Configuration
OPENAI_API_KEY=your_api_key_here
# atau
ANTHROPIC_API_KEY=your_api_key_here

# Output Settings
OUTPUT_DIR=./storage/output
TEMP_DIR=./storage/temp

# Video Settings
DEFAULT_WIDTH=1080
DEFAULT_HEIGHT=1920
DEFAULT_FPS=30
```

### Konfigurasi Subtitle

Subtitle style dapat dikonfigurasi di file Remotion (lihat folder `remotion/src/`):

| Properti          | Deskripsi                  | Default           |
| ----------------- | -------------------------- | ----------------- |
| `fontFamily`      | Font subtitle              | `Inter`           |
| `fontSize`        | Ukuran font                | `48`              |
| `fontColor`       | Warna teks                 | `#FFFFFF`         |
| `highlightColor`  | Warna highlight kata aktif | `#FFD700`         |
| `backgroundColor` | Warna background subtitle  | `rgba(0,0,0,0.6)` |
| `position`        | Posisi subtitle            | `bottom`          |

---

## Tips & Best Practices

### 1. Pilih Video yang Tepat

- Video dengan durasi 10-60 menit ideal
- Konten dengan banyak momen "quotable" atau viral
- Audio yang jelas (transkripsi lebih akurat)
- Hindari video dengan banyak noise audio
- Hindari video dengan copyright music (bisa kena strike)

### 2. Optimasi AI Prompt

Hasil analisis AI sangat bergantung pada kualitas transkripsi. Pastikan:

- Audio jelas dan tidak ada background noise berlebihan
- Speaker berbicara dengan jelas
- Untuk video multi-speaker, pertimbangkan manual review

### 3. Resolution Guide

| Platform        | Ratio | Resolution | Command           |
| --------------- | ----- | ---------- | ----------------- |
| YouTube Shorts  | 9:16  | 1080x1920  | `-W 1080 -H 1920` |
| Instagram Reels | 9:16  | 1080x1920  | `-W 1080 -H 1920` |
| TikTok          | 9:16  | 1080x1920  | `-W 1080 -H 1920` |
| Instagram Feed  | 1:1   | 1080x1080  | `-W 1080 -H 1080` |
| Twitter/X       | 16:9  | 1920x1080  | `-W 1920 -H 1080` |

### 4. Manajemen Storage

Folder `storage/` akan membesar seiring penggunaan. Tips:

- Hapus folder `storage/ingest/` yang sudah tidak dibutuhkan
- Backup folder `storage/edit/` yang berisi hasil jadi
- Pertimbangkan external storage untuk video besar

### 5. Performance

- **Ingest**: Bergantung pada durasi video dan koneksi internet
- **Transcribe**: Bergantung pada durasi audio (sekitar real-time)
- **AI Analysis**: Bergantung pada provider AI (1-30 detik)
- **Edit/Render**: Bergantung pada jumlah clip dan durasi (2-5x real-time)

---

## Troubleshooting Penggunaan

### Ingest gagal: "Video unavailable"

- Pastikan URL valid dan video tidak private/region-locked
- Coba update yt-dlp: `.\bin\yt-dlp.exe -U`
- Untuk video age-restricted, gunakan cookies browser

### Edit gagal: "Composition not found"

- Pastikan composition ID benar (copy dari output ingest)
- Cek folder `storage/ingest/{comp-id}/` ada dan lengkap

### Render lambat

- Tutup aplikasi berat lainnya
- Pastikan GPU tersedia untuk Remotion rendering
- Kurangi jumlah clip sekaligus

### Subtitle tidak muncul

- Cek file `transcript.json` valid
- Pastikan Remotion dependencies terinstall dengan benar

---

## Lihat Juga

- [INSTALLATION.md](INSTALLATION.md) -- Panduan instalasi
- [README.md](README.md) -- Overview proyek

---

<div align="center">
  <sub>Happy clipping!</sub>
</div>
