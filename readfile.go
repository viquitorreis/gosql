package main

import (
	"errors"
	"fmt"
	"os"
	_ "reflect"
	"strconv"
	"strings"
)

type MigrationBridge struct {
	store Storage
}

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
	files, err := os.ReadDir("./gosql/migrations")
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
		filename := "./gosql/migrations/" + newPrefix + "_" + cmds[1] // ARRUMAR => NOME VAI SER OQ VEM DPS DO "NEW"
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
		}

		defer func() {
			fmt.Println(FmtGreen("Arquivo criado com sucesso :D"))
			file.Close()
		}()
	}()

	// fmt.Println(cmds)

	return nil
}

func readMigrationFile(cmd string) {
	filename, err := getMigrationsLastFile()
	filename = "./gosql/migrations/" + filename
	if err != nil {
		fmt.Println(FmtRed("Error trying to GET the migration file"), err)
		return
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(FmtRed("Error trying to READ the migration file"), err)
		return
	}
	dat := string(data)
	fileByLines := getFileByLines(dat)
	if err := validateFileLines(fileByLines); err != nil {
		fmt.Println(err)
		return
	}

	// bridge := MigrationBridge{ storage: PostgresStore{
	// 	choice: "",
	// 	query: "",
	//  } }
	// fmt.Println("migration file called")
	// bridge := &MigrationBridge{}
	// br := bridge.store.RunMigration(<--- response from runMigration ---->)
	// b := bridge.(*MigrationBody)
	// bridge.runMigration(fileByLines, cmd)
	st := new(PostgresStore)
	mg := MigrationBridge{store: &PostgresStore{
		db: st.db,
	}}
	mg.runMigration(fileByLines, cmd)

}

func getFileByLines(data string) []string {
	return strings.Split(data, "\n")
}

func validateFileLines(data []string) error {
	// fmt.Println(len(data))
	if string(data[0][0]) != "-" || string(data[0][1]) != "-" {
		return errors.New(FmtRed("Files must start with --"))
	}

	return nil
}

// func getFileCommands(data []string) {

// }

func (mg *MigrationBridge) runMigration(data []string, m string) error {

	numLine, err := getCmdsLines(data)
	if err != nil {
		fmt.Println(err)
	}
	if len(numLine) > 2 {
		return err
	}

	upMigration := strings.Join(data[numLine[0]+1:numLine[1]-1], "")
	downMigration := strings.Join(data[numLine[1]+1:], "")

	if m == "up" {
		ch := &MigrationBody{
			choice: m,
			query:  upMigration,
		}
		err := mg.store.RunMigration(ch)
		if err != nil {
			return errors.New(FmtRed("Error trying to run up migration => ") + err.Error())
		}

		return err
	}

	if m == "down" {
		ch := &MigrationBody{
			choice: m,
			query:  downMigration,
		}
		err := mg.store.RunMigration(ch)
		if err != nil {
			return errors.New(FmtRed("Error trying to run down migration => ") + err.Error())
		}

		return err
	}

	return nil
}

func getCmdsLines(data []string) ([]int, error) {
	lines := []int{}
	for i, line := range data {
		if string(line) != "" {

			if i == 0 && len(line) >= 11 && line[:11] != "-- gosql Up" {
				return lines, errors.New(FmtRed("Migration files must start with '-- gosql Up'") + "\nReceived => '" + line + "'")
			} else if i == 0 && len(line) >= 11 && line[:11] == "-- gosql Up" {
				lines = append(lines, i)
			}

			if i != 0 && string(line[0]) == "-" && string(line[1]) == "-" {
				// fmt.Println(string(line))
				if i == 0 && line[:13] != "-- gosql Down" {
					return lines, errors.New(FmtRed("Migration down command must start with '-- gosql Down'") + "\nReceived => '" + line + "'")
				} else {
					lines = append(lines, i)
				}
			}

		}

	}

	if len(lines) > 2 {
		return lines, errors.New(FmtRed("Migration file have more than 2 command lines"))
	}

	if len(lines) < 2 {
		return lines, errors.New(FmtRed("Migration file have less than 2 command lines. Up and Down commands needed"))
	}

	return lines, nil
}

func GosqlCmd(cmds []string) {

	if cmds[0] != "gosql" {
		return
	}

	cmds = cmds[0:]
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
		case "new": // ESPECIFICAR O COMANDO "MIGRATION" APÓS O NEW OU SEI LA
			createMigrationFile(cmds)

		case "up":
			fmt.Println("called")
			readMigrationFile("up")

		case "down":
			readMigrationFile("down")

		default:
			fmt.Println(FmtRed("Command not found :("))

		}

	}
}
