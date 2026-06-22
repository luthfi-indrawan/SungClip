import React from "react";
import { Composition } from "remotion";
import { ClipComposition } from "./Composition";

export const Root: React.FC = () => {
  return (
    <>
      <Composition
        id="clip"
        component={ClipComposition}
        fps={30}
        defaultProps={{
          title: "missing",
          width: 0,
          height: 0,
          duration: 150,
          caption: "missing",
          video_path: "missing",
          subtitle: [],
          word_highlights: [],
        }}
        calculateMetadata={({ props }) => {
          return {
            durationInFrames: props.duration,
            width: props.width,
            height: props.height,
          };
        }}
      />
    </>
  );
};
