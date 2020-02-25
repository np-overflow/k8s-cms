/*
 * k8s-cms
 * kcmscli - k8s-cms comand line clien
 * Multipart form utilties
*/
package utils

import (
	"io"
	"os"
	"fmt"
	"bytes"
	"mime/multipart"
)

// ripped 
func CreateMultipartFormData(fieldName string , fileName string) ([]byte, *multipart.Writer) {
    var b bytes.Buffer
    var err error
    w := multipart.NewWriter(&b)
    var fw io.Writer
    file := mustOpen(fileName)
    if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		panic(err.Error())
    }
    if _, err = io.Copy(fw, file); err != nil {
		panic(err.Error())
    }
    w.Close()
    return b.Bytes(), w
}

func mustOpen(f string) *os.File {
    r, err := os.Open(f)
    if err != nil {
        pwd, _ := os.Getwd()
        fmt.Println("PWD: ", pwd)
        panic(err)
    }
    return r
}
