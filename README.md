# SungClip

> Turn long-form YouTube videos into viral-ready short clips — automatically, fast, and free.

SungClip is a CLI automation tool that transforms long YouTube videos (podcasts, interviews, livestreams, webinars) into multiple short-form videos (Shorts, Reels, TikToks) using AI-powered speech-to-text, intelligent content analysis, and automated video rendering with word-by-word subtitles.

**Built by [Luthfi Indrawan](https://github.com/luthfiindrawan).**

---

## 🎯 Why SungClip?

Video clipping is a high-income potential niche, but the process is painfully slow. A 30-minute to 1-hour video can take **1–3 hours** to manually clip, analyze, and edit.

SungClip solves this by automating the entire pipeline — from download to render — so you can focus on growing your channel, not wrestling with timelines.

### Key Goals

- ⚡ **Automate** the entire clip creation workflow
- 🤖 **Accelerate** with AI analysis and smart algorithms
- 💰 **Zero cost** processing (self-hosted, open-source)
- 🚀 **Future-ready** architecture for campaign tracking & auto-upload

---

## ✨ Features (v0.1)

- **YouTube Ingestion** — Download videos directly from YouTube URLs
- **Audio Extraction** — Automatic audio extraction for transcription
- **AI Transcription** — Speech-to-text with word-level timestamp alignment
- **AI Content Analysis** — Smart detection of viral-worthy moments
- **Automated Cutting** — Precise segment cutting via FFmpeg
- **Automated Rendering** — Generate short clips via Remotion with:
  - Word-by-word animated subtitles
  - Auto crop to vertical ratio (9:16)
  - Face tracking support
- **Simple Configuration** — Sensible defaults with optional custom paths

---

## 🏗️ Architecture

SungClip uses a multi-language stack orchestrated by Go:

| Layer                     | Language   | Responsibility                                            |
| ------------------------- | ---------- | --------------------------------------------------------- |
| **CLI & Orchestrator**    | Go         | Command handling, workflow pipeline, service coordination |
| **Transcription Service** | Python     | Whisper-based speech-to-text + word alignment             |
| **Video Renderer**        | TypeScript | Remotion composition, subtitle animation, face tracking   |

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────┐
│   Go CLI    │────►│ Python Service  │────►│  Remotion   │
│ (sungclip)  │     │  (Whisper/ML)   │     │  (Render)   │
└──────┬──────┘     └─────────────────┘     └──────┬──────┘
       │                                          │
       └──────────────────────────────────────────┘
                          ▼
                    ┌─────────────┐
                    │   Storage   │
                    │ (Local/S3)  │
                    └─────────────┘
```

---

## 🔄 How It Works

```
YouTube URL ──► Download ──► Extract Audio ──► Transcribe (Whisper)
                                                           │
                                                           ▼
                              AI Analysis ◄─── Prompt + JSON ◄─── Word-level timestamps
                                   │
                                   ▼
                    Viral moments detected (valid JSON response)
                                   │
                                   ▼
                    FFmpeg Cutting ──► Remotion Rendering
                    (auto segments)      (subtitles + crop + face tracking)
                                                          │
                                                          ▼
                                                  Ready to Upload!
```

### Step-by-step Flow

1. **Send YouTube URL** → `sungclip ingest <url>`
2. **Download & Extract** → Video downloaded, audio extracted
3. **Transcribe** → Whisper generates JSON with word-level timestamps
4. **AI Analysis** → JSON converted to AI prompt, AI returns valid JSON with clip segments
5. **Cut** → FFmpeg cuts video precisely per AI-detected segments
6. **Render** → Remotion adds subtitles, vertical crop, word highlight, face tracking
7. **Done** → Final clips ready for upload

---

## 📦 Prerequisites

### Required

- [Go](https://go.dev/dl/) 1.22+
- [Node.js](https://nodejs.org/) 20+ & npm
- [Python](https://python.org/) 3.10+
- [FFmpeg](https://ffmpeg.org/download.html)

### Python Environment Setup

```bash
cd scripts/transcript
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

### Remotion Setup

```bash
cd remotion
npm install
```

---

## 🚀 Quick Start

### 1. Clone & Build

```bash
git clone https://github.com/luthfiindrawan/sungclip.git
cd sungclip
go build -o sungclip ./cmd/sungclip
```

### 2. Configure (Optional)

```bash
# Default storage path is ./storage
# To customize, set environment variable:
export SUNGCLIP_STORAGE_PATH=/path/to/your/storage
```

### 3. Run

```bash
# Step 1: Ingest a YouTube video
./sungclip ingest https://www.youtube.com/watch?v=VIDEO_ID

# Step 2: AI generates a prompt — copy the JSON output,
# save it to a file, and optionally edit clip selections

# Step 3: Render clips
./sungclip editing clips.json
# or
./sungclip editing /path/to/clips.json
```

---

## 🛠️ Commands

### `ingest <youtube-url>`

Download a YouTube video and prepare it for processing.

```bash
./sungclip ingest https://youtube.com/watch?v=...
```

**What happens:**

- Video downloaded to `storage/uploads/`
- Audio extracted for transcription
- Whisper transcribes with word-level timestamps
- JSON prompt generated and printed to stdout

**Output example:**

```json
{
  "source": "uploads/video_id.mp4",
  "transcript": "uploads/video_id.json",
  "prompt": "AI-ready prompt with transcript context..."
}
```

### `editing <input>`

Render short clips from a JSON configuration file.

```bash
./sungclip editing clips.json
./sungclip editing /absolute/path/to/clips.json
```

**Input JSON format:**

```json
{
  "source": "uploads/video_id.mp4",
  "clips": [
    {
      "start": 120.5,
      "end": 145.0,
      "title": "The Most Important Lesson",
      "hook": "Why most people fail at..."
    }
  ],
  "style": {
    "subtitle_color": "#FFFFFF",
    "highlight_color": "#FFD700",
    "background": "blur"
  }
}
```

**Output:**

- Rendered clips in `storage/outputs/`
- Vertical 9:16 format
- Word-by-word animated subtitles
- Face tracking applied (if enabled)

---

## ⚙️ Configuration

SungClip uses sensible defaults. Configuration is loaded from environment variables:

| Variable                  | Default                 | Description                           |
| ------------------------- | ----------------------- | ------------------------------------- |
| `SUNGCLIP_STORAGE_PATH`   | `./storage`             | Root path for all file operations     |
| `SUNGCLIP_TRANSCRIPT_URL` | `http://localhost:8000` | Python transcription service endpoint |
| `SUNGCLIP_REMOTION_PATH`  | `./remotion`            | Path to Remotion project directory    |
| `SUNGCLIP_AI_MODEL`       | `gpt-4o-mini`           | AI model for content analysis         |

---

## 📁 Project Structure

```
SungClip/
├── cmd/sungclip/          # CLI entry point (Go)
├── internal/              # Core orchestration logic (Go)
│   ├── pipeline.go        # Main workflow orchestrator
│   ├── transcript.go      # Python service client
│   ├── renderer.go        # Remotion integration
│   ├── storage.go         # File storage utilities
│   └── config.go          # Configuration management
├── scripts/transcript/    # Python transcription service
│   ├── main.py            # FastAPI/HTTP server entry
│   └── requirements.txt   # Python dependencies
├── remotion/              # Remotion video renderer (TypeScript)
│   ├── src/
│   │   ├── SungClip.tsx       # Main video composition
│   │   ├── WordHighlight.tsx  # Word-by-word subtitle component
│   │   └── types.ts           # Shared type definitions
│   └── package.json
├── storage/               # File storage (gitignored)
│   ├── uploads/           # Downloaded source videos
│   ├── temp/              # Processing scratch space
│   └── outputs/           # Final rendered clips
├── go.mod
├── Makefile
├── LICENSE
└── README.md
```

---

## 🔄 Development Workflow

```bash
# Run all services for development
make dev

# Build CLI binary
make build

# Run tests
make test

# Clean generated files
make clean
```

---

## 🚧 Current Limitations

SungClip v0.1 is focused on core automation. The following features are **not yet implemented** but the architecture is designed to support them:

- ❌ Campaign tracking & trend monitoring
- ❌ Auto-upload to YouTube Shorts / TikTok / Reels
- ❌ Multi-platform source (currently YouTube only)
- ❌ Web dashboard for visual editing

**The system is architected to scale toward full automation.** Stay tuned for big moves 🚀

---

## 🗺️ Roadmap

### v0.1 (Current)

- [x] YouTube video ingestion
- [x] Audio extraction
- [x] AI transcription with word-level timestamps
- [x] AI-powered content analysis for clip detection
- [x] JSON-based clip configuration
- [x] FFmpeg automated cutting
- [x] Remotion rendering with subtitles, crop, word highlight, face tracking

### v0.2 (Planned)

- [ ] Batch processing multiple videos
- [ ] Auto-clip detection without manual JSON editing
- [ ] Custom subtitle styling templates
- [ ] Background music & padding options
- [ ] Progress bars & better CLI UX

### v1.0 (Future)

- [ ] Campaign tracking & trend analysis
- [ ] Direct upload to YouTube Shorts / TikTok / Reels
- [ ] Web dashboard for visual clip editing
- [ ] Cloud storage integration (S3)
- [ ] Queue-based distributed processing
- [ ] Multi-language transcription support

---

## 🤝 Contributing

Contributions are welcome! SungClip is built for personal use but open to the community.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

[MIT](LICENSE) © 2026 Luthfi Indrawan

---

> **Note:** SungClip is in active development (v0.1). Built with passion for automation and content creation. Feedback, issues, and contributions are highly appreciated!
