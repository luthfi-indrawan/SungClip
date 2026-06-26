import sys
import json
import cv2

from collections import defaultdict

import supervision as sv
from ultralytics import YOLO

MODEL_PATH = "models/yolov8n-face.pt"
SAMPLE_EVERY = 3


def moving_average(values, window=5):
    result = []

    for i in range(len(values)):
        start = max(0, i - window)
        end = min(len(values), i + window + 1)

        avg = sum(values[start:end]) / len(values[start:end])

        result.append(avg)

    return result


def smooth_track(track_frames):
    if len(track_frames) < 2:
        return

    xs = [f["centerX"] for f in track_frames]
    ys = [f["centerY"] for f in track_frames]
    ws = [f["width"] for f in track_frames]
    hs = [f["height"] for f in track_frames]

    smooth_x = moving_average(xs)
    smooth_y = moving_average(ys)
    smooth_w = moving_average(ws)
    smooth_h = moving_average(hs)

    for i in range(len(track_frames)):
        track_frames[i]["centerX"] = round(
            smooth_x[i],
            2
        )

        track_frames[i]["centerY"] = round(
            smooth_y[i],
            2
        )

        track_frames[i]["width"] = round(
            smooth_w[i],
            2
        )

        track_frames[i]["height"] = round(
            smooth_h[i],
            2
        )


def main():
    if len(sys.argv) != 3:
        print(
            "python main.py <video_path> <output_json>"
        )
        sys.exit(1)

    video_path = sys.argv[1]
    output_json = sys.argv[2]

    print("Loading model...")

    model = YOLO(MODEL_PATH)

    tracker = sv.ByteTrack()

    cap = cv2.VideoCapture(video_path)

    if not cap.isOpened():
        raise Exception(
            f"Cannot open video: {video_path}"
        )

    fps = cap.get(cv2.CAP_PROP_FPS)

    width = int(
        cap.get(cv2.CAP_PROP_FRAME_WIDTH)
    )

    height = int(
        cap.get(cv2.CAP_PROP_FRAME_HEIGHT)
    )

    total_frames = int(
        cap.get(cv2.CAP_PROP_FRAME_COUNT)
    )

    print(f"Resolution : {width}x{height}")
    print(f"FPS        : {fps}")
    print(f"Frames     : {total_frames}")

    frame_idx = 0

    frames = []

    track_meta = {}

    track_history = defaultdict(list)

    while True:
        success, frame = cap.read()

        if not success:
            break

        if frame_idx % SAMPLE_EVERY != 0:
            frame_idx += 1
            continue

        results = model.predict(
            frame,
            conf=0.30,
            verbose=False
        )

        detections = sv.Detections.from_ultralytics(
            results[0]
        )

        detections = tracker.update_with_detections(
            detections
        )

        frame_tracks = []

        current_time_ms = int(
            frame_idx / fps * 1000
        )

        if detections.tracker_id is not None:
            for i in range(len(detections)):
                track_id = int(
                    detections.tracker_id[i]
                )

                x1, y1, x2, y2 = (
                    detections.xyxy[i]
                )

                confidence = float(
                    detections.confidence[i]
                )

                width_box = float(
                    x2 - x1
                )

                height_box = float(
                    y2 - y1
                )

                center_x = float(
                    (x1 + x2) / 2
                )

                center_y = float(
                    (y1 + y2) / 2
                )

                track_entry = {
                    "trackId": track_id,
                    "centerX": center_x,
                    "centerY": center_y,
                    "width": width_box,
                    "height": height_box,
                    "confidence": round(
                        confidence,
                        4
                    )
                }

                frame_tracks.append(
                    track_entry
                )

                track_history[
                    track_id
                ].append(track_entry)

                if track_id not in track_meta:
                    track_meta[
                        track_id
                    ] = {
                        "trackId": track_id,
                        "firstSeenMs": current_time_ms,
                        "lastSeenMs": current_time_ms,
                    }
                else:
                    track_meta[
                        track_id
                    ]["lastSeenMs"] = (
                        current_time_ms
                    )

        frames.append({
            "frame": frame_idx,
            "timeMs": current_time_ms,
            "tracks": frame_tracks
        })

        frame_idx += 1

        if frame_idx % 300 == 0:
            print(
                f"Processed {frame_idx}/{total_frames}"
            )

    cap.release()

    print("Smoothing tracks...")

    for track_id in track_history:
        smooth_track(
            track_history[track_id]
        )

    output = {
        "videoWidth": width,
        "videoHeight": height,
        "fps": fps,
        "totalFrames": total_frames,
        "sampleEvery": SAMPLE_EVERY,
        "tracks": sorted(
            track_meta.values(),
            key=lambda x: x["trackId"]
        ),
        "frames": frames,
    }

    with open(
        output_json,
        "w",
        encoding="utf-8"
    ) as f:
        json.dump(
            output,
            f,
            indent=2,
            ensure_ascii=False,
        )

    print("")
    print("Done!")
    print(f"Output: {output_json}")
    print(
        f"Tracks found: {len(track_meta)}"
    )


if __name__ == "__main__":
    main()