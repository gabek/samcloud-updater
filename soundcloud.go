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
        err := getJson(tracklistUrl, &tracks)
        if err != nil {
            log.Println(err)
            return
        }

        track := tracks[0]
        audioURL := track.Stream + "?client_id=" + clientId
        filename := "downloads/" + GenerateSlug(soundcloudUser.Username) + "-" + GenerateSlug(track.Title) + ".mp3"
        
        if !HasPreviouslyDownloaded(audioURL) {
            log.Println("Downloading " + soundcloudUser.Username + ": " + track.Title)
            err := downloadFile(filename, audioURL)
            if err != nil {
                log.Println(err)
                return
            }
            localFilenameWithPath := "uploads/" + GenerateSlug(soundcloudUser.Username) + ".mp3"
            err = TranscodeToMP3(filename, localFilenameWithPath, soundcloudUser.Username, track.Title)
            if err != nil {
                log.Println(err)
                return
            }
            
            AddFileToUploadList(localFilenameWithPath)
            MarkFileAsDownloaded(audioURL)
        } else {
            log.Println("Nothing new for " + soundcloudUser.Username)
        }
	}
}