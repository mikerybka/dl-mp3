package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [youtube url]\n", os.Args[0])
		os.Exit(1)
	}
	ytURL := os.Args[1]
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-x", "--audio-format", "mp3", ytURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	os.Exit(cmd.ProcessState.ExitCode())
}
