#!/usr/bin/env python3
"""Generate a left-to-right scrolling video of a wide panoramic image.

Usage:
Run from this repo root:
    python make_scroll_video.py [--input PATH] [--output PATH]
                                [--duration SECONDS] [--fps FPS]
                                [--viewport-width PX] [--scale FACTOR]

Requires: Pillow, ffmpeg on PATH.
"""

import argparse
import subprocess
import time

from PIL import Image


DEFAULT_INPUT = "visualization/decks_log_scale_2022-01_2026-03.png"
DEFAULT_OUTPUT = "visualization/master_duel_history_scroll.mp4"
DEFAULT_DURATION = 56  # seconds
DEFAULT_FPS = 60
DEFAULT_VIEWPORT_WIDTH = 2400  # pixels on the source image
DEFAULT_SCALE = 1.0  # output is this fraction of the cropped size


def make_scroll_video(input_path, output_path, duration, fps, viewport_width, scale):
    img = Image.open(input_path).convert("RGB")
    src_w, src_h = img.size
    print(f"Source image: {src_w}x{src_h}")

    # Clamp viewport to image width
    viewport_width = min(viewport_width, src_w)

    out_w = round(viewport_width * scale)
    out_h = round(src_h * scale)
    # ffmpeg requires even dimensions for H.264
    out_w += out_w % 2
    out_h += out_h % 2
    print(f"Output frame: {out_w}x{out_h}")

    num_frames = duration * fps
    scroll_range = src_w - viewport_width  # total pixels to scroll
    print(f"Frames: {num_frames}, scroll range: {scroll_range}px")

    ffmpeg_cmd = [
        "ffmpeg", "-y",
        "-loglevel", "error",
        "-f", "rawvideo",
        "-vcodec", "rawvideo",
        "-s", f"{out_w}x{out_h}",
        "-pix_fmt", "rgb24",
        "-r", str(fps),
        "-i", "pipe:0",
        "-vcodec", "libx264",
        "-pix_fmt", "yuv420p",
        "-crf", "18",  # quality: lower = better, 18 is near-lossless
        "-preset", "fast",
        output_path,
    ]

    proc = subprocess.Popen(ffmpeg_cmd, stdin=subprocess.PIPE)

    total_frames = num_frames + 4 * fps  # 2s hold at start + 2s hold at end
    start = time.monotonic()

    def print_progress(done):
        elapsed = time.monotonic() - start
        if done > 0:
            remaining = elapsed / done * (total_frames - done)
            print(f"  {done}/{total_frames} frames  ~{remaining:.0f}s remaining    ", end="\r")

    try:
        # Hold the first frame for 2 seconds
        first_crop = img.crop((0, 0, viewport_width, src_h))
        if scale != 1.0:
            first_crop = first_crop.resize((out_w, out_h), Image.LANCZOS)
        first_bytes = first_crop.tobytes()
        for i in range(2 * fps):
            proc.stdin.write(first_bytes)
            if i % fps == 0:
                print_progress(i)

        for frame_index in range(num_frames):
            if scroll_range > 0:
                t = frame_index / (num_frames - 1)
                x = round(t * scroll_range)
            else:
                x = 0

            crop = img.crop((x, 0, x + viewport_width, src_h))

            if scale != 1.0:
                crop = crop.resize((out_w, out_h), Image.LANCZOS)

            proc.stdin.write(crop.tobytes())

            if frame_index % fps == 0:
                print_progress(2 * fps + frame_index)

        # Hold the last frame for 2 extra seconds
        last_crop = img.crop((scroll_range, 0, scroll_range + viewport_width, src_h))
        if scale != 1.0:
            last_crop = last_crop.resize((out_w, out_h), Image.LANCZOS)
        last_bytes = last_crop.tobytes()
        for i in range(2 * fps):
            proc.stdin.write(last_bytes)
            if i % fps == 0:
                print_progress(2 * fps + num_frames + i)

        print(f"\nDone. Saved: {output_path}")
    finally:
        proc.stdin.close()
        proc.wait()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--input", default=DEFAULT_INPUT)
    parser.add_argument("--output", default=DEFAULT_OUTPUT)
    parser.add_argument("--duration", type=int, default=DEFAULT_DURATION,
                        help="Video length in seconds")
    parser.add_argument("--fps", type=int, default=DEFAULT_FPS)
    parser.add_argument("--viewport-width", type=int, default=DEFAULT_VIEWPORT_WIDTH,
                        help="Width of the visible window on the source image (pixels)")
    parser.add_argument("--scale", type=float, default=DEFAULT_SCALE,
                        help="Output scale factor relative to viewport size")
    args = parser.parse_args()

    make_scroll_video(
        args.input,
        args.output,
        args.duration,
        args.fps,
        args.viewport_width,
        args.scale,
    )


if __name__ == "__main__":
    main()
