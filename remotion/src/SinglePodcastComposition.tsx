import React from "react";
import {
  AbsoluteFill,
  Audio,
  OffthreadVideo,
  Sequence,
  staticFile,
  interpolate,
  useCurrentFrame,
} from "remotion";

import { ClipMetadata } from "./types";
import { Subtitle } from "./components/Subtitle";
import { HookScreen } from "./components/HookScreen";

export const SinglePodcastComposition: React.FC<ClipMetadata> = ({
  hook,
  headline,

  fps,

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

  const hookFrames = Math.ceil((hook.duration_ms / 1000) * fps);

  /**
   * Frame video utama
   */
  const videoFrame = Math.max(0, frame - hookFrames);

  /**
   * Cari frame terbaik untuk hook
   * berdasarkan total area wajah terbesar
   */
  let hookPreviewFrame = 0;
  let bestScore = 0;

  for (const trackerFrame of frames_face_trackers) {
    const score = trackerFrame.tracks.reduce(
      (sum, track) => sum + track.width * track.height,
      0,
    );

    if (score > bestScore) {
      bestScore = score;
      hookPreviewFrame = trackerFrame.frame;
    }
  }

  /**
   * Curtain Animation
   */
  const curtainProgress = interpolate(videoFrame, [0, 40], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  const curtainHeight = ((1 - curtainProgress) * target_height) / 2;

  /**
   * Face Tracking
   */
  const scale = target_height / ori_height;

  const scaledWidth = ori_width * scale;

  const subtitleWidth = 800;

  const safeZoneEnd = target_width - subtitleWidth;

  let translateX = 0;

  const WINDOW_SIZE = 2;
  const DEAD_ZONE = 60;
  const FACE_PADDING = 80;

  const nearbyFrames = frames_face_trackers.filter(
    (f) => Math.abs(f.frame - videoFrame) <= WINDOW_SIZE,
  );

  if (nearbyFrames.length > 0) {
    const rights: number[] = [];

    for (const trackerFrame of nearbyFrames) {
      if (trackerFrame.tracks.length === 0) {
        continue;
      }

      const right = Math.max(
        ...trackerFrame.tracks.map((track) => track.centerX + track.width / 2),
      );

      rights.push(right);
    }

    if (rights.length > 0) {
      const avgRight = rights.reduce((a, b) => a + b, 0) / rights.length;

      const scaledRight = avgRight * scale;

      const faceRight = scaledRight + FACE_PADDING;

      const overflow = faceRight - safeZoneEnd;

      if (overflow > DEAD_ZONE) {
        translateX = -(overflow - DEAD_ZONE);
      }
    }
  }

  const maxTranslate = Math.max(0, scaledWidth - target_width);

  translateX = Math.max(-maxTranslate, Math.min(0, translateX));

  /**
   * Fade transition
   * dari Hook ke Main Video
   */
  const revealProgress = interpolate(videoFrame, [0, 15], [1, 0], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  return (
    <AbsoluteFill>
      {/* ========================= */}
      {/* HOOK */}
      {/* ========================= */}

      {hook.duration_ms > 0 && (
        <Sequence from={0} durationInFrames={hookFrames}>
          <>
            <HookScreen
              videoPath={video_path}
              startFrame={hookPreviewFrame}
              headline={headline}
            />

            <Audio src={staticFile(hook.audio_path)} />
          </>
        </Sequence>
      )}

      {/* ========================= */}
      {/* MAIN VIDEO */}
      {/* ========================= */}

      <Sequence from={hookFrames}>
        <AbsoluteFill
          style={{
            overflow: "hidden",
          }}
        >
          <OffthreadVideo
            src={staticFile(video_path)}
            style={{
              position: "absolute",

              width: scaledWidth,
              height: target_height,

              transform: `translateX(${translateX}px)`,

              objectFit: "cover",

              filter: "contrast(1.08) saturate(1.1) brightness(0.96)",
            }}
          />

          {/* Subtitle */}

          <Subtitle words={subtitle} highlights={word_highlights} />

          {/* Curtain Top */}

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

          {/* Curtain Bottom */}

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

          {/* Hook -> Video Transition */}

          <div
            style={{
              position: "absolute",
              inset: 0,

              background: `rgba(0,0,0,${revealProgress})`,

              pointerEvents: "none",

              zIndex: 10000,
            }}
          />
        </AbsoluteFill>
      </Sequence>
    </AbsoluteFill>
  );
};
