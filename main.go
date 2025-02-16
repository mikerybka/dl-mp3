package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	SpotifyAPIURL = "https://api.spotify.com/v1/tracks/"
	YouTubeAPIURL = "https://www.googleapis.com/youtube/v3/search"
)

type SpotifyTrack struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}

type YouTubeSearchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

func getSpotifyTrack(spotifyID, token string) (*SpotifyTrack, error) {
	url := SpotifyAPIURL + spotifyID
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get track info: %s", resp.Status)
	}

	var track SpotifyTrack
	if err := json.NewDecoder(resp.Body).Decode(&track); err != nil {
		return nil, err
	}
	return &track, nil
}

func searchYouTube(query, apiKey string) (*YouTubeSearchResponse, error) {
	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("q", query)
	params.Set("maxResults", "1")
	params.Set("type", "video")
	params.Set("key", apiKey)

	url := YouTubeAPIURL + "?" + params.Encode()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube search failed: %s", resp.Status)
	}

	var result YouTubeSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func main() {
	spotifyToken := os.Getenv("SPOTIFY_API_TOKEN")
	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")

	if len(os.Args) < 2 {
		log.Fatal("Usage: %s <spotify_url>", os.Args[0])
	}

	// Parse spotify URL
	spotifyURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	spotifyID := strings.TrimPrefix(spotifyURL.Path, "/track/")

	// Look up song info on Spotify
	track, err := getSpotifyTrack(spotifyID, spotifyToken)
	if err != nil {
		log.Fatalf("Error fetching Spotify track: %v", err)
	}

	// Search for YouTube video
	artistNames := []string{}
	for _, artist := range track.Artists {
		artistNames = append(artistNames, artist.Name)
	}
	query := fmt.Sprintf("%s %s", track.Name, strings.Join(artistNames, " "))
	searchResults, err := searchYouTube(query, youtubeAPIKey)
	if err != nil {
		log.Fatalf("Error searching YouTube: %v", err)
	}
	if len(searchResults.Items) == 0 {
		log.Println("No results found.")
		return
	}
	youtubeID := searchResults.Items[0].ID.VideoID

	// Download mp3
	youtubeURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", youtubeID)
	mp3file := spotifyID + ".mp3"
	tmpFile := "dl.mp3"
	cmd := exec.Command("yt-dlp",
		"-o", tmpFile,
		"-f", "bestaudio",
		"-x",
		"--audio-format", "mp3",
		youtubeURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Set mp3 metadata
	artists := []string{}
	for _, a := range track.Artists {
		artists = append(artists, a.Name)
	}
	cmd = exec.Command("ffmpeg",
		"-i", tmpFile,
		"-metadata", fmt.Sprintf("title=%s", track.Name),
		"-metadata", fmt.Sprintf("artist=%s", strings.Join(artists, "; ")),
		mp3file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Delete temp file
	err = os.Remove(tmpFile)
	if err != nil {
		log.Fatal(err)
	}
}
