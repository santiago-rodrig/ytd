package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

var config map[string]map[string]string

func main() {
	flag.Parse()
	err := verifyCommands()
	if err != nil {
		log.Fatal(err)
	}
	err = readConfig()
	if err != nil {
		log.Fatal(err)
	}
	for dir, urls := range config {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Chdir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for url, fileName := range urls {
			for _, entry := range dirEntries {
				if !entry.IsDir() {
					if entry.Name() == fileName {
						continue
					}
					err = downloadSong(url)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

func readConfig() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(path.Join(homedir, ".config", "ytd", "ytd.toml"))
	if err != nil {
		return err
	}
	err = toml.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}

func verifyCommands() error {
	_, err := exec.LookPath("youtube-dl")
	if err != nil {
		return err
	}
	return nil
}

func readURLs(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func downloadSong(url string) error {
	downloadCmd := exec.Command("youtube-dl",
		"-x",
		"--audio-format",
		"mp3",
		"--no-playlist",
		"-o",
		"%(title)s.%(ext)s",
		url,
	)
	downloadCmd.Stdout = os.Stdout
	downloadCmd.Stderr = os.Stderr
	err := downloadCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
