package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func upload() {
	fmt.Println("Uploading...")

	args := []string{"SAM/ImportUtil", "-a", "refresh", "-f", uploadListFile, "-i", stationID, "-l", stationLogin, "-t", "MUS"}

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

	args := []string{"SAM/PlaylistUtil", "-l", stationLogin, "-i", stationID, "-f", uploadListFile, "-nd", "-p", stationPlaylist}

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

func samToolsCheck() {
	if !FileExists("SAM/PlaylistUtil") || !FileExists("SAM/ImportUtil") {
		log.Fatal("ERROR: SAM Broadcaster Import Utilities must be available in the SAM directory.")
	}
}
