/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * utilties
*/
package main

import (
	"os"
	"fmt"
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
	file, err := os.OpenFile(path, os.O_WRONLY , 0644)
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
	
