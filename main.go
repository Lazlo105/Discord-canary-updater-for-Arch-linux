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
	"strings"
)

// Константы для путей
const (
	tmpDir            = "/tmp"
	discordInstallDir = "/opt/discord-canary"
	downloadURL       = "https://discord.com/api/download/canary?platform=linux&format=tar.gz"
	archiveName       = "discord-canary.tar.gz"
)

// downloader скачивает файл по указанному URL и сохраняет его по заданному пути.
func downloader(url, filePath string) error {
	log.Println("Начинаем загрузку из Discord в", filePath)

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Println("Загрузка...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка статуса: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Println("Загрузка завершена")
	return nil
}

// extractor распаковывает tar.gz архив и возвращает путь к извлеченной папке.
func extractor(archivePath string) (string, error) {
	log.Println("Распаковка архива", archivePath, "...")

	file, err := os.Open(archivePath)
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
	destDir := filepath.Dir(archivePath)
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

		// Определяем имя первой папки в архиве (например, "DiscordCanary")
		if extractedFolderName == "" {
			// Разделяем путь по слэшам и берем первый элемент
			parts := strings.Split(header.Name, "/")
			if len(parts) > 0 && parts[0] != "" {
				extractedFolderName = parts[0]
			} else {
				extractedFolderName = header.Name
			}
		}

		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return "", err
			}

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}

			_, err = io.Copy(outFile, tr)
			outFile.Close() // Закрываем файл даже при ошибке копирования
			if err != nil {
				return "", err
			}
		}
	}

	log.Println("Распаковка завершена")
	return filepath.Join(destDir, extractedFolderName), nil
}

// updater переносит файлы из временной папки в целевую директорию установки.
func updater(sourceDir, destDir string) error {
	log.Println("Удаление старой версии...")

	if err := os.RemoveAll(destDir); err != nil {
		return err
	}
	log.Println("Старая версия удалена")

	log.Println("Установка новой версии...")

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(sourceDir, entry.Name())
		dstPath := filepath.Join(destDir, entry.Name())

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

	log.Println("Установка завершена")

	log.Println("Очистка временной папки...")
	if err := os.RemoveAll(sourceDir); err != nil {
		return err
	}

	return nil
}

// copyDir рекурсивно копирует содержимое директории.
func copyDir(src, dst string) error {
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

// copyFile копирует файл из одного места в другое с сохранением прав доступа.
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

	info, err := source.Stat()
	if err != nil {
		return err
	}

	if err := destination.Chmod(info.Mode()); err != nil {
		return err
	}

	return nil
}

func main() {
	// Формируем полные пути
	archivePath := filepath.Join(tmpDir, archiveName)

	// 1. Скачиваем архив во временную папку /tmp
	log.Println("Шаг 1: Загрузка архива")
	if err := downloader(downloadURL, archivePath); err != nil {
		log.Fatal("Ошибка загрузки:", err)
	}

	// 2. Распаковываем архив
	log.Println("Шаг 2: Распаковка архива")
	extractedPath, err := extractor(archivePath)
	if err != nil {
		log.Fatal("Ошибка распаковки:", err)
	}
	log.Println("Архив распакован в:", extractedPath)

	// 3. Переносим файлы в директорию установки
	log.Println("Шаг 3: Установка Discord Canary")
	if err := updater(extractedPath, discordInstallDir); err != nil {
		log.Fatal("Ошибка установки:", err)
	}

	// 4. Удаляем архив (файл в /tmp удалится сам после перезагрузки, но удалим сразу для чистоты)
	log.Println("Шаг 4: Удаление временного архива")
	if err := os.Remove(archivePath); err != nil {
		log.Printf("Предупреждение: не удалось удалить архив %s: %v", archivePath, err)
	}

	log.Println("Discord Canary успешно обновлен!")
}
