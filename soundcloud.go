package main

import (
    "sync"
    "log"
)

type SoundcloudTrack struct {
    Kind string `json:"kind"`
    Id   int `json:"id"`
    Title   string `json:"title"`
    Description string `json:"description"`
    Stream string `json:"stream_url"`
}

func processSoundcloud() {
    // This is a hard-coded client id I found elsewhere on the internet.
    // If it ever gets turned off it'll need to be replaced.
    clientId := "175c043157ffae2c6d5fed16c3d95a4c"

    soundcloudUsers := getConfig().SoundcloudUsers
	var wg sync.WaitGroup
	wg.Add(len(soundcloudUsers))

	for _, soundcloudUser := range soundcloudUsers {
        var tracks []SoundcloudTrack
        tracklistUrl := "https://api.soundcloud.com/users/" + soundcloudUser.UserId + "/tracks?client_id=" + clientId
        getJson(tracklistUrl, &tracks)
        track := tracks[0]
        audio := track.Stream + "?client_id=" + clientId
    	filename := "downloads/" + GenerateSlug(soundcloudUser.UserId) + "-" + GenerateSlug(track.Title) + ".mp3"

    	if !FileExists(filename) {
		    log.Println("Downloading " + soundcloudUser.Title + ": " + track.Title)
		    downloadFile(filename, audio)
            newFilename := "uploads/" + GenerateSlug(soundcloudUser.UserId) + ".mp3"
    		TranscodeToMP3(filename, newFilename, soundcloudUser.UserId, track.Title)
            AddFileToUploadList(newFilename)
        }
    }
}