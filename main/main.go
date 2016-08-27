package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var config = getConfig()
var uploadListFile = "toUpload.m3u"
var stationLogin = "\"" + config.Station.Username + ";" + config.Station.Password + "\""
var stationID = strconv.Itoa(config.Station.Id)

func main() {
	fmt.Println("*** The Bat Station mix updater ***")

	os.Remove("toUpload.m3u")

	processPodcasts()
	processMixcloud()

	uploadListFile := "toUpload.m3u"
	if FileExists(uploadListFile) {
		upload()
		addToStation()
	}
}

func upload() {
	fmt.Println("Uploading...")

	args := []string{"SAM/ImportUtil", "-a", "refresh", "-f", uploadListFile, "-i", stationID, "-l", stationLogin, "-t", "MUS"}

	// fmt.Println(args)
	cmd := exec.Command("bash", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
}

func addToStation() {
	fmt.Println("Adding to station...")

	args := []string{"SAM/PlaylistUtil", "-l", stationLogin, "-i", stationID, "-f", uploadListFile, "-nd", "-p", "Mixes"}

	cmd := exec.Command("bash", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
}
