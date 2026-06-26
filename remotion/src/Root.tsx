import React from "react";
import { Composition } from "remotion";
import { SinglePodcastComposition } from "./SinglePodcastComposition";

export const Root: React.FC = () => {
  return (
    <>
      <Composition
        id="single-podcast"
        component={SinglePodcastComposition}
        defaultProps={{
          title: "missing",
          headline: "missing",
          hook: {
            text: "missing",
            audio_path: "missing",
            duration_ms: 0,
          },
          fps: 30,
          target_width: 0,
          target_height: 0,
          ori_width: 0,
          ori_height: 0,
          total_frames: 300,
          caption: "missing",
          video_path: "missing",
          subtitle: [],
          word_highlights: [],
          frames_face_trackers: [],
        }}
        calculateMetadata={({ props }) => {
          const hookFrames = Math.ceil(
            (props.hook.duration_ms / 1000) * props.fps,
          );

          return {
            durationInFrames: hookFrames + props.total_frames,

            width: props.target_width,
            height: props.target_height,

            fps: props.fps,
          };
        }}
      />
    </>
  );
};
