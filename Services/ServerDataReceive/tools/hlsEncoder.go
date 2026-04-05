package tools

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)


func EncodeVideoToHLS(videoUID string) {

	input := "./MediaData/videos/rawVideos/" + videoUID + ".mp4" // change to your file path
	outputDir := "./MediaData/videos/hlsEncodedVideos/" +  videoUID
	log.Println(input)
	// Create output folder if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal("Error creating output folder:", err)
	}

	// Output playlist path
	outputM3U8 := filepath.Join(outputDir, "playlist.m3u8")

	// FFmpeg command: compress + HLS in one step
	cmd := exec.Command("./Services/ServerDataReceive/tools/ffmpeg",
		"-i", input,

		"-r", "30",                 // better UX than 24 for mobile
		"-g", "60",                 // 2x FPS

		"-c:v", "libx264",
		"-preset", "veryfast",      // faster encoding (important for scale)
		"-crf", "25",               // slightly lower quality → huge savings
		"-maxrate", "800k",         // cap bitrate (VERY important)
		"-bufsize", "1200k",

		"-vf", "scale=720:-2",      // force 720p (TikTok-like)

		"-c:a", "aac",
		"-b:a", "96k",              // reduce audio bitrate

		"-movflags", "+faststart",  // faster playback start

		"-hls_time", "4",           // shorter segments → faster start
		"-hls_list_size", "0",
		"-hls_flags", "independent_segments",

		"-f", "hls",
		outputM3U8,
	)

	// Optional: stream ffmpeg output to terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	log.Println("Processing video to HLS...")
	if err := cmd.Run(); err != nil {
		log.Fatal("FFmpeg failed:", err)
	}

	log.Println("HLS files created in:", outputDir)
}
