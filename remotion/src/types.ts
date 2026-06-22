export type Word = {
  word: string;
  start_ms: number;
  end_ms: number;
};

export type ClipMetadata = {
  title: string;
  width: number;
  height: number;
  duration: number;
  caption: string;
  video_path: string;
  subtitle: Word[];
  word_highlights: string[];
};
