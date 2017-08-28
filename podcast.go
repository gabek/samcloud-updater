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

			urlFilename := path.Base(audioURL)
			originalFilename := urlFilename
			generatedFilename := GenerateSlug(podcast.Title) + ".mp3"
			filename := "downloads/" + originalFilename

			if !HasPreviouslyDownloaded(audioURL) {
				log.Println("Downloading " + podcast.Title + ": " + episodeName)

				downloadFile(filename, audioURL)

				localFilenameWithPath := "uploads/" + generatedFilename
				TranscodeToMP3(filename, localFilenameWithPath, podcast.Title, episodeName)

				AddFileToUploadList(localFilenameWithPath)
				MarkFileAsDownloaded(audioURL)
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
