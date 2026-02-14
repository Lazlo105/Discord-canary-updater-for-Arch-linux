package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Downloader(url, filePath string) error {
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
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Println("Download complete")
	return nil
}

func Extracter(archive_path string) (string, error) {
	log.Println("Unpacking ", archive_path, " ...")

	file, err := os.Open(archive_path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	destDir := filepath.Dir(archive_path)
	var extractedFolderName string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		target := filepath.Join(destDir, header.Name)

		// В архиве Discord первый элемент обычно "DiscordCanary"
		if extractedFolderName == "" {
			parts := filepath.SplitList(header.Name)
			if len(parts) > 0 {
				extractedFolderName = parts[0]
			} else {
				extractedFolderName = header.Name
			}
		}

		if header.Typeflag == tar.TypeDir {
			os.MkdirAll(target, 0755)
		} else {
			os.MkdirAll(filepath.Dir(target), 0755)

			out, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			io.Copy(out, tr)
			out.Close()
		}
	}

	log.Println("Unpacking complete")

	return filepath.Join(destDir, extractedFolderName), nil
}

func Updater(source_dir, dist_dir string) error {
	log.Println("Deleting old version...")

	err := os.RemoveAll(dist_dir)
	if err != nil {
		return err
	}
	log.Println("Deleting complete")

	log.Println("Installing new version...")

	if err := os.MkdirAll(dist_dir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(source_dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(source_dir, entry.Name())
		dstPath := filepath.Join(dist_dir, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	log.Println("Installing complete")

	log.Println("Cleaning up temp source folder...")
	os.RemoveAll(source_dir)

	return nil
}

func copyDir(src string, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	info, _ := source.Stat()
	destination.Chmod(info.Mode())

	return nil
}

func main() {
	fileName := "discord-canary.tar.gz"
	fullPathToArchive := filepath.Join("/tmp", fileName)
	fullPathToDiscord := "/opt/discord-canary"

	// 1. Скачиваем
	err := Downloader("https://discord.com/api/download/canary?platform=linux&format=tar.gz", fullPathToArchive)
	if err != nil {
		log.Fatal("Download error:", err)
	}

	// 2. Распаковываем
	extractedPath, err := Extracter(fullPathToArchive)
	if err != nil {
		log.Fatal("Extract error:", err)
	}
	log.Println("Extracted to:", extractedPath)

	// 3. Переносим файлы
	err = Updater(extractedPath, fullPathToDiscord)
	if err != nil {
		log.Fatal("Update error:", err)
	}

	// 4. Удаляем архив
	os.Remove(fullPathToArchive)

	log.Println("Discord Canary updated")
}
