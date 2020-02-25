/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * utilties
*/
package utils

import (
	"os"
	"bufio"
	"strings"
)

/* CSV Utilties */
// load & return a column of the given name in the csv file located as the given csv path.
// the file should be seperated with given seperator
// the csv file should have a first header line with column names
func LoadCSVColumn(csvPath string, colName string, seperator string) []string {
	// open csv 
	file, err := os.Open(csvPath)
	if err != nil {
		panic(err.Error())
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
		panic("loadCSVColumn: column not found in CSV file: " + colName)
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

