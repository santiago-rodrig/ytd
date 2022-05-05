package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args[1:]) < 1 {
		log.Fatal("must provide a url")
	}
	url, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	_, err = exec.LookPath("youtube-dl")
	if err != nil {
		log.Fatal(err)
	}
	downloadCmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", "--no-playlist", "-o", "%(title)s.%(ext)s", url.String())
	downloadCmd.Stdout = os.Stdout
	err = downloadCmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
