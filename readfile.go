package main

import (
	"fmt"
	"os"
	"strconv"
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

	filename := &FileName{}
	for _, file := range files {
		fileExt := file.Name()[len(file.Name())-4:]
		if fileExt == ".sql" {
			// fmt.Println(file.Name()[len(file.Name())-4:])
			filename.name = append(filename.name, file.Name())
		} else {
			return nil, fmt.Errorf("Somente arquivos .sql são suportados: %s", file.Name())
		}
	}

	return filename.name, nil
}

func createMigrationFile(cmds []string) error {
	// validação dos arquivos - prefix
	filenames, err := getDirFilenames()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	lastFile := filenames[len(filenames)-1]
	lfPrefix := lastFile[:4]
	intLfPrefix, err := strconv.Atoi(lfPrefix)
	if err != nil {
		fmt.Println("Os arquivos devem ter prefixo no padrão 0000: ", err)
		return nil
	}

	func() {
		newPrefix := fmt.Sprintf("%04d", intLfPrefix+1)
		filename := "./database/migrations/" + newPrefix + "." + cmds[1] // ARRUMAR => NOME VAI SER OQ VEM DPS DO "NEW"
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
		}

		defer func() {
			fmt.Println(FmtGreen("Arquivo criado com sucesso :D"))
			// fmt.Println(string("\033[34m"), "Arquivo criado com sucesso :D", string("\033[0m"))
			file.Close()
		}()
	}()

	fmt.Println(cmds)

	return nil
}

func GosqlCmd(cmds []string) {

	if cmds[0] != "gosql" {
		return
	}

	command := &Comands{}
	for i, cmd := range cmds {
		if i != 0 {
			command.cmds = append(command.cmds, cmd)
		}
	}

	handleGosqlCmds(command.cmds)
}

func handleGosqlCmds(cmds []string) { // criar interface para retornar funcoes aq
	for _, cmd := range cmds {
		if cmd == "new" {
			fmt.Println("test")
			createMigrationFile(cmds)
		}
	}
}
