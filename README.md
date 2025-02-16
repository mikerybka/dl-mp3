# dl-mp3

## Requirements
- ffmpeg
- yt-dlp

## Setup

### YouTube API
1. Log into [Google Cloud Console](https://console.cloud.google.com) and make sure YouTube Data API v3 is enabled.
1. Go to "APIs & Services" > "Credentials" and create a new API key.

### Spotify API
1. Log into the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard) and create a new app.
1. Name and Redirect URLs aren't important.
1. Go to your app settings to get Client ID and Client Secret.
1. Run `curl -X "POST" -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=client_credentials" -u "<client_id>:<client_secret>" https://accounts.spotify.com/api/token` in the terminal to get an API token.


## Install
```
go install github.com/mikerybka/dl-mp3@latest
```

## Usage
```
SPOTIFY_API_TOKEN=<spotify_api_token> YOUTUBE_API_KEY=<youtube_api_key> dl-mp3 <spotify_url>
```
