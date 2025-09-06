package comm

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"

	localopt "github.com/yupsh/comm/opt"
)

// Flags represents the configuration options for the comm command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Comm creates a new comm command with the given parameters
func Comm(parameters ...any) yup.Command {
	return command(opt.Args[string, Flags](parameters...))
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	// Check for cancellation before starting
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if len(c.Positional) != 2 {
		fmt.Fprintln(stderr, "comm: need exactly 2 files")
		return fmt.Errorf("need exactly 2 files")
	}

	file1Name := c.Positional[0]
	file2Name := c.Positional[1]

	// Open files
	var file1, file2 io.ReadCloser
	var err error

	if file1Name == "-" {
		file1 = io.NopCloser(stdin)
	} else {
		file1, err = os.Open(file1Name)
		if err != nil {
			fmt.Fprintf(stderr, "comm: %s: %v\n", file1Name, err)
			return err
		}
		defer file1.Close()
	}

	// Check for cancellation after opening first file
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	if file2Name == "-" {
		file2 = io.NopCloser(stdin)
	} else {
		file2, err = os.Open(file2Name)
		if err != nil {
			fmt.Fprintf(stderr, "comm: %s: %v\n", file2Name, err)
			return err
		}
		defer file2.Close()
	}

	return c.compareFiles(ctx, file1, file2, stdout, stderr)
}

func (c command) compareFiles(ctx context.Context, file1, file2 io.Reader, output, stderr io.Writer) error {
	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)

	var line1, line2 string
	var hasLine1, hasLine2 bool
	var counts [3]int // [unique to file1, unique to file2, common]
	lineCount := 0

	// Read first lines with context awareness
	if yup.ScanWithContext(ctx, scanner1) {
		line1 = scanner1.Text()
		hasLine1 = true
	}
	if yup.ScanWithContext(ctx, scanner2) {
		line2 = scanner2.Text()
		hasLine2 = true
	}

	for hasLine1 || hasLine2 {
		// Check for cancellation periodically (every 1000 lines for efficiency)
		lineCount++
		if lineCount%1000 == 0 {
			if err := yup.CheckContextCancellation(ctx); err != nil {
				return err
			}
		}

		if !hasLine1 {
			// Only file2 has remaining lines
			if !bool(c.Flags.SuppressColumn2) {
				c.outputLine(output, line2, 2)
			}
			counts[1]++
			if yup.ScanWithContext(ctx, scanner2) {
				line2 = scanner2.Text()
			} else {
				hasLine2 = false
			}
		} else if !hasLine2 {
			// Only file1 has remaining lines
			if !bool(c.Flags.SuppressColumn1) {
				c.outputLine(output, line1, 1)
			}
			counts[0]++
			if yup.ScanWithContext(ctx, scanner1) {
				line1 = scanner1.Text()
			} else {
				hasLine1 = false
			}
		} else {
			// Both files have lines - compare them
			cmp := strings.Compare(line1, line2)
			if cmp < 0 {
				// line1 < line2: line1 is unique to file1
				if !bool(c.Flags.SuppressColumn1) {
					c.outputLine(output, line1, 1)
				}
				counts[0]++
				if yup.ScanWithContext(ctx, scanner1) {
					line1 = scanner1.Text()
				} else {
					hasLine1 = false
				}
			} else if cmp > 0 {
				// line1 > line2: line2 is unique to file2
				if !bool(c.Flags.SuppressColumn2) {
					c.outputLine(output, line2, 2)
				}
				counts[1]++
				if yup.ScanWithContext(ctx, scanner2) {
					line2 = scanner2.Text()
				} else {
					hasLine2 = false
				}
			} else {
				// line1 == line2: common line
				if !bool(c.Flags.SuppressColumn3) {
					c.outputLine(output, line1, 3)
				}
				counts[2]++
				// Advance both
				if yup.ScanWithContext(ctx, scanner1) {
					line1 = scanner1.Text()
				} else {
					hasLine1 = false
				}
				if yup.ScanWithContext(ctx, scanner2) {
					line2 = scanner2.Text()
				} else {
					hasLine2 = false
				}
			}
		}
	}

	// Check if context was cancelled
	if err := yup.CheckContextCancellation(ctx); err != nil {
		return err
	}

	// Check for scan errors
	if err := scanner1.Err(); err != nil {
		fmt.Fprintf(stderr, "comm: error reading file1: %v\n", err)
		return err
	}
	if err := scanner2.Err(); err != nil {
		fmt.Fprintf(stderr, "comm: error reading file2: %v\n", err)
		return err
	}

	// Output totals if requested
	if bool(c.Flags.Total) {
		fmt.Fprintf(output, "%d\t%d\t%d\ttotal\n", counts[0], counts[1], counts[2])
	}

	return nil
}

func (c command) outputLine(output io.Writer, line string, column int) {
	switch column {
	case 1:
		// Column 1: unique to file1
		fmt.Fprintln(output, line)
	case 2:
		// Column 2: unique to file2 (indented)
		fmt.Fprintf(output, "\t%s\n", line)
	case 3:
		// Column 3: common to both (double indented)
		fmt.Fprintf(output, "\t\t%s\n", line)
	}
}

func (c command) String() string {
	return fmt.Sprintf("comm %v", c.Positional)
}
