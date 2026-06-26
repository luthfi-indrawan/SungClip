export type Word = {
  word: string;
  start_ms: number;
  end_ms: number;
};

export type Frame = {
  frame: number;
  timeMs: number;
  tracks: Track[];
};

export type Track = {
  trackId: number;
  centerX: number;
  centerY: number;
  width: number;
  height: number;
  confidence: number;
};

export type HookMetadata = {
  text: string;
  audio_path: string;
  duration_ms: number;
};

export type ClipMetadata = {
  title: string;
  headline: string;
  hook: HookMetadata;
  fps: number;
  target_width: number;
  target_height: number;
  ori_width: number;
  ori_height: number;
  total_frames: number;
  video_path: string;
  caption: string;
  subtitle: Word[];
  word_highlights: string[];
  frames_face_trackers: Frame[];
};
