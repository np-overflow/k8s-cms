/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * Archive utilties tests
*/
package utils

import (
	"os"
	"testing"
	"archive/tar"
	"compress/gzip"
	fp "path/filepath"
)


func TestTGZ(t *testing.T) {
	// create test files 
	testDir := fp.Join(os.TempDir(), "test")
	os.MkdirAll(fp.Join(testDir, "dir"), 0700)
	defer os.RemoveAll(testDir)
	file, _ := os.Create(fp.Join(testDir, "dir", "nested"))
	file.Close()

	// archive director
	tgzPath := fp.Join(os.TempDir(), "test.tgz")
	tgz := NewTGZ(tgzPath)
	tgz.ArchiveDir(testDir)
	tgz.Close()
	//defer os.Remove(tgzPath)

	// check contents of tar archive 
	file, _ = os.Open(tgzPath)
	defer file.Close()
	gr, _  := gzip.NewReader(file)
	tr := tar.NewReader(gr)
	header, _ := tr.Next()
	path := header.Name
	if path != fp.Join("dir", "nested") {
		t.Error("Contents of tar archive is not consistent in test")
	}
}

