import argparse
import asyncio
import json

import edge_tts
from mutagen.mp3 import MP3

DEFAULT_VOICE = "id-ID-ArdiNeural"


async def generate(
    text: str,
    output_path: str,
    voice: str,
):
    communicate = edge_tts.Communicate(
        text=text,
        voice=voice
    )

    await communicate.save(output_path)

    audio = MP3(output_path)

    duration_ms = int(
        audio.info.length * 1000
    )

    return duration_ms


def main():
    parser = argparse.ArgumentParser()

    parser.add_argument(
        "--text",
        required=True,
    )

    parser.add_argument(
        "--output",
        required=True,
    )

    parser.add_argument(
        "--voice",
        default=DEFAULT_VOICE,
    )

    args = parser.parse_args()

    duration_ms = asyncio.run(
        generate(
            text=args.text,
            output_path=args.output,
            voice=args.voice,
        )
    )

    print(
        json.dumps(
            {
                "audio_path": args.output,
                "duration_ms": duration_ms,
            }
        )
    )


if __name__ == "__main__":
    main()