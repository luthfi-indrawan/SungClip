import React from "react";
import {
  AbsoluteFill,
  OffthreadVideo,
  interpolate,
  staticFile,
  useCurrentFrame,
} from "remotion";

type Props = {
  videoPath: string;
  startFrame: number;
  headline: string;
};

export const HookScreen: React.FC<Props> = ({
  videoPath,
  startFrame,
  headline,
}) => {
  const frame = useCurrentFrame();

  const zoom = interpolate(frame, [0, 120], [1.05, 1.15], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  const opacity = interpolate(frame, [0, 12], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  const translateY = interpolate(frame, [0, 12], [60, 0], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
  });

  return (
    <AbsoluteFill>
      <OffthreadVideo
        src={staticFile(videoPath)}
        muted
        startFrom={startFrame}
        style={{
          width: "100%",
          height: "100%",
          objectFit: "cover",

          transform: `scale(${zoom})`,

          filter:
            "grayscale(35%) brightness(0.45) contrast(1.25) saturate(0.8)",
        }}
      />

      {/* cinematic dark overlay */}

      <div
        style={{
          position: "absolute",
          inset: 0,

          background:
            "linear-gradient(to bottom, rgba(0,0,0,0.35), rgba(0,0,0,0.75))",
        }}
      />

      {/* vignette */}

      <div
        style={{
          position: "absolute",
          inset: 0,

          background: `
            radial-gradient(
              circle,
              rgba(0,0,0,0) 40%,
              rgba(0,0,0,0.45) 100%
            )
          `,
        }}
      />

      <AbsoluteFill
        style={{
          justifyContent: "center",
          paddingLeft: 80,
          paddingRight: 80,
        }}
      >
        <div
          style={{
            opacity,
            transform: `translateY(${translateY}px)`,
          }}
        >
          <div
            style={{
              width: 140,
              height: 10,

              borderRadius: 999,

              background: "#FFD700",

              marginBottom: 36,
            }}
          />

          <div
            style={{
              color: "#FFFFFF",

              fontSize: 92,

              fontWeight: 900,

              lineHeight: 1.05,

              textTransform: "uppercase",

              maxWidth: 950,

              textShadow: "0px 8px 30px rgba(0,0,0,0.7)",
            }}
          >
            {headline}
          </div>

          <div
            style={{
              marginTop: 30,

              color: "rgba(255,255,255,0.85)",

              fontSize: 28,

              fontWeight: 700,

              letterSpacing: 4,

              textTransform: "uppercase",
            }}
          >
            Podcast Highlight
          </div>
        </div>
      </AbsoluteFill>
    </AbsoluteFill>
  );
};
