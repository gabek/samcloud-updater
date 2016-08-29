package main

import (
	"bytes"
	_ "crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/bogem/id3v2"
)

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func setID3TagsForFile(filepath string, artist string, title string) {
	mp3File, err := id3v2.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer mp3File.Close()

	mp3File.SetArtist(artist)
	mp3File.SetTitle(title)

	mp3File.Save()
}

func htmlForURL(url string) string {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0")

	resp, _ := client.Do(req)
	bytes, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	return string(bytes)
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func TranscodeToMP3(originalFile string, destinationFile string) {
	args := []string{"-y", "-i", originalFile, "-q:a", "2", destinationFile}

	cmd := exec.Command("ffmpeg", args...)
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

func AddFileToUploadList(filename string) {
	uploadFile := "toUpload.m3u"

	f, err := os.OpenFile(uploadFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println(err)
	}
	n, err := io.WriteString(f, filename+"\n")
	if err != nil {
		fmt.Println(n, err)
	}
	f.Close()
}

func GenerateSlug(str string) (slug string) {
	return strings.Map(func(r rune) rune {
		switch {
		case r == ' ', r == '-':
			return '-'
		case r == '_', unicode.IsLetter(r), unicode.IsDigit(r):
			return r
		default:
			return -1
		}
		return -1
	}, strings.ToLower(strings.TrimSpace(str)))
}
