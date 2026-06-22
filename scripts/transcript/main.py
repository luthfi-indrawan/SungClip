import json
import os
import sys
import time
import traceback

from faster_whisper import WhisperModel


def log_json(data: dict):
    """Print JSON ke stdout (Go bisa parse) + tetap readable."""
    print(json.dumps(data, ensure_ascii=False))
    sys.stdout.flush()


def main():
    if len(sys.argv) < 3:
        log_json({
            "error": "usage: main.py <audio_path> <output_path>"
        })
        sys.exit(1)

    audio_path = sys.argv[1]
    output_path = sys.argv[2]

    # Validasi audio exist
    if not os.path.exists(audio_path):
        log_json({"error": f"audio file not found: {audio_path}"})
        sys.exit(1)

    # Buat directory output kalau belum ada
    output_dir = os.path.dirname(output_path)
    if output_dir:
        os.makedirs(output_dir, exist_ok=True)

    try:
        log_json({"status": "loading_model"})

        started_at = time.time()

        model = WhisperModel(
            "small",
            device="cpu",
            compute_type="int8",
        )

        log_json({
            "status": "model_loaded",
            "elapsed_ms": int((time.time() - started_at) * 1000)
        })

        log_json({
            "status": "transcribing",
            "audio": audio_path,
            "output": output_path,
        })

        transcribe_start = time.time()

        segments, info = model.transcribe(
            audio_path,
            language="id",
            word_timestamps=True,
        )

        result = []
        count = 0

        for segment in segments:
            count += 1
            words = []

            if segment.words:
                for word in segment.words:
                    words.append({
                        "word": word.word.strip(),
                        "start_ms": int(word.start * 1000),
                        "end_ms": int(word.end * 1000),
                    })

            result.append({
                "start_ms": int(segment.start * 1000),
                "end_ms": int(segment.end * 1000),
                "text": segment.text.strip(),
                "words": words,
            })

        elapsed = time.time() - transcribe_start

        # Save transcript
        with open(output_path, "w", encoding="utf-8") as f:
            json.dump(result, f, ensure_ascii=False, indent=2)

        log_json({
            "status": "done",
            "output": output_path,
            "segments": count,
            "language": info.language,
            "language_probability": round(info.language_probability, 4),
            "elapsed_seconds": round(elapsed, 2),
        })

    except Exception as e:
        log_json({
            "error": str(e),
            "traceback": traceback.format_exc() if os.getenv("DEBUG") else None,
        })
        sys.exit(1)


if __name__ == "__main__":
    main()