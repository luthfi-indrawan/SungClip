import { useCurrentFrame, useVideoConfig, interpolate, Easing } from "remotion";

export type Word = {
  word: string;
  start_ms: number;
  end_ms: number;
};

export const Subtitle = ({
  words,
  highlights,
}: {
  words: Word[];
  highlights: string[];
}) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  const currentMs = (frame / fps) * 1000;

  const normalize = (s: string) => s.replace(/[^\w]/g, "").toLowerCase();

  const activeWord = words.find(
    (w) => currentMs >= w.start_ms && currentMs <= w.end_ms,
  );

  if (!activeWord) {
    return null;
  }

  const isHighlight = highlights.some(
    (h) => normalize(h) === normalize(activeWord.word),
  );

  const wordDurationMs = activeWord.end_ms - activeWord.start_ms;

  const animDurationFrames = Math.max((wordDurationMs / 1000) * fps * 0.4, 5);

  const wordStartFrame = (activeWord.start_ms / 1000) * fps;

  const wordEndFrame = wordStartFrame + animDurationFrames;

  const progress = interpolate(frame, [wordStartFrame, wordEndFrame], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
    easing: Easing.out(Easing.back(1.2)),
  });

  const scale = interpolate(progress, [0, 1], [0.5, 1]);

  return (
    <div
      style={{
        position: "absolute",
        bottom: 180,
        left: 0,
        right: 0,

        display: "flex",
        justifyContent: "center",
        alignItems: "center",

        padding: "0 60px",

        fontFamily: "'Inter', 'Segoe UI', 'Arial', sans-serif",

        fontSize: 72,
        fontWeight: 900,

        textAlign: "center",

        textShadow: "0 2px 12px rgba(0,0,0,0.9), 0 0 30px rgba(0,0,0,0.5)",
      }}
    >
      <span
        style={{
          color: isHighlight ? "#FFD700" : "#FFFFFF",

          opacity: progress,
          transform: `scale(${scale})`,
          display: "inline-block",
          whiteSpace: "nowrap",
        }}
      >
        {activeWord.word}
      </span>
    </div>
  );
};
