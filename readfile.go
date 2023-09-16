package main

import (
	"fmt"
	"os"
)

func ReadFile(filename string) (string, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("erro ao ler o arquivo %s ", err)
	}

	// tem que converter o binário do arquivo para string
	fileContent := string(file)
	return fileContent, nil
}
