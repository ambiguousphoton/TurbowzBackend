package tools

import (
	"os/exec"
)

func FirstFrameThumbnail(videoURL string) error {
	input := "./MediaData/videos/rawVideos/" + videoURL + ".mp4"
	output := "./MediaData/thumbnails/" + videoURL + ".jpg"

		cmd := exec.Command("./Services/ServerDataReceive/tools/ffmpeg",
		"-ss", "00:00:01",          // 👈 seek BEFORE input (fast seek)
		"-i", input,

		"-frames:v", "1",

		"-vf", "scale=720:1280:force_original_aspect_ratio=decrease," +
			"pad=720:1280:(ow-iw)/2:(oh-ih)/2:black",

		"-q:v", "20",               // 👈 smaller size, still good quality
		"-an",                      // no audio processing (faster)
		"-sn",                      // no subtitles

		output,
	)

	return cmd.Run()
}