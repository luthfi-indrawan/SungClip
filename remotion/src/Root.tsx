import React from "react";
import { Composition } from "remotion";
import { ClipComposition } from "./Composition";

export const Root: React.FC = () => {
  return (
    <>
      <Composition
        id="clip"
        component={ClipComposition}
        defaultProps={{
          title: "missing",
          target_width: 0,
          target_height: 0,
          ori_width: 0,
          ori_height: 0,
          total_frames: 300,
          fps: 30,
          caption: "missing",
          video_path: "missing",
          subtitle: [],
          word_highlights: [],
          frames_face_trackers: [],
        }}
        calculateMetadata={({ props }) => {
          return {
            durationInFrames: props.total_frames,
            width: props.target_width,
            height: props.target_height,
            fps: props.fps,
          };
        }}
      />
    </>
  );
};
