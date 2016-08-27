package main

import (
	"fmt"
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
	fmt.Println("*** The Bat Station mix updater ***")

	os.Remove(uploadListFile)

	processPodcasts()
	processMixcloud()

	if FileExists(uploadListFile) {
		upload()
		addToStation()
		os.Remove(uploadListFile)
	}
}

func setup() {
	// If they don't already exist, create the directories we'll username
	// for storing the downloaded audio files.
	_ = os.Mkdir("./audio", os.ModePerm)
	_ = os.Mkdir("./transcode", os.ModePerm)
}
