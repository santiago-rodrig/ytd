package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

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
		dirEntriesMap := make(map[string]struct{})
		for _, entry := range dirEntries {
			dirEntriesMap[entry.Name()] = struct{}{}
		}

		if err != nil {
			log.Fatal(err)
		}
		err = os.Chdir(dir)
		if err != nil {
			log.Fatal(err)
		}
		for url, fileName := range urls {
			isPresent := false
			for _, entry := range dirEntries {
				if !entry.IsDir() {
					if entry.Name() == fileName {
						isPresent = true
						break
					}
				}
			}
			if !isPresent {
				err = downloadSong(url)
				if err != nil {
					log.Fatal(err)
				}
				dirEntries, err := os.ReadDir(dir)
				if err != nil {
					log.Fatal(err)
				}
				for _, entry := range dirEntries {
					_, ok := dirEntriesMap[entry.Name()]
					if !ok {
						config[dir][url] = entry.Name()
						break
					}
				}
			}
		}
	}
	err = writeConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func writeConfig() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path.Join(homedir, ".config", "ytd", "ytd.toml"), os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := toml.NewEncoder(f)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return nil
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
	err = toml.Unmarshal(data, &config)
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
