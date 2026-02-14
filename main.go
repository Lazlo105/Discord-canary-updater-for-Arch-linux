package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fileName := "discord-canary.tar.gz"
	fullPath := filepath.Join("/tmp", fileName)

	err := downloader("https://discord.com/api/download/canary?platform=linux&format=tar.gz", fullPath)
	if err != nil {
		log.Fatal(err)
	}
}

func downloader(url, filePath string) error {
	log.Println("Init download from discord.com to", filePath)

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Println("Downloading...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(string(resp.Status))
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Println("Download complete")
	return nil
}
