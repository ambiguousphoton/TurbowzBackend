package tools

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func FirstFrameThumbnail(videoURL string) error {
	input := "./MediaData/videos/rawVideos/" + videoURL + ".mp4"
	output := "./MediaData/thumbnails/" + videoURL + ".jpg"

	// ✅ Ensure output directory exists
	err := os.MkdirAll(filepath.Dir(output), os.ModePerm)
	if err != nil {
		log.Printf("❌ Failed to create thumbnail directory: %v", err)
		return err
	}

	// ✅ Check if input exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Printf("❌ Input video not found: %s", input)
		return err
	}

	cmd := exec.Command("./Services/ServerDataReceive/tools/ffmpeg",
		"-y",                      // overwrite output
		"-ss", "00:00:01",
		"-i", input,

		"-frames:v", "1",

		"-vf", "thumbnail,scale=720:-2", // safer filter

		"-q:v", "20",
		"-an",
		"-sn",

		output,
	)

	// ✅ Capture BOTH stdout + stderr
	outputBytes, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("❌ FFmpeg thumbnail failed for video %s", videoURL)
		log.Printf("📁 Input: %s", input)
		log.Printf("📁 Output: %s", output)
		log.Printf("⚠️ Error: %v", err)
		log.Printf("📜 FFmpeg Output:\n%s", string(outputBytes))
		return err
	}

	log.Printf("✅ Thumbnail created successfully for video %s", videoURL)
	return nil
}