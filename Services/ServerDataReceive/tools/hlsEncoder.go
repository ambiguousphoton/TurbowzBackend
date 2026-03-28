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
		"-i", input,                // input file
		"-c:v", "libx264",          // H.264 video codec
		"-preset", "fast",          // encoding speed/efficiency tradeoff
		"-crf", "23",                // quality (lower = better quality, bigger file)
		"-c:a", "aac",              // audio codec
		"-b:a", "128k",             // audio bitrate
		"-hls_time", "6",           // segment length (seconds)
		"-hls_list_size", "0",      // include all segments in playlist
		"-f", "hls",                // output format
		outputM3U8,                 // playlist output
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
