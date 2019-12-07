package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
)

const testFile = "test/main.go"

var backupFile = fmt.Sprintf("%v.backup", testFile)

var testMutations = []Mutation{
	Mutation{token.ADD, token.SUB},
	Mutation{token.SUB, token.ADD},
}

type Mutation struct {
	token         token.Token
	mutationToken token.Token
}

func main() {
	err := os.Rename(testFile, backupFile)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = os.Rename(backupFile, testFile)
	}()

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, testFile, nil, 0600)

	if err != nil {
		log.Fatal(err)
	}

	srcFile, err := os.Create(testFile)
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(backupFile, srcFile)

	for _, mutation := range testMutations {
		ast.Inspect(astFile, func(x ast.Node) bool {
			s, ok := x.(*ast.BinaryExpr)
			if !ok {
				return true
			}

			if s.Op == mutation.token {
				s.Op = mutation.mutationToken
			}
			return false
		})

		if err := printer.Fprint(srcFile, fset, astFile); err != nil {
			log.Fatal(err)
		}

		output, err := exec.Command("go", "test").Output()
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(output))
	}

}

func createBackup(fileName, backupFileName string) *os.File {
	err := os.Rename(fileName, backupFileName)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(backupFileName)
	if err != nil {
		log.Fatal(err)
	}

	return f
}
