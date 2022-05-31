package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	//"sync"

	"github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
)

// TODO POSSIBLE TO MAKE GOROUTINE?
func generateZIPFile(qrValues []string, saveName string) {
	currentTime := time.Now()
	formatTime := fmt.Sprintf(currentTime.Format("20060102150405"))

	saveName = saveName + "-" + formatTime + ".zip"

	qrZIP, err := os.Create("files/" + saveName)
	if err != nil {
		panic(err)
	}

	defer qrZIP.Close()
	zipWriter := zip.NewWriter(qrZIP)

	for _, qrValue := range qrValues {
		qrName := qrValue[len(qrValue)-10:] + ".png"
		var png []byte
		png, err := qrcode.Encode(qrValue, qrcode.Medium, 500)

		if err != nil {
			panic(err)
		}

		pngFile := bytes.NewReader(png)

		qrFile, err := zipWriter.Create(qrName)
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(qrFile, pngFile); err != nil {
			panic(err)
		}
	}
	zipWriter.Close()

	fmt.Printf("%v werd aangemaakt\n", saveName)
}

func indexOfColumn(columns []string, columnLetter string) int {
	for i := range columns {
		if columns[i] == columnLetter {
			return i
		}
	}
	return -1
}

func checkFormatValue(cellValue string) (string, bool) {
	matchDots, _ := regexp.MatchString("^[0-9]{4}.[0-9]{3}.[0-9]{3}$", cellValue)
	matchNoDots, _ := regexp.MatchString("^[0-9]{10}$", cellValue)

	if matchDots == true || matchNoDots == true {
		if matchNoDots == true {
			return cellValue, true
		} else {
			cellValueNoDots := strings.Split(cellValue, ".")
			cellValue = strings.Join(cellValueNoDots, "")
			return cellValue, true
		}
	}
	return cellValue, false
}

func selectColumn(sheetValues [][]string) []string {
	columns := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	var columnLetter string
	fmt.Println("In welke kolomletter staan de ondernemingsnummers: [A]")
	fmt.Scan(&columnLetter)
	index := indexOfColumn(columns, strings.ToUpper(columnLetter))
	//fmt.Println(index)

	// ask if there is a title column
	var hasTitle string
	for {
		fmt.Println("Heeft jouw kolom een titel? [y] [n]")
		fmt.Scan(&hasTitle)
		hasTitle = strings.ToLower(hasTitle)
		if hasTitle == "y" || hasTitle == "n" {
			break
		}
	}

	// 2 url options
	var chooseOption string
	for {
		fmt.Println("Onmiddellijk naar betaalflow [1] of de flow doorlopen [2]")
		fmt.Scan(&chooseOption)
		if chooseOption == "1" {
			chooseOption = "direct-payment"
			break
		} else if chooseOption == "2" {
			chooseOption = "member-flow"
			break
		}
	}

	// TODO: put this in seperate function?
	// TODO: regex check ondernemingsnummer format?
	var cellValues []string
	var wrongValues []string

	for i, cellValue := range sheetValues[index] {
		if i == 0 && hasTitle == "y" {
			continue
		}

		// check if the format is correct
		correctCellValue, isCorrect := checkFormatValue(cellValue)

		if isCorrect == false {
			correctCellValue = fmt.Sprintf("Rij %v - %v", i+1, correctCellValue)
			wrongValues = append(wrongValues, correctCellValue)
		} else {
			// generate URL
			correctCellValue = fmt.Sprintf("https://www.google.com/%v/%v", chooseOption, correctCellValue)
			cellValues = append(cellValues, correctCellValue)
		}

	}
	if len(wrongValues) > 0 {
		// create a logfile
		f, err := os.Create("files/log.txt")
		if err != nil {
			panic(err)
		}

		defer f.Close()

		for _, wrongValue := range wrongValues {
			_, err := f.WriteString(wrongValue + "\n")
			if err != nil {
				panic(err)
			}
		}
		var ignoreError string
		for {
			fmt.Printf("%v waarden (van de %v) staan in een foutief formaat.\nDoorgaan? [y] [n]\n", len(wrongValues), len(cellValues)+len(wrongValues))
			fmt.Scan(&ignoreError)
			ignoreError = strings.ToLower(ignoreError)
			if ignoreError == "y" {
				break
			} else if ignoreError == "n" {
				fmt.Println("Je kan de foute waarde in de log-file vinden.")
				os.Exit(1)
			}
		}
	}
	return cellValues
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

func openSourceFile(fileName string) [][]string {
	// Open File
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("Oops, er ging iets verkeerd")
		os.Exit(1)
	}

	// Close file when done
	defer f.Close()

	// get the name of the sheet
	sheetName := getSheetName(f)
	//fmt.Println(sheetName)

	// get all values of the sheet
	sheetValues := getValues(f, sheetName)
	//fmt.Println(sheetValues)

	return sheetValues
}

func fileNameNoExtension(fileName string) string {
	saveName := strings.Split(fileName, ".")
	return saveName[0]
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

	saveName := fileNameNoExtension(fileName)

	sheetValues := openSourceFile(fileName)

	// get the values of a specific column
	qrValues := selectColumn(sheetValues)

	generateZIPFile(qrValues, saveName)
}
