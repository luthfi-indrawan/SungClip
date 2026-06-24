import sys
import json
import cv2

from ultralytics import YOLO

# Ambil sampel setiap N frame
SAMPLE_EVERY = 3

# YOLO model
MODEL_NAME = "yolov8n.pt"


def moving_average(values, window=5):
    result = []

    for i in range(len(values)):
        start = max(0, i - window)
        end = min(len(values), i + window + 1)

        avg = sum(values[start:end]) / len(values[start:end])

        result.append(avg)

    return result


def main():
    if len(sys.argv) != 3:
        print("Usage:")
        print("python main.py <video_path> <output_json>")
        sys.exit(1)

    video_path = sys.argv[1]
    output_path = sys.argv[2]

    print("Loading model...")

    model = YOLO(MODEL_NAME)

    cap = cv2.VideoCapture(video_path)

    if not cap.isOpened():
        raise Exception(f"Cannot open video: {video_path}")

    fps = cap.get(cv2.CAP_PROP_FPS)
    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
    total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))

    print(f"Resolution : {width}x{height}")
    print(f"FPS        : {fps}")
    print(f"Frames     : {total_frames}")

    frames_data = []

    frame_idx = 0

    while True:
        success, frame = cap.read()

        if not success:
            break

        if frame_idx % SAMPLE_EVERY != 0:
            frame_idx += 1
            continue

        results = model.track(
            frame,
            persist=True,
            verbose=False,
            classes=[0],  # person
        )

        best_person = None
        best_area = 0

        for result in results:
            boxes = result.boxes

            if boxes is None:
                continue

            for box in boxes:
                x1, y1, x2, y2 = box.xyxy[0].tolist()

                w = x2 - x1
                h = y2 - y1

                area = w * h

                if area > best_area:
                    best_area = area

                    best_person = {
                        "x1": x1,
                        "y1": y1,
                        "x2": x2,
                        "y2": y2,
                        "width": w,
                        "height": h,
                    }

        if best_person:
            center_x = (
                best_person["x1"] +
                best_person["x2"]
            ) / 2

            center_y = (
                best_person["y1"] +
                best_person["y2"]
            ) / 2

            frames_data.append({
                "frame": frame_idx,
                "timeMs": int(frame_idx / fps * 1000),
                "centerX": center_x,
                "centerY": center_y,
                "width": best_person["width"],
                "height": best_person["height"],
            })

        frame_idx += 1

        if frame_idx % 300 == 0:
            print(
                f"Processed {frame_idx}/{total_frames}"
            )

    cap.release()

    print("Smoothing...")

    if len(frames_data) > 0:
        xs = [f["centerX"] for f in frames_data]
        ys = [f["centerY"] for f in frames_data]

        smooth_x = moving_average(xs, 5)
        smooth_y = moving_average(ys, 5)

        for i in range(len(frames_data)):
            frames_data[i]["centerX"] = round(
                smooth_x[i],
                2
            )

            frames_data[i]["centerY"] = round(
                smooth_y[i],
                2
            )

            frames_data[i]["width"] = round(
                frames_data[i]["width"],
                2
            )

            frames_data[i]["height"] = round(
                frames_data[i]["height"],
                2
            )

    output = {
        "videoWidth": width,
        "videoHeight": height,
        "fps": fps,
        "totalFrames": total_frames,
        "sampleEvery": SAMPLE_EVERY,
        "frames": frames_data,
    }

    with open(
        output_path,
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
    print(f"Output: {output_path}")
    print(f"Frames tracked: {len(frames_data)}")


if __name__ == "__main__":
    main()