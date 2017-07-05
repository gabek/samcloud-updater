package main

import (
	"log"
	"os"
	"strconv"
)

var config = getConfig()
var uploadListFile = "toUpload.m3u"
var stationLogin = "\"" + config.Station.Username + ";" + config.Station.Password + "\""
var stationID = strconv.Itoa(config.Station.Id)
var stationPlaylist = config.Station.Playlist

func init() {
	samToolsCheck()
	setup()
}

func main() {
	log.Println("*** SAM Broadcaster Cloud Updater ***")

	os.Remove(uploadListFile)

	log.Println("Podcasts:")
	processPodcasts()
	log.Println("Mixcloud:")
	processMixcloud()
	log.Println("Soundcloud:")
	processSoundcloud()

	if FileExists(uploadListFile) {
		upload()
		addToStation()
		os.Remove(uploadListFile)
	}
}

func setup() {
	// If they don't already exist, create the directories we'll username
	// for storing the downloaded audio files.
	_ = os.Mkdir("./uploads", os.ModePerm)
	_ = os.Mkdir("./downloads", os.ModePerm)
}
