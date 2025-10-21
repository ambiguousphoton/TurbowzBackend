package tools

import (
	"os/exec"
)

func FirstFrameThumbnail(video_url string) error {
	input := "./MediaData/videos/rawVideos/" + video_url + ".mp4" 
	output := "./MediaData/thumbnails/" +  video_url + ".jpg" 

	// -ss 0 → start at the very beginning
	// -frames:v 1 → take one frame
	cmd := exec.Command("./Services/ServerDataReceive/tools/ffmpeg",
		"-ss", "0",
		"-i", input,
		"-frames:v", "1",
		// "-vf", "scale=w=1280:h=720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2",   //   ----> forcing  aspect ratio
		"-q:v", "2", // quality (lower is better)
		output,
	)
	
	return cmd.Run()
}