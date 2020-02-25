/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * Archive utilties
*/
package utils

import (
	"io"
	"os"
	"strings"
	"path/filepath"
	"archive/tar"
	"compress/gzip"
)

type TGZ struct {
	file *os.File
	gzipWriter *gzip.Writer
	tarWriter *tar.Writer
}

// create & return a new tar gzipped archive at the given archive path. 
func NewTGZ(archivePath string) TGZ {
	// create file and writer used to write tgz
	// tar -> gzip -> file 
	file, err := os.Create(archivePath)
    if err != nil {
        panic(err.Error())
    }
	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)

	return TGZ {
		file,
		gzipWriter,
		tarWriter,
	}
}


// recursively archive the given directory into the tar archive
func (tgz *TGZ) ArchiveDir(dirPath string) {
	// check directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		panic(err.Error())
	}
	
	

	// collect paths & file info from directory
	var filePaths []string
	var headers []*tar.Header
	err := filepath.Walk(dirPath, func(filePath string,  info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip non regular files
		if !info.Mode().IsRegular() {
			return nil
		}
	
		filePaths = append(filePaths, filePath)
		
		// constructor tar header for file 
		fileArchivePath := strings.Replace(filePath, dirPath, "", 1)
		fileArchivePath = strings.TrimPrefix(
			fileArchivePath, string(filepath.Separator))
		header, err := tar.FileInfoHeader(info, filePath)
		if err != nil {
			panic(err.Error())
		}
		header.Name = fileArchivePath
		
		headers = append(headers, header)
		return err
	})

	// write header & file contents into tar writer
	for i := 0; i < len(filePaths); i++ {
		tgz.tarWriter.WriteHeader(headers[i])
	
		file, err := os.Open(filePaths[i])
		if err != nil {
			panic(err.Error())
		}
		io.Copy(tgz.tarWriter, file)
		file.Close()
	}
    if err != nil {
        panic(err.Error())
    }
}

// close all writers in the tgz file 
func (tgz *TGZ) Close() {
	tgz.tarWriter.Close()
	tgz.gzipWriter.Close()
	tgz.file.Close()
}

/* archive utils */
// make a gzipped tar archive of the given dir at the given path
// and writes the gzipped tar archive at the given path
// TODO: rewrite in go make this portable
func MakeTGZ(srcDir string, archivePath string) {
	tgz := NewTGZ(archivePath)
	tgz.ArchiveDir(srcDir)
	tgz.Close()
}
