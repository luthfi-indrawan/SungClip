import {
  AbsoluteFill,
  OffthreadVideo,
  Sequence,
  staticFile,
  interpolate,
  useCurrentFrame,
} from "remotion";

import { ClipMetadata } from "./types";
import { Subtitle } from "./components/Subtitle";

export const ClipComposition: React.FC<ClipMetadata> = ({
  video_path,
  subtitle,
  word_highlights,

  frames_face_trackers,

  ori_width,
  ori_height,

  target_width,
  target_height,
}) => {
  const frame = useCurrentFrame();

  const curtainProgress = interpolate(frame, [0, 40], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  const curtainHeight = ((1 - curtainProgress) * target_height) / 2;

  const scale = target_height / ori_height;

  const scaledWidth = ori_width * scale;

  const subtitleWidth = 980;

  const safeZoneEnd = target_width - subtitleWidth;

  let translateX = 0;

  if (frames_face_trackers.length > 0) {
    let maxFaceRight = 0;

    for (const tracker of frames_face_trackers) {
      const scaledCenterX = tracker.centerX * scale;

      const scaledPersonWidth = tracker.width * scale;

      const faceRight = scaledCenterX + scaledPersonWidth / 2 + 120;

      maxFaceRight = Math.max(maxFaceRight, faceRight);
    }

    const overflow = maxFaceRight - safeZoneEnd;

    if (overflow > 0) {
      translateX = -overflow;
    }
  }

  const maxTranslate = scaledWidth - target_width;

  translateX = Math.max(-maxTranslate, Math.min(0, translateX));

  return (
    <AbsoluteFill
      style={{
        overflow: "hidden",
      }}
    >
      <Sequence from={0}>
        <OffthreadVideo
          src={staticFile(video_path)}
          style={{
            position: "absolute",

            width: scaledWidth,
            height: target_height,

            transform: `translateX(${translateX}px)`,

            objectFit: "cover",
          }}
        />
      </Sequence>

      <Sequence from={0}>
        <Subtitle words={subtitle} highlights={word_highlights} />
      </Sequence>
      <div
        style={{
          position: "absolute",

          top: 0,
          left: 0,
          right: 0,

          height: curtainHeight,

          background: "#000",

          zIndex: 9999,
        }}
      />

      <div
        style={{
          position: "absolute",

          bottom: 0,
          left: 0,
          right: 0,

          height: curtainHeight,

          background: "#000",

          zIndex: 9999,
        }}
      />
    </AbsoluteFill>
  );
};
