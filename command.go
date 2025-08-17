package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[string, flags]

func Comm(parameters ...any) yup.Command {
	return command(yup.Initialize[string, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Need two file paths to compare
		if len(p.Positional) < 2 {
			_, _ = fmt.Fprintf(stderr, "comm: missing operand\n")
			return fmt.Errorf("comm requires two files to compare")
		}

		file1Path := p.Positional[0]
		file2Path := p.Positional[1]

		// Read both files
		lines1, err := readFileLines(file1Path)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "comm: %s: %v\n", file1Path, err)
			return err
		}

		lines2, err := readFileLines(file2Path)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "comm: %s: %v\n", file2Path, err)
			return err
		}

		// Compare sorted files line by line
		i, j := 0, 0
		for i < len(lines1) || j < len(lines2) {
			var output string

			if i >= len(lines1) {
				// Only lines from file2 remain
				if !bool(p.Flags.SuppressColumn2) {
					if !bool(p.Flags.SuppressColumn1) {
						output = "\t"
					}
					output += lines2[j]
				}
				j++
			} else if j >= len(lines2) {
				// Only lines from file1 remain
				if !bool(p.Flags.SuppressColumn1) {
					output = lines1[i]
				}
				i++
			} else if lines1[i] < lines2[j] {
				// Line only in file1
				if !bool(p.Flags.SuppressColumn1) {
					output = lines1[i]
				}
				i++
			} else if lines1[i] > lines2[j] {
				// Line only in file2
				if !bool(p.Flags.SuppressColumn2) {
					if !bool(p.Flags.SuppressColumn1) {
						output = "\t"
					}
					output += lines2[j]
				}
				j++
			} else {
				// Line in both files
				if !bool(p.Flags.SuppressColumn3) {
					if !bool(p.Flags.SuppressColumn1) {
						output = "\t"
					}
					if !bool(p.Flags.SuppressColumn2) {
						output += "\t"
					}
					output += lines1[i]
				}
				i++
				j++
			}

			if output != "" {
				_, _ = fmt.Fprintln(stdout, output)
			}
		}

		return nil
	}
}

// readFileLines reads all lines from a file
func readFileLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
