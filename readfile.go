package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

/// VALIDAR EXTENSAO DOS ARQUIVOS NO DIRETORIO DE DESTINO SEPARADAMENTE

func getDirFilenames(reorder ...bool) ([]string, []int, error) {
	files, err := os.ReadDir("./gosql/migrations")
	if err != nil {
		return nil, nil, errors.New("Errors trying to read migrations directory")
	}

	filename := &FileName{}
	filesPrefix := []int{}
	for i, file := range files {
		fileExt := file.Name()[len(file.Name())-4:]
		intPrefix, _ := strconv.Atoi(file.Name()[:4])

		if i > 0 && len(reorder) == 0 && intPrefix-1 != filesPrefix[i-1] {
			fmt.Println("File isnt in sequential ordering =>", filesPrefix[i-1])
			log.Fatal(FmtRed("\nFiles must have a sequencial ordering name. If you want gosql to reorder all your files run 'gosql new reorder'")) // ============================== precisa retornar o erro
		}
		filesPrefix = append(filesPrefix, intPrefix)
		if fileExt == ".sql" {
			filename.name = append(filename.name, file.Name())
		} else {
			log.Fatal(errors.New(FmtRed("Only .sql files are supported: " + file.Name())))
		}
	}

	return filename.name, filesPrefix, nil

}

func getMigrationsLastFile() (string, error) {
	filenames, _, err := getDirFilenames()
	if err != nil {
		return "", err
		// return "", errors.New((FmtRed("Error trying to get last migration file") + err.Error()))
	}
	lastFile := filenames[len(filenames)-1]
	return lastFile, nil
}

func reorderUserFiles() {
	files, prefixes, _ := getDirFilenames(true)
	lowestWrongIdx := 0

	var rnm func(oldFile, newFile string) error
	rnm = func(oldFile, newFile string) error {
		err := os.Rename(oldFile, newFile)
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			fmt.Println(FmtGreen("File successfully renamed"))
			reorderUserFiles()
		}
		return nil
	}
	for i, file := range files {

		absFname, _ := filepath.Abs(file)
		idxBeginFile := strings.LastIndex(absFname, "/")
		newFileName := absFname[:idxBeginFile] + "/gosql/migrations/"
		oldFileName := absFname[:idxBeginFile] + "/gosql/migrations/" + fmt.Sprintf("%04d", prefixes[i]) + file[4:]

		if i == 0 {
			fmt.Println("file => ", file)
			if file[:4] != "0001" {
				fmt.Println(absFname)
				newFileName += fmt.Sprintf("%04d", 0001) + file[4:]
				rnm(oldFileName, newFileName)
				continue
			}
		}

		if i > 0 && prefixes[i] > 0 && prefixes[i] != prefixes[i-1]+1 && lowestWrongIdx < len(files) {
			fmt.Println("File is different", prefixes[i])

			fmt.Println("absFilename =>", absFname)
			newFileName += fmt.Sprintf("%04d", prefixes[i-1]+1) + file[4:]
			lowestWrongIdx = prefixes[i]
			fmt.Println("prefixes[i]=>", prefixes[i])
			fmt.Println("oldFileName =>", oldFileName)
			i = prefixes[i-1] + 1
			fmt.Println("nfName => ", newFileName)
			rnm(oldFileName, newFileName)
		}
	}
}

func createMigrationFile(cmds []string) error {
	// validação dos arquivos - prefix
	fmt.Println(cmds)
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
		filename := "./gosql/migrations/" + newPrefix + "." + cmds[1] + ".sql" // ARRUMAR => NOME VAI SER OQ VEM DPS DO "NEW"
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(FmtRed("Error trying to create file: " + err.Error()))
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
	files, _, err := getDirFilenames()
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
	if string(data[0][0]) != "-" || string(data[0][1]) != "-" {
		return errors.New(FmtRed("Files must start with --"))
	}

	return nil
}

func runDesiredMigrationSqlCmd(data []string, mtype string) error {

	st, err := NewPostgresStore()
	if err != nil {
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
	switch cmds[0] {
	case "new":
		handleNewCmd(cmds[1:])
	case "migration":
		handleMigrationCmd(cmds)

	default:
		fmt.Println(FmtRed("Command not found :(. Run gosql --help."))

	}

}

func handleNewCmd(cmds []string) {

	if len(cmds) == 0 {
		fmt.Println(FmtRed("Not enough parameters to run new command. Run 'gosql new --help'"))
		return
	}

	for _, cmd := range cmds {
		switch cmd {
		case "migration":
			createMigrationFile(cmds)
			return

		case "query":
			fmt.Println("Not yet implemented")

		case "reorder":
			reorderUserFiles()

		default:
			fmt.Println(FmtRed("Command not found. Run 'gosql new --help'\n")+"Received:", cmd)
		}

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
