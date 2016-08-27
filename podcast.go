package main

import (
	"fmt"
	"path"

	"github.com/SlyMarbo/rss"
)

func processPodcasts() {
	podcasts := getConfig().Podcasts

	for _, podcast := range podcasts {
		fmt.Print(podcast.Title + "...")

		episodeName, audioURL := itemForRssFeedURL(podcast.URL)

		filename := ""
		filetype := path.Ext(audioURL)

		if filetype == ".m4a" {
			filename = "transcode/" + GenerateSlug(episodeName) + ".m4a"
		} else {
			filename = "audio/" + GenerateSlug(episodeName) + ".mp3"
		}

		if !FileExists(filename) {
			fmt.Println(episodeName)

			downloadFile(filename, audioURL)

			if filetype == ".m4a" {
				newFilename := "audio/" + GenerateSlug(episodeName) + ".mp3"
				TranscodeToMP3(filename, newFilename)
				filename = newFilename
			}

			setID3TagsForFile(filename, podcast.Title, episodeName)
			AddFileToUploadList(filename)
		} else {
			fmt.Println("episode exists")
		}
	}
}

func itemForRssFeedURL(url string) (string, string) {
	feed, err := rss.Fetch(url)
	if err != nil {
		// handle error.
	}

	episode := feed.Items[0]
	episodeAudio := episode.Enclosures[0].Url
	episodeTitle := episode.Title
	return episodeTitle, episodeAudio
}