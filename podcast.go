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
			err, episodeName, audioURL := itemForRssFeedURL(podcast.URL)
			if err != nil {
				return
			}

			urlFilename := path.Base(audioURL)
			originalFilename := urlFilename
			generatedFilename := GenerateSlug(podcast.Title) + ".mp3"
			filename := "downloads/" + originalFilename

			if !HasPreviouslyDownloaded(audioURL) {
				log.Println("Downloading " + podcast.Title + ": " + episodeName)

				err := downloadFile(filename, audioURL, 0)
				if err != nil {
					return
				}

				localFilenameWithPath := "uploads/" + generatedFilename
				err = TranscodeToMP3(filename, localFilenameWithPath, podcast.Title, episodeName)
				if err != nil {
					return
				}

				AddFileToUploadList(localFilenameWithPath)
				MarkFileAsDownloaded(audioURL)
			} else {
				log.Println("Nothing new for " + podcast.Title)
			}
		}(podcast)
	}
	wg.Wait()
}

func itemForRssFeedURL(url string) (error, string, string) {
	feed, err := rss.Fetch(url)
	if err != nil {
		log.Println(err)
		return err, "", ""
	}

	episode := feed.Items[0]
	episodeAudio := episode.Enclosures[0].Url
	episodeTitle := episode.Title
	return nil, episodeTitle, episodeAudio
}
