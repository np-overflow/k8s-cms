/*
 * k8s-cms
 * kcmscli - k8s-cms comand line clien
 * Multipart form utilties
*/
package utils

import (
	"io"
	"os"
	"bytes"
	"strings"
	"mime/multipart"
)

// create multipart formdata with the given args
// fields - string fields to include in the multipart form data
// fileFields - file fields to include in the form data
// returns the content type and buffer with formdata
func NewMultipartData(fields map[string]string, fileFields map[string]*os.File) (
	string, *bytes.Buffer) {
	// setup to write formdata to buffer
	var formdata bytes.Buffer

	multipart := multipart.NewWriter(&formdata)
	defer multipart.Close()
	// write form data fields
	for name, data := range fields {
		dataWriter, err := multipart.CreateFormField(name)
		if err != nil {
			panic(err.Error())
		}
		io.Copy(dataWriter, strings.NewReader(data))
	}
	
	// write file fields 
	for name, file := range fileFields {
		fileWriter, err := multipart.CreateFormFile(name, file.Name())
		if err != nil {
			panic(err.Error())
		}
		io.Copy(fileWriter, file)
	}
	
	return multipart.FormDataContentType(), &formdata
}
