package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func indexOfColumn(columns []string, columnLetter string) int {
	for i := range columns {
		if columns[i] == columnLetter {
			return i
		}
	}
	return -1
}

func selectColumn(sheetValues [][]string) []string {
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	var columnLetter string
	fmt.Println("In welke kolomletter staan de ondernemingsnummers: [A]")
	fmt.Scan(&columnLetter)
	index := indexOfColumn(columns, strings.ToUpper(columnLetter))
	//fmt.Println(index)

	var hasHead string
	fmt.Println("Heeft jouw kolom een titel? [y] [n]")
	fmt.Scan(&hasHead)

	// TODO: put this in sperate function?
	// TODO: regex check ondernemingsnummer format?
	var ondernemingsnummers []string

	for i, rowCell := range sheetValues[index] {
		if i == 0 && strings.ToLower(hasHead) == "y" {
			continue
		}
		ondernemingsnummers = append(ondernemingsnummers, rowCell)
	}
	return ondernemingsnummers
}

func getValues(f *excelize.File, sheetName string) [][]string {
	cols, err := f.GetCols(sheetName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cols
}

func getSheetName(f *excelize.File) string {
	return f.GetSheetName(0)
}

func askFileName() string {
	fmt.Println("Geef de naam en extensie van jouw bestand [lijst.xlsx]")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fileName := scanner.Text()
	return fileName
}

func main() {
	fileName := askFileName()
	fmt.Println(fileName)

	// Open File
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("Oops, er ging iets verkeerd")
		return
	}

	// Close file when done
	defer f.Close()

	// get the name of the sheet
	sheetName := getSheetName(f)
	//fmt.Println(sheetName)

	// get all values of the sheet
	sheetValues := getValues(f, sheetName)
	//fmt.Println(sheetValues)

	// get the values of a specific column
	columnValue := selectColumn(sheetValues)
	fmt.Println(columnValue)
}
