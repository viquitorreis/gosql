package main

import (
	"fmt"
	"os"
)

func main() {

	filename := FileName{
		name: os.Args[1],
	}

	readRes, err := ReadFile(filename.name)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo %s", err)
		os.Exit(1)
	}

	fmt.Println(readRes)

}
