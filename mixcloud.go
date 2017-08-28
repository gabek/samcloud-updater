package main

import (
	"log"
	"path"
	"strings"
	"sync"
	
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

func processMixcloud() {
	mixcloudUsers := getConfig().MixcloudUsers
	mixcloudTags := getConfig().MixcloudTags

	var wg sync.WaitGroup
	mixCount := len(mixcloudUsers) + len(mixcloudTags)
	wg.Add(mixCount)

	// Users
	for _, user := range mixcloudUsers {
		mixcloudURL := "http://www.mixcloud.com/" + user.Username

		go func(mixcloudURL string) {
			defer wg.Done()
			processMixcloudURL(mixcloudURL, nil)
		}(mixcloudURL)
	}

	// Tags
	for _, tag := range mixcloudTags {
		mixcloudURL := "https://www.mixcloud.com/discover/" + tag + "/?order=popular"

		go func(mixcloudURL string) {
			defer wg.Done()
			processMixcloudURL(mixcloudURL, &tag)
		}(mixcloudURL)
	}

	wg.Wait()
}

func processMixcloudURL(mixcloudURL string, finalFilename *string) {
	mixcloudDetails := detailsForMixcloudURL(mixcloudURL)
	if mixcloudDetails == nil {
		return
	}

	extension := path.Ext(mixcloudDetails.AudioURL)
	filename := "downloads/" + GenerateSlug(mixcloudDetails.Username) + "-" + GenerateSlug(mixcloudDetails.EpisodeTitle) + extension

	if !HasPreviouslyDownloaded(mixcloudDetails.OriginalTrackURL) {
		log.Println("Downloading " + mixcloudDetails.Username + ": " + mixcloudDetails.EpisodeTitle)
		downloadFile(filename, mixcloudDetails.AudioURL)

		// If the caller passed in a specific filename that we should
		// save this data to use that.  Otherwise save using the username.
		var newFilename = "uploads/"
		if finalFilename == nil {
			newFilename += GenerateSlug(mixcloudDetails.Username) + ".mp3"
		} else {
			newFilename += *finalFilename + ".mp3"
		}

		TranscodeHLSToMp3(mixcloudDetails.AudioURL, newFilename, mixcloudDetails.Username, mixcloudDetails.EpisodeTitle)
		AddFileToUploadList(newFilename)
		MarkFileAsDownloaded(mixcloudDetails.OriginalTrackURL)
	} else {
		log.Println("Nothing new for " + mixcloudDetails.Username)
	}
}

func detailsForMixcloudURL(url string) *MixcloudDetails {
    defer func() {
        if err := recover(); err != nil {
			log.Println("Unable to process " + url + ". There may be something wrong with Mixcloud.com or your host may have been banned from accessing it.")
        }
    }()

	htmlString := htmlForURL(url)

	root, _ := html.Parse(strings.NewReader(htmlString))

	episode := scrape.FindAllNested(root, scrape.ByClass("card"))[0]
	episodeInfo := scrape.FindAllNested(episode, scrape.ByClass("play-button"))[0]
	previewURL := (scrape.Attr(episodeInfo, "m-preview"))

	// We keep a reference to the web URL of this Mixcloud episode to log
	// that it has been downloaded.  This is due to the actual URL of the audio
	// changing from one request to the next.
	episodeURL := "https://www.mixcloud.com" + (scrape.Attr(episodeInfo, "m-url"))

	// If we don't have any audio to handle then this mix is of no use to us
	if previewURL == "" {
		return nil
	}

	episodeTitle := (scrape.Attr(episodeInfo, "m-title"))
	username := (scrape.Attr(episodeInfo, "m-owner-name"))
	return &MixcloudDetails{episodeTitle, fullAudioFromPreviewURL(previewURL), username, episodeURL}
}

func fullAudioFromPreviewURL(previewURL string) string {
	length := strings.Index(previewURL, "preview")
	server := previewURL[0 : length-1]
	cdnServer := strings.Replace(server, "audiocdn", "testcdn", -1)
	mixIdentifier := previewURL[length+9 : len(previewURL)]
	audioFilePath := strings.Replace(mixIdentifier, ".mp3", ".m4a", -1)
	fullURL := cdnServer + "/secure/hls/" + audioFilePath + "/index.m3u8"

	return fullURL
}

type MixcloudDetails struct {
	EpisodeTitle string
	AudioURL     string
	Username     string
	OriginalTrackURL	string
}
