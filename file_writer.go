package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// writeArticlesToFile formats and writes all article data to a file
func writeArticlesToFile(articles []Article, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write header
	ist := time.FixedZone("IST", 5*3600+30*60)
	fmt.Fprintf(file, "LeetCode Discuss - Latest %d Articles\n", len(articles))
	fmt.Fprintf(file, "Fetched on: %s\n", time.Now().In(ist).Format("2006-01-02 15:04:05 MST"))
	fmt.Fprintf(file, "%s\n\n", strings.Repeat("=", 80))

	for i, article := range articles {
		fmt.Fprintf(file, "%s\n", strings.Repeat("═", 80))
		fmt.Fprintf(file, "Article #%d\n", i+1)
		fmt.Fprintf(file, "%s\n\n", strings.Repeat("═", 80))

		// Basic article info
		fmt.Fprintf(file, "UUID: %s\n", article.UUID)
		fmt.Fprintf(file, "Title: %s\n", article.Title)
		fmt.Fprintf(file, "Slug: %s\n", article.Slug)
		fmt.Fprintf(file, "Article Type: %s\n", article.ArticleType)
		fmt.Fprintf(file, "Posted: %s\n", formatStringTimestamp(article.CreatedAt))
		fmt.Fprintf(file, "Updated: %s\n", formatStringTimestamp(article.UpdatedAt))
		fmt.Fprintf(file, "URL: https://leetcode.com/discuss/post/%d/%s/\n", article.TopicId, article.Slug)
		fmt.Fprintf(file, "Author: %s\n", article.Author.UserName)

		// Summary
		if article.Summary != "" {
			fmt.Fprintf(file, "\n--- Summary ---\n")
			fmt.Fprintf(file, "%s\n", article.Summary)
		}

		// Tags
		if len(article.Tags) > 0 {
			fmt.Fprintf(file, "\n--- Tags ---\n")
			for _, tag := range article.Tags {
				fmt.Fprintf(file, "  - %s (%s) [%s]\n", tag.Name, tag.Slug, tag.TagType)
			}
		}

		// Reactions
		if len(article.Reactions) > 0 {
			fmt.Fprintf(file, "\n--- Reactions ---\n")
			for _, reaction := range article.Reactions {
				fmt.Fprintf(file, "  %s: %d\n", reaction.ReactionType, reaction.Count)
			}
		}

		fmt.Fprintf(file, "\n")
	}

	return nil
}
