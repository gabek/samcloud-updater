package main

import (
	"log"
	"path"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

func processMixcloudUsers() {
	mixcloudUsers := getConfig().MixcloudUsers

	for _, user := range mixcloudUsers {
		log.Print(user.Title + "...")
		mixcloudURL := "http://www.mixcloud.com/" + user.Username
		processMixcloudURL(mixcloudURL, nil)
	}
}

func processMixcloudTags() {
	mixcloudTags := getConfig().MixcloudTags

	for _, tag := range mixcloudTags {
		log.Print("Mixcloud tag: " + tag + "...")
		mixcloudURL := "https://www.mixcloud.com/discover/" + tag + "/?order=popular"
		processMixcloudURL(mixcloudURL, &tag)
	}
}

func processMixcloudURL(mixcloudURL string, finalFilename *string) {
	mixcloudDetails := detailsForMixcloudURL(mixcloudURL)
	if mixcloudDetails == nil {
		return
	}

	extension := path.Ext(mixcloudDetails.AudioURL)
	filename := "downloads/" + GenerateSlug(mixcloudDetails.Username) + "-" + GenerateSlug(mixcloudDetails.EpisodeTitle) + extension

	if !FileExists(filename) {
		log.Println("Downloading " + mixcloudDetails.EpisodeTitle)
		downloadFile(filename, mixcloudDetails.AudioURL)

		// If the caller passed in a specific filename that we should
		// save this data to use that.  Otherwise save using the username.
		var newFilename = "uploads/"
		if finalFilename == nil {
			newFilename += GenerateSlug(mixcloudDetails.Username) + ".mp3"
		} else {
			newFilename += *finalFilename + ".mp3"
		}

		TranscodeToMP3(filename, newFilename, mixcloudDetails.Username, mixcloudDetails.EpisodeTitle)
		AddFileToUploadList(newFilename)
	}
}

func detailsForMixcloudURL(url string) *MixcloudDetails {
	htmlString := htmlForURL(url)

	root, _ := html.Parse(strings.NewReader(htmlString))
	episode := scrape.FindAllNested(root, scrape.ByClass("card-elements-container"))[0]
	episodeInfo := scrape.FindAllNested(episode, scrape.ByClass("play-button"))[0]
	previewURL := (scrape.Attr(episodeInfo, "m-preview"))

	// If we don't have any audio to handle then this mix is of no use to us
	if previewURL == "" {
		return nil
	}

	episodeTitle := (scrape.Attr(episodeInfo, "m-title"))
	username := (scrape.Attr(episodeInfo, "m-owner-name"))
	return &MixcloudDetails{episodeTitle, fullAudioFromPreviewURL(previewURL), username}
}

func fullAudioFromPreviewURL(previewURL string) string {
	length := strings.Index(previewURL, "preview")
	server := previewURL[0 : length-1]
	cdnServer := strings.Replace(server, "audiocdn", "stream", -1)
	mixIdentifier := previewURL[length+9 : len(previewURL)]
	audioFilePath := strings.Replace(mixIdentifier, ".mp3", ".m4a", -1)
	fullURL := cdnServer + "/c/m4a/64/" + audioFilePath
	return fullURL
}

type MixcloudDetails struct {
	EpisodeTitle string
	AudioURL     string
	Username     string
}
