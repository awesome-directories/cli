package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"

	"github.com/awesome-directories/cli/pkg/models"
)

// ExportToCSV exports directories to CSV format
func ExportToCSV(directories []models.Directory, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close CSV file")
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Name",
		"URL",
		"Description",
		"Categories",
		"Pricing",
		"Link Type",
		"Domain Rating",
		"Organic Traffic",
		"Organic Keywords",
		"Helpful Votes",
		"Submission URL",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write rows
	for _, dir := range directories {
		row := []string{
			dir.Name,
			dir.URL,
			dir.Description,
			strings.Join(dir.Categories, ", "),
			dir.Pricing,
			dir.LinkType,
			strconv.Itoa(dir.DomainRating),
			strconv.Itoa(dir.OrganicTraffic),
			strconv.Itoa(dir.OrganicKeywords),
			strconv.Itoa(dir.HelpfulCount),
			dir.SubmissionURL,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}

// ExportToJSON exports directories to JSON format
func ExportToJSON(directories []models.Directory, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close JSON file")
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(directories); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}

// ExportToMarkdown exports directories to Markdown format
func ExportToMarkdown(directories []models.Directory, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create Markdown file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Markdown file")
		}
	}()

	if _, err := fmt.Fprintf(file, "# Awesome Directories Export\n\n"); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := fmt.Fprintf(file, "Total directories: %d\n\n", len(directories)); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}
	if _, err := fmt.Fprintf(file, "---\n\n"); err != nil {
		return fmt.Errorf("failed to write separator: %w", err)
	}

	// Group by category
	categoryMap := make(map[string][]models.Directory)
	for _, dir := range directories {
		for _, cat := range dir.Categories {
			categoryMap[cat] = append(categoryMap[cat], dir)
		}
	}

	// Write by category
	for category, dirs := range categoryMap {
		if _, err := fmt.Fprintf(file, "## %s\n\n", category); err != nil {
			return fmt.Errorf("failed to write category: %w", err)
		}

		for _, dir := range dirs {
			if _, err := fmt.Fprintf(file, "### [%s](%s)\n\n", dir.Name, dir.URL); err != nil {
				return fmt.Errorf("failed to write directory name: %w", err)
			}
			if _, err := fmt.Fprintf(file, "%s\n\n", dir.Description); err != nil {
				return fmt.Errorf("failed to write description: %w", err)
			}

			if _, err := fmt.Fprintf(file, "- **Pricing:** %s\n", dir.Pricing); err != nil {
				return fmt.Errorf("failed to write pricing: %w", err)
			}
			if _, err := fmt.Fprintf(file, "- **Link Type:** %s\n", dir.LinkType); err != nil {
				return fmt.Errorf("failed to write link type: %w", err)
			}

			if dir.DomainRating > 0 {
				if _, err := fmt.Fprintf(file, "- **Domain Rating:** %d\n", dir.DomainRating); err != nil {
					return fmt.Errorf("failed to write domain rating: %w", err)
				}
			}

			if dir.HelpfulCount > 0 {
				if _, err := fmt.Fprintf(file, "- **Helpful Votes:** %d\n", dir.HelpfulCount); err != nil {
					return fmt.Errorf("failed to write helpful votes: %w", err)
				}
			}

			if dir.SubmissionURL != "" {
				if _, err := fmt.Fprintf(file, "- **Submission URL:** %s\n", dir.SubmissionURL); err != nil {
					return fmt.Errorf("failed to write submission URL: %w", err)
				}
			}

			if _, err := fmt.Fprintf(file, "\n"); err != nil {
				return fmt.Errorf("failed to write newline: %w", err)
			}
		}
	}

	return nil
}
