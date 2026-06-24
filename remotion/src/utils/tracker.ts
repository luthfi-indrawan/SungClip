import { Frame } from "../types";

export const getFaceTracker = (
  currentFrame: number,
  trackers: Frame[],
): Frame | null => {
  if (trackers.length === 0) {
    return null;
  }

  let nearest = trackers[0];
  let minDiff = Math.abs(trackers[0].frame - currentFrame);

  for (const tracker of trackers) {
    const diff = Math.abs(tracker.frame - currentFrame);

    if (diff < minDiff) {
      minDiff = diff;
      nearest = tracker;
    }
  }

  return nearest;
};
