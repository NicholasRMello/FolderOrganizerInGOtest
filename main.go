package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var categories = map[string][]string{
	"Imagens":     {".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp"},
	"Documentos":  {".pdf", ".doc", ".docx", ".txt", ".xls", ".xlsx", ".ppt", ".pptx", ".odt"},
	"Vídeos":      {".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv"},
	"Compactados": {".zip", ".rar", ".7z", ".tar", ".gz"},
	"Áudios":      {".mp3", ".wav", ".flac", ".ogg", ".m4a"},
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
	return "Outros"
}

func main() {
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		fmt.Print("Informe o caminho da pasta: ")
		fmt.Scanln(&path)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Erro ao ler diretório: %v\n", err)
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
				fmt.Printf("Erro ao criar pasta %s: %v\n", category, err)
				continue
			}
		}

		newPath := filepath.Join(newDir, file.Name())
		err := os.Rename(oldPath, newPath)
		if err != nil {
			// Caso falhe o Rename (ex: volumes diferentes), tenta copiar e apagar
			if err := moveFile(oldPath, newPath); err != nil {
				fmt.Printf("Erro ao mover %s: %v\n", file.Name(), err)
				continue
			}
		}

		report[category]++
	}

	fmt.Println("\n--- Relatório de Organização ---")
	if len(report) == 0 {
		fmt.Println("Nenhum arquivo organizado.")
	} else {
		for cat, count := range report {
			fmt.Printf("%s: %d arquivo(s)\n", cat, count)
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

	inputFile.Close() // Fecha antes de remover
	return os.Remove(sourcePath)
}
