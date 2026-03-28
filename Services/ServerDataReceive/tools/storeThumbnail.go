package tools

import (
	"os/exec"
)

func FirstFrameThumbnail(videoURL string) error {
	input := "./MediaData/videos/rawVideos/" + videoURL + ".mp4"
	output := "./MediaData/thumbnails/" + videoURL + ".jpg"

	cmd := exec.Command("./Services/ServerDataReceive/tools/ffmpeg",
		"-ss", "0",
		"-i", input,
		"-frames:v", "1",
		"-vf", "scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2:black",
		"-q:v", "10", // lower quality → smaller file (range: 1 best, 31 worst)
		output,
	)

	return cmd.Run()
}