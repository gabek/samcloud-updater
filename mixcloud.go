package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

func processMixcloud() {
	mixcloudUsers := getConfig().Mixcloud

	for _, user := range mixcloudUsers {
		fmt.Print(user.Title + "...")

		episodeTitle, url := detailsForMixcloudUser(user.Username)
		baseFile := path.Base(url)
		filename := "transcode/" + baseFile

		if !FileExists(filename) {
			fmt.Println(episodeTitle)
			downloadFile(filename, url)
			newFilename := "audio/" + GenerateSlug(episodeTitle) + ".mp3"
			TranscodeToMP3(filename, newFilename)
			setID3TagsForFile(newFilename, user.Title, episodeTitle)
			AddFileToUploadList(newFilename)
		} else {
			fmt.Println("mix exists")
		}
	}
}

func detailsForMixcloudUser(user string) (string, string) {
	url := "http://www.mixcloud.com/" + user
	htmlString := htmlForURL(url)

	root, _ := html.Parse(strings.NewReader(htmlString))

	episode := scrape.FindAllNested(root, scrape.ByClass("card-elements-container"))[0]
	episodeInfo := scrape.FindAllNested(episode, scrape.ByClass("play-button"))[0]
	previewURL := (scrape.Attr(episodeInfo, "m-preview"))
	episodeTitlte := (scrape.Attr(episodeInfo, "m-title"))
	return episodeTitlte, fullAudioFromPreviewURL(previewURL)
}

func fullAudioFromPreviewURL(previewURL string) string {
	length := strings.Index(previewURL, "preview")
	server := previewURL[0 : length-1]
	cdnServer := strings.Replace(server, "audiocdn", "stream", -1)
	mixIdentifier := previewURL[length+9 : len(previewURL)]
	audioFilePath := strings.Replace(mixIdentifier, ".mp3", ".m4a", -1)
	fullUrl := cdnServer + "/c/m4a/64/" + audioFilePath
	return fullUrl
}