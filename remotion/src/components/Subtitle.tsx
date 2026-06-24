import { useCurrentFrame, useVideoConfig } from "remotion";

export type Word = {
  word: string;
  start_ms: number;
  end_ms: number;
};

type Chunk = {
  startIndex: number;
  endIndex: number;
  lines: number[][];
};

const normalize = (s: string) => s.replace(/[^\w]/g, "").toLowerCase();

const MAX_CHUNK_CHARS = 42;
const MAX_WORDS_PER_CHUNK = 10;
const MAX_LINES = 3;

const LINE_OFFSETS = [0, 70, 30];

function estimateWordWidth(word: string) {
  return word.length * 42;
}

function buildChunkLines(
  words: Word[],
  start: number,
  end: number,
): number[][] {
  const MAX_LINE_WIDTH = 650;

  const lines: number[][] = [];

  let currentLine: number[] = [];
  let currentWidth = 0;

  for (let i = start; i <= end; i++) {
    const width = estimateWordWidth(words[i].word) + 30;

    if (currentWidth + width > MAX_LINE_WIDTH && currentLine.length > 0) {
      lines.push(currentLine);

      currentLine = [i];
      currentWidth = width;
    } else {
      currentLine.push(i);
      currentWidth += width;
    }
  }

  if (currentLine.length > 0) {
    lines.push(currentLine);
  }

  if (lines.length <= MAX_LINES) {
    return lines;
  }

  const flattened = lines.flat();

  const wordsPerLine = Math.ceil(flattened.length / MAX_LINES);

  return [
    flattened.slice(0, wordsPerLine),
    flattened.slice(wordsPerLine, wordsPerLine * 2),
    flattened.slice(wordsPerLine * 2),
  ];
}

function buildChunks(words: Word[]): Chunk[] {
  const chunks: Chunk[] = [];

  let start = 0;

  while (start < words.length) {
    let end = start;
    let chars = 0;

    while (end < words.length) {
      const gap = end > start ? words[end].start_ms - words[end - 1].end_ms : 0;

      if (gap > 700) {
        break;
      }

      chars += words[end].word.length + 1;

      if (chars > MAX_CHUNK_CHARS) {
        break;
      }

      if (end - start + 1 > MAX_WORDS_PER_CHUNK) {
        break;
      }

      end++;
    }

    const actualEnd = Math.max(start, end - 1);

    chunks.push({
      startIndex: start,
      endIndex: actualEnd,
      lines: buildChunkLines(words, start, actualEnd),
    });

    start = actualEnd + 1;
  }

  return chunks;
}

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

  const activeIndex = words.findIndex(
    (w) => currentMs >= w.start_ms && currentMs <= w.end_ms,
  );

  const chunks = buildChunks(words);

  const activeChunk =
    chunks.find(
      (chunk) =>
        activeIndex >= chunk.startIndex && activeIndex <= chunk.endIndex,
    ) ?? null;

  if (!activeChunk) {
    return (
      <div
        style={{
          position: "absolute",
          top: 0,
          right: 0,

          width: 980,
          height: "100%",

          background:
            "linear-gradient(to left, rgba(0,0,0,0.94) 0%, rgba(0,0,0,0.78) 45%, transparent 100%)",
        }}
      />
    );
  }

  return (
    <div
      style={{
        position: "absolute",

        top: 0,
        right: 0,

        width: 980,
        height: "100%",

        display: "flex",
        alignItems: "center",

        padding: "0 80px",

        boxSizing: "border-box",

        background:
          "linear-gradient(to left, rgba(0,0,0,0.94) 0%, rgba(0,0,0,0.78) 45%, transparent 100%)",

        pointerEvents: "none",

        fontFamily: "'Inter','Segoe UI','Arial',sans-serif",
      }}
    >
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          width: "100%",
        }}
      >
        {activeChunk.lines.map((line, lineIndex) => (
          <div
            key={lineIndex}
            style={{
              marginLeft: LINE_OFFSETS[lineIndex] ?? 0,

              marginBottom: 14,

              fontSize: 72,

              fontWeight: 900,

              lineHeight: 1.05,

              display: "flex",
              flexWrap: "wrap",
            }}
          >
            {line.map((wordIndex) => {
              const word = words[wordIndex];

              const visible = wordIndex <= activeIndex;

              if (!visible) {
                return null;
              }

              const isHighlight = highlights.some(
                (h) => normalize(h) === normalize(word.word),
              );

              return (
                <span
                  key={wordIndex}
                  style={{
                    color: isHighlight ? "#FFD700" : "#FFFFFF",

                    marginRight: 12,

                    textShadow: "0 4px 18px rgba(0,0,0,0.95)",
                  }}
                >
                  {word.word}
                </span>
              );
            })}
          </div>
        ))}
      </div>
    </div>
  );
};
