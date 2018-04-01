package main

import (
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
	"strconv"
	"strings"
	"sync"
)

var targetArchive = "Death Guild"

func processDNALounge() {
	url := "https://www.dnalounge.com/webcast/archive/"
	htmlString := htmlForURL(url)

	root, _ := html.Parse(strings.NewReader(htmlString))
	archives := scrape.FindAllNested(root, scrape.ByClass("ebody"))

	var wg sync.WaitGroup
	archiveCount := 1
	wg.Add(archiveCount)

	// Loop over the archives to find Death Guild
	for _, archive := range archives {
		archiveInfo := scrape.FindAllNested(archive, scrape.ByClass("etitle"))[0]
		archiveTitle := scrape.Text(archiveInfo)

		if archiveTitle == targetArchive {
			archiveDate := scrape.FindAllNested(archive, scrape.ByClass("edate"))[0]
			dateString := scrape.Text(archiveDate)
			archiveIdentifier := targetArchive + "-" + dateString
			if HasPreviouslyDownloaded(archiveIdentifier) {
				log.Println("Nothing new for " + targetArchive)
				return
			}

			go func(mixcloudURL string) {
				defer wg.Done()
				processDNALoungeArchive(archive, archiveIdentifier, targetArchive, dateString)
			}(archiveIdentifier)

			wg.Wait()

			break
		}
	}
}

func processDNALoungeArchive(archive *html.Node, archiveIdentifier string, artist string, title string) {
	audioLinkContainer := scrape.FindAllNested(archive, scrape.ByClass("lbox"))[0]
	audioLinks := scrape.FindAllNested(audioLinkContainer, scrape.ByClass("l2"))

	var urls []string

	for _, link := range audioLinks {
		chunkTitle := scrape.Text(link)

		if chunkTitle == "10:30" {
			urlNode := scrape.FindAllNested(link, scrape.ByTag(atom.A))[0]
			playlistURL := scrape.Attr(urlNode, "href")
			audioURL := strings.Replace(playlistURL, "m3u", "mp3", -1)
			if !stringInSlice(audioURL, urls) {
				urls = append(urls, audioURL)
			}
		}
	}

	var filenames []string

	for i, url := range urls {
		log.Println("Downloading " + url)
		filename := "downloads/" + GenerateSlug(archiveIdentifier+"-"+strconv.Itoa(i)) + ".mp3"

		maxFileSize := 60000000 // in bytes
		err := downloadFile(filename, url, maxFileSize)
		if err != nil {
			log.Println(err)
		}
		filenames = append(filenames, filename)
	}

	destination := "uploads/" + GenerateSlug(artist) + ".mp3"
	TranscodeFilesToMP3(filenames, destination, artist, title)
	AddFileToUploadList(destination)
	MarkFileAsDownloaded(archiveIdentifier)
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
