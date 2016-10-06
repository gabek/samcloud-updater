package main

import (
	"log"
	"path"
	"sync"

	"github.com/SlyMarbo/rss"
)

func processPodcasts() {
	podcasts := getConfig().Podcasts
	var wg sync.WaitGroup
	wg.Add(len(podcasts))

	for _, podcast := range podcasts {

		go func(podcast Podcast) {
			defer wg.Done()
			episodeName, audioURL := itemForRssFeedURL(podcast.URL)

			filetype := path.Ext(audioURL)
			urlFilename := path.Base(audioURL)
			filename := "downloads/" + GenerateSlug(urlFilename) + filetype

			if !FileExists(filename) {
				log.Println("Downloading " + podcast.Title + ": " + episodeName)

				downloadFile(filename, audioURL)

				newFilename := "uploads/" + GenerateSlug(podcast.Title) + ".mp3"
				TranscodeToMP3(filename, newFilename, podcast.Title, episodeName)
				filename = newFilename

				AddFileToUploadList(filename)
			} else {
				log.Println("Nothing new for " + podcast.Title)
			}
		}(podcast)
	}
	wg.Wait()
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
