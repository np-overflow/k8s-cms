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
