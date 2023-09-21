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
		return nil, errors.New("Errors trying to read migrations directory")
	}

	filename := &FileName{}
	for _, file := range files {
		fileExt := file.Name()[len(file.Name())-4:]
		if fileExt == ".sql" {
			filename.name = append(filename.name, file.Name())
		} else {
			return nil, errors.New("Only .sql files are supported: " + file.Name())
		}
	}

	return filename.name, nil
}
func getMigrationsLastFile() (string, error) {
	filenames, err := getDirFilenames()
	if err != nil {
		return "", errors.New((FmtRed("Error trying to get last migration file") + err.Error()))
	}

	lastFile := filenames[len(filenames)-1]
	return lastFile, nil
}

func createMigrationFile(cmds []string) error {
	// validação dos arquivos - prefix
	lastFile, err := getMigrationsLastFile()
	if err != nil {
		return err
	}
	lfPrefix := lastFile[:4]
	intLfPrefix, err := strconv.Atoi(lfPrefix)
	if err != nil {
		return errors.New(("Migrations files must have a 0000 (NNNN) prefix pattern: " + err.Error()))
	}

	func() {
		newPrefix := fmt.Sprintf("%04d", intLfPrefix+1)
		filename := "./gosql/migrations/" + newPrefix + "_" + cmds[1] // ARRUMAR => NOME VAI SER OQ VEM DPS DO "NEW"
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(FmtRed("Errors trying to create file: " + err.Error()))
			return
		}

		defer func() {
			fmt.Println(FmtGreen("Migration file successfully created :D"))
			file.Close()
		}()
	}()

	return nil
}

func getFileByPrefix(prefix string) string {
	files, err := getDirFilenames()
	if err != nil {
		return err.Error()
	}
	for _, file := range files {
		if prefix == file[:4] {
			fmt.Println(file)
			return file
		}
	}
	return ""

}

func targetMigrationFile(cmd []string) error {
	fmt.Println("cmds in targetMigrationFile => ", cmd)

	if len(cmd) < 3 {
		filename, err := getMigrationsLastFile()
		if err != nil {
			return err
		}
		mgFile := readMigrationFile(filename)
		runDesiredMigrationSqlCmd(mgFile, cmd[1])

	} else if len(cmd) == 3 {
		filename := getFileByPrefix(cmd[2])
		mgFile := readMigrationFile(filename)
		runDesiredMigrationSqlCmd(mgFile, cmd[1])

	} else {
		return errors.New(FmtRed("More params passed than needed"))
	}

	return nil

}

func readMigrationFile(file string) []string {
	filename := "./gosql/migrations/" + file

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(FmtRed("Error trying to READ the migration file"), err)
		return []string{err.Error()} // ARRUMAR PRA RETORNAR O ERRO DPS -----------------------------------
	}
	dat := string(data)
	fileByLines := getFileByLines(dat)
	if err := validateFileLines(fileByLines); err != nil {
		fmt.Println(err)
		return []string{err.Error()} // ARRUMAR PRA RETORNAR O ERRO DPS -----------------------------------
	}

	return fileByLines

	// st, err := NewPostgresStore()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// mg := MigrationBridge{store: &PostgresStore{
	// 	db: st.db,
	// }}
	// mg.runMigration(fileByLines, cmd[0])

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

// FUNC runMigration está =>
// PEGANDO COMANDO UP NO ARQUIVO SQL
// PEGANDO COMANDO DOWN NO ARQUIVO SQL
// RODANDO MIGRATION UP
// RODANDO MIGRATION DOWN

func runDesiredMigrationSqlCmd(data []string, mtype string) error {

	st, err := NewPostgresStore()
	if err != nil {
		// fmt.Println(err)
		return err
	}
	mg := MigrationBridge{store: &PostgresStore{
		db: st.db,
	}}

	numLine, err := getCmdsLines(data)
	if err != nil {
		return err // ARRUMAR PRA RETORNAR O ERRO DPS -----------------------------------
	}
	if len(numLine) > 2 {
		return err // ARRUMAR PRA RETORNAR O ERRO DPS -----------------------------------
	}

	if mtype == "up" {
		upMigration := strings.Join(data[numLine[0]+1:numLine[1]-1], "")
		mg.runMigration(upMigration)
		return nil
	}

	if mtype == "down" {
		downMigration := strings.Join(data[numLine[1]+1:], "")
		mg.runMigration(downMigration)
		return nil
	}

	return errors.New(FmtRed("Wrong migration type") + mtype)

}

func (mg *MigrationBridge) runMigration(query string) error {

	mgb := &MigrationBody{
		query: query,
	}
	err := mg.store.RunMigration(mgb)
	if err != nil {
		return errors.New(FmtRed("Error trying to run up migration => ") + err.Error())
	}

	// if m == "up" {
	// 	mgb := &MigrationBody{
	// 		choice: m,
	// 		query:  upMigration,
	// 	}
	// 	err := mg.store.RunMigration(mgb)
	// 	if err != nil {
	// 		return errors.New(FmtRed("Error trying to run up migration => ") + err.Error())
	// 	}

	// 	return err
	// }

	// if m == "down" {
	// 	mgb := &MigrationBody{
	// 		choice: m,
	// 		query:  downMigration,
	// 	}
	// 	err := mg.store.RunMigration(mgb)
	// 	if err != nil {
	// 		return errors.New(FmtRed("Error trying to run down migration => ") + err.Error())
	// 	}

	// 	return err
	// }

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
	switch cmds[0] {
	case "new": // ESPECIFICAR O COMANDO "MIGRATION" APÓS O NEW OU SEI LA
		createMigrationFile(cmds)

	case "migration":
		handleMigrationCmd(cmds)

	default:
		fmt.Println(FmtRed("Command not found :(. Run gosql --help."))

	}

}

func handleMigrationCmd(cmds []string) {
	fmt.Println("called handle migration")
	for _, cmd := range cmds {
		switch cmd {
		case "up":
			targetMigrationFile(cmds)

		case "down":
			targetMigrationFile(cmds)

			// default:
			// 	fmt.Println(FmtRed("Type of migration is wrong :("), cmd)
		}
	}
}
