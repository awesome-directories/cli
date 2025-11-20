package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

var (
	colorsEnabled = true

	// Color schemes
	SuccessColor = color.New(color.FgGreen, color.Bold)
	ErrorColor   = color.New(color.FgRed, color.Bold)
	WarningColor = color.New(color.FgYellow, color.Bold)
	InfoColor    = color.New(color.FgCyan)
	MutedColor   = color.New(color.FgHiBlack)
	BoldColor    = color.New(color.Bold)

	// DR (Domain Rating) color thresholds
	HighDRColor   = color.New(color.FgGreen)
	MediumDRColor = color.New(color.FgYellow)
	LowDRColor    = color.New(color.FgRed)
)

// DisableColors disables colored output
func DisableColors() {
	colorsEnabled = false
	color.NoColor = true
}

// EnableColors enables colored output
func EnableColors() {
	colorsEnabled = true
	color.NoColor = false
}

// Success prints a success message
func Success(format string, args ...interface{}) {
	if colorsEnabled {
		if _, err := SuccessColor.Printf("✓ "+format+"\n", args...); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to print success message: %v\n", err)
		}
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// Error prints an error message
func Error(format string, args ...interface{}) {
	if colorsEnabled {
		if _, err := ErrorColor.Fprintf(os.Stderr, "✗ "+format+"\n", args...); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to print error message: %v\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}

// Warning prints a warning message
func Warning(format string, args ...interface{}) {
	if colorsEnabled {
		if _, err := WarningColor.Printf("⚠ "+format+"\n", args...); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to print warning message: %v\n", err)
		}
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// Info prints an info message
func Info(format string, args ...interface{}) {
	if colorsEnabled {
		if _, err := InfoColor.Printf("ℹ "+format+"\n", args...); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to print info message: %v\n", err)
		}
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// Muted prints a muted message
func Muted(format string, args ...interface{}) {
	if colorsEnabled {
		if _, err := MutedColor.Printf(format+"\n", args...); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to print muted message: %v\n", err)
		}
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// Bold prints a bold message
func Bold(format string, args ...interface{}) {
	if colorsEnabled {
		if _, err := BoldColor.Printf(format+"\n", args...); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to print bold message: %v\n", err)
		}
	} else {
		fmt.Printf(format+"\n", args...)
	}
}

// FormatDR formats a domain rating with color
func FormatDR(dr *int) string {
	if dr == nil {
		return MutedColor.Sprint("N/A")
	}

	value := *dr
	var colorFunc *color.Color

	switch {
	case value >= 70:
		colorFunc = HighDRColor
	case value >= 40:
		colorFunc = MediumDRColor
	default:
		colorFunc = LowDRColor
	}

	if colorsEnabled {
		return colorFunc.Sprint(value)
	}
	return strconv.Itoa(value)
}

// FormatPricing formats pricing type with color
func FormatPricing(pricing string) string {
	if !colorsEnabled {
		return pricing
	}

	switch strings.ToLower(pricing) {
	case "free":
		return HighDRColor.Sprint(pricing)
	case "freemium":
		return MediumDRColor.Sprint(pricing)
	case "paid":
		return LowDRColor.Sprint(pricing)
	default:
		return pricing
	}
}

// FormatLinkType formats link type with color
func FormatLinkType(linkType string) string {
	if !colorsEnabled {
		return linkType
	}

	switch strings.ToLower(linkType) {
	case "dofollow":
		return HighDRColor.Sprint(linkType)
	case "nofollow":
		return MutedColor.Sprint(linkType)
	default:
		return linkType
	}
}

// Table represents a simple table
type Table struct {
	writer  *tabwriter.Writer
	headers []string
	rows    [][]string
}

// CreateTable creates a formatted table
func CreateTable(headers []string) *Table {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	return &Table{
		writer:  w,
		headers: headers,
		rows:    [][]string{},
	}
}

// Row adds a row to the table
func (t *Table) Row(cols ...string) {
	t.rows = append(t.rows, cols)
}

// String renders the table
func (t *Table) String() string {
	if len(t.headers) > 0 {
		for i, h := range t.headers {
			if i > 0 {
				if _, err := fmt.Fprint(t.writer, "\t"); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to write tab: %v\n", err)
				}
			}
			if _, err := fmt.Fprint(t.writer, BoldColor.Sprint(h)); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write header: %v\n", err)
			}
		}
		_, err := fmt.Fprintln(t.writer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write newline: %v\n", err)
		}

		for i := range t.headers {
			if i > 0 {
				if _, err := fmt.Fprint(t.writer, "\t"); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to write tab: %v\n", err)
				}
			}
			if _, err := fmt.Fprint(t.writer, strings.Repeat("-", len(t.headers[i])+2)); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write separator: %v\n", err)
			}
		}
		_, err = fmt.Fprintln(t.writer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write newline: %v\n", err)
		}
	}

	for _, row := range t.rows {
		for i, col := range row {
			if i > 0 {
				if _, err := fmt.Fprint(t.writer, "\t"); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to write tab: %v\n", err)
				}
			}
			if _, err := fmt.Fprint(t.writer, col); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write column: %v\n", err)
			}
		}
		_, err := fmt.Fprintln(t.writer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write newline: %v\n", err)
		}
	}

	if err := t.writer.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush table writer: %v\n", err)
	}
	return ""
}

// TruncateString truncates a string to maxLen
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
