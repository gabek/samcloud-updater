package main

import (
	"bytes"
	"bufio"
	_ "crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"unicode"
	"encoding/json"
	"time"
)

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	userAgent := getConfig().UserAgent	
	req.Header.Set("User-Agent", userAgent)

	resp, _ := client.Do(req)

	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func htmlForURL(url string) string {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	userAgent := getConfig().UserAgent	
	req.Header.Set("User-Agent", userAgent)

	resp, _ := client.Do(req)
	bytes, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	return string(bytes)
}

func getJson(url string, target interface{}) error {
	client := &http.Client{}

    r, err := client.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func HasPreviouslyDownloaded(name string) bool {
	filenameCache := getConfig().DownloadLog

	if !FileExists(filenameCache) {
		return false
	}

	f, err := ioutil.ReadFile(filenameCache)
	if err != nil {
        fmt.Print(err)
    }

	stringData := string(f)

	return strings.Contains(stringData, name)
}

func MarkFileAsDownloaded(name string) {
	filenameCache := getConfig().DownloadLog

	file, err := os.OpenFile(filenameCache, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Print(err)
	}
	defer file.Close()

	timestampString := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("%s\t%s", timestampString, name)
	w := bufio.NewWriter(file)
	fmt.Fprintln(w, logLine)
	w.Flush()
}

func TranscodeToMP3(originalFile string, destinationFile string, artist string, title string) {
	args := []string{"-i", originalFile, "-acodec", "libmp3lame", "-ab", "128k", "-metadata", "artist=" + artist, "-metadata", "title=" + title, "-y", destinationFile}

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

func TranscodeHLSToMp3(url string, destinationFile string, artist string, title string) {
	userAgentHeader := "User-Agent: " + getConfig().UserAgent	

	args := []string{"-headers", userAgentHeader, "-i", url, "-c", "copy", "-acodec", "libmp3lame", "-ab", "128k", "-metadata", "artist=" + artist, "-metadata", "title=" + title, "-y", destinationFile}
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
	}, strings.ToLower(strings.TrimSpace(str)))
}
