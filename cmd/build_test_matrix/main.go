package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type GithubActionTestMatrix struct {
	Include []TestSuitePair `json:"include"`
}

type TestSuitePair struct {
	Test       string `json:"test"`
	EntryPoint string `json:"entrypoint"`
}

const (
	testPrefix     = "Test"
	testFileSuffix = "_test.go"
	e2eDir         = "e2e"
	entryPointEnv  = "TEST_ENTRYPOINT"
	exclusionsEnv  = "TEST_EXCLUSIONS"
)

func main() {
	testFunc := os.Getenv(entryPointEnv)
	exclusions := strings.Split(os.Getenv(exclusionsEnv), ",")
	if len(exclusions) == 1 && exclusions[0] == "" {
		exclusions = nil
	}

	matrix, err := generateTestMatrix(e2eDir, testFunc, exclusions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating matrix: %s\n", err)
		os.Exit(1)
	}

	if output, err := json.Marshal(matrix); err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling json: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println(string(output))
	}
}

func generateTestMatrix(rootDir, suite string, excluded []string) (GithubActionTestMatrix, error) {
	testSuites := make(map[string][]string)
	fset := token.NewFileSet()

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !strings.HasSuffix(path, testFileSuffix) {
			return err
		}

		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parse error in %s: %w", path, err)
		}

		suiteName, tests, err := parseTestFile(f)
		if err != nil || suiteName == "" || isExcluded(suiteName, excluded) {
			return err
		}

		if suite == "" || suiteName == suite {
			testSuites[suiteName] = tests
		}

		return nil
	})

	if err != nil {
		return GithubActionTestMatrix{}, err
	}

	return buildMatrix(testSuites), nil
}

func parseTestFile(file *ast.File) (string, []string, error) {
	var suiteName string
	var testCases []string

	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			name := fd.Name.Name
			if !strings.HasPrefix(name, testPrefix) {
				return true
			}

			params := fd.Type.Params.List
			if len(params) == 1 {
				if suiteName == "" {
					suiteName = name
				}
			} else if len(params) == 0 {
				testCases = append(testCases, name)
			}
		}
		return true
	})

	if suiteName == "" {
		return "", nil, fmt.Errorf("no test suite found in %s", file.Name.Name)
	}

	return suiteName, testCases, nil
}

func isExcluded(name string, excluded []string) bool {
	for _, ex := range excluded {
		if ex == name {
			return true
		}
	}
	return false
}

func buildMatrix(testSuites map[string][]string) GithubActionTestMatrix {
	matrix := GithubActionTestMatrix{
		Include: make([]TestSuitePair, 0),
	}

	for suite, tests := range testSuites {
		for _, test := range tests {
			matrix.Include = append(matrix.Include, TestSuitePair{
				Test:       test,
				EntryPoint: suite,
			})
		}
	}

	return matrix
}
