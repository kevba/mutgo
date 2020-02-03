package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type Mutation struct {
	token         token.Token
	mutationToken token.Token
}

func (m Mutation) String() string {
	return fmt.Sprintf("%v to %v", m.token, m.mutationToken)
}

var allMutations = []Mutation{
	Mutation{token.ADD, token.SUB},
	Mutation{token.SUB, token.ADD},

	Mutation{token.ADD, token.MUL},
	Mutation{token.MUL, token.ADD},

	Mutation{token.SUB, token.MUL},
	Mutation{token.MUL, token.SUB},
}

func applyMutation(m Mutation, srcFileName string) error {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, srcFileName, nil, 0600)
	if err != nil {
		return fmt.Errorf("could not parse ast for source file %v: %v", srcFileName, err)
	}

	applyMutationToAST(m, astFile)

	srcFile, err := os.OpenFile(srcFileName, os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("could not open source file %v for mutating: %v", srcFileName, err)
	}

	if err := printer.Fprint(srcFile, fset, astFile); err != nil {
		return fmt.Errorf("could not mutate source file %v: %v", srcFileName, err)
	}

	return nil
}

func applyMutationToAST(m Mutation, astFile *ast.File) {
	ast.Inspect(astFile, func(x ast.Node) bool {
		s, ok := x.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		if s.Op == m.token {
			s.Op = m.mutationToken
		}
		return false
	})
}
