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
		testResult, err := testWithMutation(m, files)
		if err != nil {
			log.Println(err)
		}

		if testResult.testSucceeded {
			log.Printf("WARNING: Test succeeded with mutation: %v", m)
		} else {
			log.Printf("SUCCESS: Test failed with mutation: %v", m)
		}
	}
}

func testWithMutation(m Mutation, files []os.FileInfo) (TestResult, error) {
	var totalChanges int
	result := TestResult{testSucceeded: false}

	for _, fileInfo := range files {
		createBackup(fileInfo)
		defer restoreBackup(fileInfo)

		changes, err := applyMutation(m, fileInfo.Name())
		if err != nil {
			return result, fmt.Errorf("failed to apply mutation: %v", err)
		}

		totalChanges += changes
	}

	if totalChanges == 0 {
		// The mutation had no effect, so the tests will still succeed.
		return result, nil
	}
	_, err := exec.Command("go", "test", "-count=1", ".").CombinedOutput()

	if err != nil {
		// log.Printf("tests failed, which is expected with mutation %v", m)
		return result, nil
	}
	result.testSucceeded = true

	return result, nil
}

type TestResult struct {
	testSucceeded bool
}
