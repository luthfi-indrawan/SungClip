# Panduan Instalasi SungClip

Panduan lengkap untuk menginstal dan menyiapkan SungClip di sistem Anda.

---

## Daftar Isi

- [Prasyarat](#prasyarat)
- [Instalasi Cepat (Recommended)](#instalasi-cepat-recommended)
- [Instalasi Manual](#instalasi-manual)
- [Verifikasi Instalasi](#verifikasi-instalasi)
- [Troubleshooting](#troubleshooting)

---

## Prasyarat

Sebelum menginstal SungClip, pastikan sistem Anda memiliki:

| Software    | Versi Minimum | Download                            |
| ----------- | ------------- | ----------------------------------- |
| **Node.js** | v18.x         | [nodejs.org](https://nodejs.org/)   |
| **Go**      | v1.21         | [go.dev](https://go.dev/dl/)        |
| **Python**  | v3.9          | [python.org](https://python.org/)   |
| **Git**     | --            | [git-scm.com](https://git-scm.com/) |

> **Catatan:** FFmpeg dan yt-dlp akan di-download otomatis oleh script setup.

### Verifikasi Prasyarat

Buka PowerShell/Terminal dan jalankan:

```powershell
node --version    # Harus keluar: v18.x.x atau lebih tinggi
go version        # Harus keluar: go version go1.21.x ...
python --version  # Harus keluar: Python 3.9.x atau lebih tinggi
```

---

## Instalasi Cepat (Recommended)

### Windows (PowerShell)

```powershell
# 1. Clone repository
git clone <repository-url>
cd SungClip

# 2. Jalankan setup otomatis
.\setup.ps1

# 3. Selesai! SungClip siap digunakan
.\sungclip.exe --help
```

### Windows (Batch - Fallback)

Jika PowerShell tidak tersedia atau ada masalah execution policy:

```cmd
# 1. Clone repository
git clone <repository-url>
cd SungClip

# 2. Jalankan setup batch
.\setup.bat
```

### Linux / macOS

```bash
# 1. Clone repository
git clone <repository-url>
cd SungClip

# 2. Jalankan setup
chmod +x setup.sh
./setup.sh

# 3. Selesai!
./sungclip --help
```

### Apa yang Dilakukan Script Setup?

Script setup akan otomatis:

1. **Mengecek prasyarat** -- Node.js, Go, Python
2. **Membuat folder `bin/`** -- Jika belum ada
3. **Download FFmpeg** -- Dari [gyan.dev](https://www.gyan.dev/ffmpeg/builds/) (Windows) atau cek sistem (Linux/Mac)
4. **Download yt-dlp** -- Dari [GitHub releases](https://github.com/yt-dlp/yt-dlp/releases)
5. **Setup Python virtual environment** -- Di `scripts/transcript/.venv/`
6. **Install Python dependencies** -- Dari `requirements.txt`
7. **Install Node.js dependencies** -- `npm install` di folder `remotion/`
8. **Setup Go modules** -- `go mod tidy`
9. **Build binary** -- Menghasilkan `sungclip.exe` (Windows) atau `sungclip` (Linux/Mac)

---

## Instalasi Manual

Jika Anda ingin menginstal secara manual atau script otomatis gagal:

### 1. Clone Repository

```powershell
git clone <repository-url>
cd SungClip
```

### 2. Setup Binary Dependencies

#### FFmpeg

**Windows:**

```powershell
# Download dari https://www.gyan.dev/ffmpeg/builds/
# Extract dan copy ffmpeg.exe ke folder bin/
```

**Linux:**

```bash
sudo apt update
sudo apt install ffmpeg
```

**macOS:**

```bash
brew install ffmpeg
```

#### yt-dlp

**Windows:**

```powershell
# Download dari https://github.com/yt-dlp/yt-dlp/releases
# Copy yt-dlp.exe ke folder bin/
```

**Linux/macOS:**

```bash
# Akan di-handle oleh setup script
# Atau install via package manager
```

### 3. Setup Python Environment

```powershell
cd scripts\transcript

# Buat virtual environment
python -m venv .venv

# Aktifkan virtual environment
.venv\Scripts\Activate.ps1    # Windows PowerShell
.venv\Scripts\activate.bat    # Windows CMD
source .venv/bin/activate     # Linux/Mac

# Install dependencies
pip install -r requirements.txt

cd ..\..
```

### 4. Setup Remotion (Node.js)

```powershell
cd remotion
npm install
cd ..
```

### 5. Setup Go & Build

```powershell
# Download dependencies
go mod tidy

# Build aplikasi
go build -o sungclip.exe .    # Windows
go build -o sungclip .        # Linux/Mac
```

---

## Verifikasi Instalasi

Setelah instalasi selesai, verifikasi dengan:

```powershell
# 1. Cek binary SungClip
.\sungclip.exe --help

# 2. Cek FFmpeg
.\bin\ffmpeg.exe -version

# 3. Cek yt-dlp
.\bin\yt-dlp.exe --version

# 4. Cek Python venv
.\scripts\transcript\.venv\Scripts\python.exe --version
```

Output yang diharapkan:

```
SungClip - AI-powered short video clip generator

Usage:
  sungclip [command]

Available Commands:
  ingest      Ingest a video from URL for processing
  edit        Generate edited video clips from ingested content
  help        Help about any command

Flags:
  -h, --help   help for sungclip
```

---

## Troubleshooting

### PowerShell Execution Policy

**Error:** `cannot be loaded because running scripts is disabled on this system`

**Solusi:**

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Python tidak ditemukan

**Error:** `Python is not installed`

**Solusi:**

- Pastikan Python terinstall dan ada di PATH
- Coba gunakan `py` sebagai alternatif: `py --version`
- Jika menggunakan Microsoft Store Python, nonaktifkan alias di Settings > Apps > Advanced app settings (Windows 11)

### FFmpeg download gagal

**Error:** `Failed to download ffmpeg`

**Solusi:**

1. Download manual dari [ffmpeg.org](https://ffmpeg.org/download.html)
2. Extract dan copy `ffmpeg.exe` ke folder `bin/`

### Node modules corrupt

**Error:** `Cannot find module` atau error saat `npm install`

**Solusi:**

```powershell
cd remotion
Remove-Item -Recurse -Force node_modules
Remove-Item package-lock.json
npm cache clean --force
npm install
```

### Go build error

**Error:** `go: module SungClip: not found`

**Solusi:**

```powershell
# Pastikan go.mod sudah ada
go mod init SungClip   # Jika belum ada
go mod tidy
go build -o sungclip.exe .
```

### yt-dlp tidak bisa download video

**Error:** `ERROR: [youtube] ... Sign in to confirm you're not a bot`

**Solusi:**

- Update yt-dlp: `.\bin\yt-dlp.exe -U`
- Gunakan cookies: `.\bin\yt-dlp.exe --cookies-from-browser chrome <url>`

---

## Butuh Bantuan?

Jika mengalami masalah yang tidak tercantum di atas:

1. Cek [USAGE.md](USAGE.md) untuk panduan penggunaan
2. Buka issue di repository GitHub
3. Cek log error dengan flag verbose (jika tersedia)

---

<div align="center">
  <sub>Selamat membuat clip!</sub>
</div>
