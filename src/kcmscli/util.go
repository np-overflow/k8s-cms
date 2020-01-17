/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * utilties
*/
package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"io/ioutil"
)

/* utilties */

// kill program due to a fatal error detailed by given messagea
func die(message string) {
	fmt.Printf("FATAL: %s\n", message)
	os.Exit(1)
}

/* I/O */
// read the file at the given path as []byte
// returns read bytes
func readBytes(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		die(err.Error())
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		die(err.Error())
	}
	
	return bytes
}

// write given bytes to the file at the given path
func writeBytes(bytes []byte, path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE , 0644)
	if err != nil {
		die(err.Error())
	}
	file.Truncate(0)
	file.Seek(0, 0) // whence 0: - with refrence to the start of the file
	
	_, err = file.Write(bytes)
	if err != nil {
		die(err.Error())
	}
	
	file.Sync()
	file.Close()
}

/* CSV Utilties */
// load & return a column of the given name in the csv file located as the given csv path.
// the file should be seperated with given seperator
// the csv file should have a first header line with column names
func loadCSVColumn(csvPath string, colName string, seperator string) []string {
	// open csv 
	file, err := os.Open(csvPath)
	if err != nil {
		die(err.Error())
	}

	// read header line
	lineScanner := bufio.NewScanner(file)
	lineScanner.Scan() 
	header := lineScanner.Text()
	
	// find position colName of in the header line
	colIdx := -1
	colNames := strings.Split(header, seperator)
	for i, c := range colNames {
		if colName == c {
			colIdx = i
			break
		}
	}
	
	if colIdx == -1 {
		die("loadCSVColumn: column not found in CSV file: " + colName)
	}

	// collect values for the specified column
	var colValues []string
	for lineScanner.Scan() {
		line := lineScanner.Text()
		rowValues := strings.Split(line, seperator)
		colValues = append(colValues, rowValues[colIdx])
	}
	
	return colValues
}
