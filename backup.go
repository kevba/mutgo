package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var testFilePostfix = "_test.go"

func findSrcFilesInPackage(path string) []os.FileInfo {
	filesInfo := []os.FileInfo{}

	dir, _ := ioutil.ReadDir(path)
	for _, f := range dir {
		if !f.IsDir() {
			if strings.Contains(f.Name(), testFilePostfix) {
				continue
			}

			filesInfo = append(filesInfo, f)
		}
	}

	return filesInfo
}

func createBackup(file os.FileInfo) {
	f, _ := ioutil.ReadFile(file.Name())

	ioutil.WriteFile(getBackupName(file.Name()), f, file.Mode())
}

func restoreBackup(file os.FileInfo) {
	backupName := getBackupName(file.Name())

	f, _ := ioutil.ReadFile(backupName)
	ioutil.WriteFile(file.Name(), f, file.Mode())

	os.Remove(backupName)
}

func getBackupName(name string) string {
	return fmt.Sprintf("%v.og", name)
}
