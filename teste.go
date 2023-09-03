package main

import "github.com/Projeto-Pindorama/cbr2pdf/lib"
import "log"
import "os"

func main() {
	tempdir, _ := os.MkdirTemp("", "teste")
	err := lib.Unzip('x', os.Args[1], tempdir)
	if err != nil {
		log.Printf("%s\n", err)
	}
}
