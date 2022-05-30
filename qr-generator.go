package main

import (
	"fmt"
)

func askFileName() string {
	fmt.Println("Geef de naam en extensie van jouw bestand [lijst.xlsx]")
	var fileName string
	fmt.Scan(&fileName)
	return fileName
}

func main() {
	fileName := askFileName()
	fmt.Println(fileName)
}
