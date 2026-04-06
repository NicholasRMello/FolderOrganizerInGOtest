package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var categories = map[string][]string{
	"Images":     {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp"},
	"Documents":  {".pdf", ".doc", ".docx", ".txt", ".xls", ".xlsx", ".ppt", ".pptx", ".odt"},
	"Videos":     {".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv"},
	"Compressed": {".zip", ".rar", ".7z", ".tar", ".gz"},
	"Audios":     {".mp3", ".wav", ".flac", ".ogg", ".m4a"},
}

func getCategory(ext string) string {
	ext = strings.ToLower(ext)
	for cat, exts := range categories {
		for _, e := range exts {
			if e == ext {
				return cat
			}
		}
	}
	return "Others"
}

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		fmt.Print("Enter the folder path: ")
		fmt.Scanln(&path)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	report := make(map[string]int)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		oldPath := filepath.Join(path, file.Name())
		ext := filepath.Ext(file.Name())
		category := getCategory(ext)

		newDir := filepath.Join(path, category)
		if _, err := os.Stat(newDir); os.IsNotExist(err) {
			err := os.MkdirAll(newDir, os.ModePerm)
			if err != nil {
				fmt.Printf("Error creating folder %s: %v\n", category, err)
				continue
			}
		}

		newPath := filepath.Join(newDir, file.Name())
		err := os.Rename(oldPath, newPath)
		if err != nil {
			// Fallback if Rename fails (e.g. across different disk volumes)
			if err := moveFile(oldPath, newPath); err != nil {
				fmt.Printf("Error moving %s: %v\n", file.Name(), err)
				continue
			}
		}

		report[category]++
	}

	fmt.Println("\n--- Organization Report ---")
	if len(report) == 0 {
		fmt.Println("No files organized.")
	} else {
		for cat, count := range report {
			fmt.Printf("%s: %d file(s)\n", cat, count)
		}
	}
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return err
	}

	inputFile.Close()
	return os.Remove(sourcePath)
}
