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

func getMigrationsLastFile() (string, error) {
	filenames, err := getDirFilenames()
	if err != nil {
		fmt.Println(FmtRed(err.Error()))
		return "", nil
	}

	lastFile := filenames[len(filenames)-1]
	return lastFile, nil
}

func createMigrationFile(cmds []string) error {
	// validação dos arquivos - prefix
	lastFile, err := getMigrationsLastFile()
	if err != nil {
		fmt.Println(FmtRed(err.Error()))
		return nil
	}

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
			file.Close()
		}()
	}()

	fmt.Println(cmds)

	return nil
}

func readMigrationFile() {
	fmt.Println("Called out just like that :b")

	filename, err := getMigrationsLastFile()
	filename = "./database/migrations/" + filename
	if err != nil {
		fmt.Println(FmtRed("Error trying to GET the migration file"), err)
		return
	}

	// f, err := os.Open("./database/migrations" + filename)
	// if err != nil {
	// 	fmt.Println(FmtRed("Error trying to OPENING the migration file"), err)
	// 	return
	// }
	// defer f.Close()

	// buf := make([]byte, 8)
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(FmtRed("Error trying to READ the migration file"), err)
		return
	}
	fmt.Println(string(data))
}

func GosqlCmd(cmds []string) {

	if cmds[0] != "gosql" {
		return
	}

	cmds = cmds[0:]
	// fmt.Println("cmds => ", len(cmds[1:]))
	if len(cmds[1:]) == 0 {
		fmt.Println(FmtRed("Gosql needs arguments in order to work"))
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
	fmt.Println("comands => ", cmds)
	for _, cmd := range cmds {
		switch cmd {
		case "new":
			createMigrationFile(cmds)

		case "up":
			readMigrationFile()

		default:
			fmt.Println(FmtRed("Command not found :("))

		}

	}
}
