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
	episodeTitle, url, username := detailsForMixcloudURL(mixcloudURL)
	extension := path.Ext(url)
	filename := "downloads/" + GenerateSlug(username) + "-" + GenerateSlug(episodeTitle) + extension

	if !FileExists(filename) {
		log.Println("Downloading " + episodeTitle)
		downloadFile(filename, url)

		// If the caller passed in a specific filename that we should
		// save this data to use that.  Otherwise save using the username.
		var newFilename = ""
		if finalFilename == nil {
			newFilename = "uploads/" + GenerateSlug(username) + ".mp3"
		} else {
			newFilename = "uploads/" + *finalFilename + ".mp3"
		}

		TranscodeToMP3(filename, newFilename)
		setID3TagsForFile(newFilename, username, episodeTitle)
		AddFileToUploadList(newFilename)
	}
}

func detailsForMixcloudURL(url string) (string, string, string) {
	htmlString := htmlForURL(url)

	root, _ := html.Parse(strings.NewReader(htmlString))

	episode := scrape.FindAllNested(root, scrape.ByClass("card-elements-container"))[0]
	episodeInfo := scrape.FindAllNested(episode, scrape.ByClass("play-button"))[0]
	previewURL := (scrape.Attr(episodeInfo, "m-preview"))
	episodeTitle := (scrape.Attr(episodeInfo, "m-title"))
	username := (scrape.Attr(episodeInfo, "m-owner-name"))
	return episodeTitle, fullAudioFromPreviewURL(previewURL), username
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
