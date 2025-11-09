package command_test

import (
	"testing"

	"github.com/gloo-foo/testable/assertion"
	"github.com/gloo-foo/testable/run"
	command "github.com/yupsh/comm"
)

func TestComm_Basic(t *testing.T) {
	result := run.Quick(command.Comm("testdata/file1.txt", "testdata/file2.txt"))
	assertion.NoError(t, result.Err)
	// Should show: lines only in file1, lines only in file2, lines in both
	assertion.Count(t, result.Stdout, 8) // Actual output count
}

func TestComm_Suppress1(t *testing.T) {
	// Suppress column 1 (lines only in file1)
	result := run.Quick(command.Comm("testdata/file1.txt", "testdata/file2.txt", command.SuppressColumn1))
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5) // Only file2-only and common lines
}

func TestComm_SuppressColumn2(t *testing.T) {
	// Suppress column 2 (lines only in file2)
	result := run.Quick(command.Comm("testdata/file1.txt", "testdata/file2.txt", command.SuppressColumn2))
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5) // Only file1-only and common lines
}

func TestComm_SuppressColumn3(t *testing.T) {
	// Suppress column 3 (lines in both files)
	result := run.Quick(command.Comm("testdata/file1.txt", "testdata/file2.txt", command.SuppressColumn3))
	assertion.NoError(t, result.Err)
	assertion.Count(t, result.Stdout, 5) // Only unique lines from each file
}

func TestComm_MissingFile(t *testing.T) {
	result := run.Quick(command.Comm("nonexistent.txt", "testdata/file2.txt"))
	assertion.Error(t, result.Err)
}

