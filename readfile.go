package main

import (
	"fmt"
	"os"
)

// func ReadFile(f *FileName) (string, error) {
// 	// filename := FileName{
// 	// 	name: os.Args[1],
// 	// }

// 	file, err := os.ReadFile(filename.name)
// 	if err != nil {
// 		return "", fmt.Errorf("erro ao ler o arquivo %s ", err)
// 	}

// 	// tem que converter o binário do arquivo para string
// 	fileContent := string(file)
// 	return fileContent, nil
// }

func getDirFilenames() ([]string, error) {
	files, err := os.ReadDir("./database/migrations")
	if err != nil {
		fmt.Println("Erro ao ler diretórios: ", err)
		return nil, nil
	}

	// filename := FileName{
	// 	name: []string{},
	// }
	filename := &FileName{}
	for _, file := range files {
		fileExt := file.Name()[len(file.Name())-4:]
		if fileExt == ".sql" {
			fmt.Println(file.Name()[len(file.Name())-4:])
			filename.name = append(filename.name, file.Name())
		} else {
			return nil, fmt.Errorf("Arquivo com extensão errada: %s", string(fileExt))
		}
	}

	fmt.Printf("%+v\n", filename)

	return filename.name, nil
}
