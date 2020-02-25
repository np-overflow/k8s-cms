/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * IO utilties
*/

package utils

import (
	"os"
	"io/ioutil"
)

/* I/O */
// read the file at the given path as []byte
// returns read bytes
func ReadBytes(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}
	
	return bytes
}

// write given bytes to the file at the given path
func WriteBytes(bytes []byte, path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE , 0644)
	if err != nil {
		panic(err.Error())
	}
	file.Truncate(0)
	file.Seek(0, 0) // whence 0: - with refrence to the start of the file
	
	_, err = file.Write(bytes)
	if err != nil {
		panic(err.Error())
	}
	
	file.Sync()
	file.Close()
}

// create a temp directory with the given prefix
// returns the path the temp directory
func MakeTempDir(prefix string) string {
	workDir, err := ioutil.TempDir(os.TempDir(), prefix)
	if err != nil {
		panic(err.Error())
	}
	
	return workDir
}

