import { AbsoluteFill, OffthreadVideo, Sequence, staticFile } from "remotion";

import { ClipMetadata } from "./types";
import { Subtitle } from "./components/Subtitle";

export const ClipComposition: React.FC<ClipMetadata> = ({
  video_path,
  subtitle,
  word_highlights,
}) => {
  return (
    <AbsoluteFill>
      <Sequence from={0}>
        <OffthreadVideo
          src={staticFile(video_path)}
          style={{
            position: "absolute",
            width: "100%",
            height: "100%",
            objectFit: "cover",
          }}
        />
      </Sequence>

      <Sequence from={0}>
        <div
          style={{
            position: "absolute",
            bottom: 0,
            left: 0,
            right: 0,
            height: "35%",
            background:
              "linear-gradient(to top, rgba(0,0,0,0.85) 0%, rgba(0,0,0,0.4) 50%, transparent 100%)",
            pointerEvents: "none",
          }}
        />

        <Subtitle words={subtitle} highlights={word_highlights} />
      </Sequence>
    </AbsoluteFill>
  );
};
