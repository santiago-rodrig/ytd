package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/exec"
)

var playlist = flag.Bool("playlist", false, "should the whole playlist be downloaded?")

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("must provide a url")
	}
	url, err := url.Parse(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	_, err = exec.LookPath("youtube-dl")
	if err != nil {
		log.Fatal(err)
	}
	var cmdArgs []string
	if *playlist {
		cmdArgs = []string{
			"-x",
			"--audio-format",
			"mp3",
			"-o",
			"%(title)s.%(ext)s",
			url.String(),
		}
	} else {
		cmdArgs = []string{
			"-x",
			"--audio-format",
			"mp3",
			"--no-playlist",
			"-o",
			"%(title)s.%(ext)s",
			url.String(),
		}
	}
	downloadCmd := exec.Command("youtube-dl", cmdArgs...)
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	err = downloadCmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
