package main

import (
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	regularFile = "not_test_file.go"
	testFile1   = "first_go_file_test.go"
	testFile2   = "second_go_file_test.go"
	fileMode   = os.FileMode(0o600)
)

func TestGenerateTestMatrix(t *testing.T) {
	tests := map[string]struct {
		setup    func(t *testing.T) string
		suite    string
		excluded []string
		want     GithubActionTestMatrix
		wantErr  bool
	}{
		"empty directory": {
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
		},
		"single test file": {
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeTestFile(t, dir, testFile1, "FeeMiddleware", "TestA", "TestB")
				return dir
			},
			want: GithubActionTestMatrix{
				Include: []TestSuitePair{
					{EntryPoint: "TestFeeMiddlewareTestSuite", Test: "TestA"},
					{EntryPoint: "TestFeeMiddlewareTestSuite", Test: "TestB"},
				},
			},
		},
		"multiple test files": {
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeTestFile(t, dir, testFile1, "FeeMiddleware", "TestA", "TestB")
				writeTestFile(t, dir, testFile2, "Transfer", "TestC", "TestD")
				return dir
			},
			want: GithubActionTestMatrix{
				Include: []TestSuitePair{
					{EntryPoint: "TestTransferTestSuite", Test: "TestC"},
					{EntryPoint: "TestFeeMiddlewareTestSuite", Test: "TestA"},
					{EntryPoint: "TestFeeMiddlewareTestSuite", Test: "TestB"},
					{EntryPoint: "TestTransferTestSuite", Test: "TestD"},
				},
			},
		},
		"non-test files ignored": {
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeTestFile(t, dir, regularFile, "FeeMiddleware", "TestA", "TestB")
				return dir
			},
			want: GithubActionTestMatrix{Include: []TestSuitePair{}},
		},
		"duplicate suites error": {
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				content := `package foo
func SuiteOne(t *testing.T) { suite.Run(t, new(FeeMiddlewareTestSuite)) }
func SuiteTwo(t *testing.T) { suite.Run(t, new(FeeMiddlewareTestSuite)) }
type FeeMiddlewareTestSuite struct {}`
				err := os.WriteFile(path.Join(dir, testFile1), []byte(content), fileMode)
				assert.NoError(t, err)
				return dir
			},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := tc.setup(t)
			got, err := generateTestMatrix(dir, tc.suite, tc.excluded)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assertMatrixEqual(t, tc.want, got)
		})
	}
}

func assertMatrixEqual(t *testing.T, want, got GithubActionTestMatrix) {
	sortMatrix := func(m *GithubActionTestMatrix) {
		sort.SliceStable(m.Include, func(i, j int) bool {
			if m.Include[i].EntryPoint == m.Include[j].EntryPoint {
				return m.Include[i].Test < m.Include[j].Test
			}
			return m.Include[i].EntryPoint < m.Include[j].EntryPoint
		})
	}

	sortMatrix(&want)
	sortMatrix(&got)
	assert.Equal(t, want.Include, got.Include)
}

func writeTestFile(t *testing.T, dir, filename, suiteName, test1, test2 string) {
	template := `package foo

func TestSUITETestSuite(t *testing.T) {
	suite.Run(t, new(SUITETestSuite))
}

type SUITETestSuite struct {}

func (s *SUITETestSuite) TEST1() {}
func (s *SUITETestSuite) TEST2() {}

func (s *SUITETestSuite) helper() {}
func helper() {}`

	content := strings.NewReplacer(
		"SUITE", suiteName,
		"TEST1", test1,
		"TEST2", test2,
	).Replace(template)

	err := os.WriteFile(path.Join(dir, filename), []byte(content), fileMode)
	assert.NoError(t, err)
}
