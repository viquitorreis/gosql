package main

import (
	"fmt"
	"os"
)

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
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", filenames)

	fmt.Println("args em main.go => ", os.Args[1:])

	GosqlCmd(os.Args[1:])

	fmt.Println("teste")
}
