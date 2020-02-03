package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	filesInfo := findSrcFilesInPackage(".")

	mutations := allMutations

	runMutationTests(mutations, filesInfo)
}

func runMutationTests(mutations []Mutation, files []os.FileInfo) {
	for _, m := range mutations {
		testSucceeded, err := testWithMutation(m, files)
		if err != nil {
			log.Println(err)
		}

		if testSucceeded {
			log.Printf("WARNING: Test succeeded with mutation: %v", m)
		} else {
			log.Printf("SUCCESS: Test failed with mutation: %v", m)
		}
	}
}

func testWithMutation(m Mutation, files []os.FileInfo) (bool, error) {
	for _, fileInfo := range files {
		createBackup(fileInfo)
		defer restoreBackup(fileInfo)

		err := applyMutation(m, fileInfo.Name())
		if err != nil {
			return false, fmt.Errorf("failed to apply mutation: %v", err)
		}
	}

	_, err := exec.Command("go", "test", "-count=1", ".").CombinedOutput()

	if err != nil {
		// log.Printf("tests failed, which is expected with mutation %v", m)
		return false, nil
	}

	return true, nil
}
