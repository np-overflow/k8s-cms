/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * Archive utilties
*/
package utils

import (
	"os"
	"os/exec"
)
	

/* archive utils */
// make a gzipped tar archive of the given dir at the given path
// and writes the gzipped tar archive at the given path
// TODO: rewrite in go make this portable
func MakeTGZ(srcDir string, archivePath string) {
	curdir, _ := os.Getwd()
	os.Chdir(srcDir)
	tarArgs := []string {"czf", archivePath, "." }
	tarCmd := exec.Command("tar", tarArgs...)
	_, err := tarCmd.Output()

    if err != nil {
        panic(err.Error())
    }
	
	os.Chdir(curdir)
}
