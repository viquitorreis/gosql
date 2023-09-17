package main

import "fmt"

// "fmt"
// "os"

func main() {

	// readRes, err := ReadFile()
	// if err != nil {
	// 	fmt.Printf("Erro ao ler o arquivo %s", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(readRes)

	filenames, err := getDirFilenames()
	if err != nil {
		println(err)
		return
	}
	fmt.Println(len(filenames))
	fmt.Printf("%+v\n", filenames)

}
