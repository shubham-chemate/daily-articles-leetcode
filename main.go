package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	ist := time.FixedZone("IST", 5*3600+30*60)

	// Read configuration from environment variables
	sendGridAPIKey := strings.TrimSpace(os.Getenv("SENDGRID_API_KEY"))
	fromEmail := strings.TrimSpace(os.Getenv("FROM_EMAIL"))
	fromName := strings.TrimSpace(os.Getenv("FROM_NAME"))
	toEmailsStr := os.Getenv("TO_EMAILS")                          // Comma-separated list
	enableFileOutput := os.Getenv("ENABLE_FILE_OUTPUT") != "false" // Default to true

	// Parse recipient emails
	var toEmails []string
	if toEmailsStr != "" {
		toEmails = strings.Split(toEmailsStr, ",")
		for i := range toEmails {
			toEmails[i] = strings.TrimSpace(toEmails[i])
		}
	}

	// Validate configuration
	enableEmail := sendGridAPIKey != "" && fromEmail != "" && len(toEmails) > 0
	if !enableEmail && !enableFileOutput {
		fmt.Fprintf(os.Stderr, "Error: Either email or file output must be enabled\n")
		os.Exit(1)
	}

	// Read last processed timestamp from file
	lastProcessed, err := readLastProcessedTimestamp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading last processed timestamp: %v\n", err)
		os.Exit(1)
	}

	var cutoffTime time.Time
	if lastProcessed.IsZero() {
		// First run - fetch articles from last 24 hours
		cutoffTime = time.Now().Add(-24 * time.Hour)
		fmt.Println("First run - fetching articles from last 24 hours...")
	} else {
		cutoffTime = lastProcessed
		fmt.Printf("Last processed: %s\n", lastProcessed.In(ist).Format("2006-01-02 03:04 PM MST"))
	}

	fmt.Printf("Fetching articles published after %s...\n", cutoffTime.In(ist).Format("2006-01-02 03:04 PM MST"))

	// Fetch all articles after cutoff time using pagination
	articles, err := fetchArticlesAfterTime(cutoffTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching discuss articles: %v\n", err)
		os.Exit(1)
	}

	if len(articles) == 0 {
		fmt.Println("No new articles found.")
		return
	}

	fmt.Printf("Found %d articles published after cutoff time.\n", len(articles))

	// Print article summary
	for i, article := range articles {
		creationTime := formatStringTimestamp(article.CreatedAt)
		fmt.Printf("\n%d. %s\n", i+1, article.Title)
		fmt.Printf("   Created: %s\n", creationTime)
		fmt.Printf("   URL: https://leetcode.com/discuss/post/%d/%s/\n", article.TopicId, article.Slug)
	}

	// Send email if configured
	if enableEmail {
		fmt.Println("\nSending email...")
		subject := fmt.Sprintf("📚 LeetCode Daily Digest - %d New Articles", len(articles))
		htmlContent := generateHTMLEmail(articles, ist)

		if fromName == "" {
			fromName = "LeetCode Articles Bot"
		}

		err = sendEmailViaSendGrid(sendGridAPIKey, fromEmail, fromName, toEmails, subject, htmlContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending email: %v\n", err)
			// Don't exit, continue with file output if enabled
		} else {
			fmt.Printf("✓ Successfully sent email to: %s\n", strings.Join(toEmails, ", "))
		}
	}

	// Write to file if enabled
	if enableFileOutput {
		// Ensure fetched_articles directory exists
		if err := os.MkdirAll("fetched_articles", 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating fetched_articles directory: %v\n", err)
			os.Exit(1)
		}

		filename := fmt.Sprintf("fetched_articles/leetcode_articles_%s.txt", time.Now().In(ist).Format("2006-01-02_15-04-05"))
		err = writeArticlesToFile(articles, filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing articles to file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Successfully saved %d articles to %s\n", len(articles), filename)
	}

	// Update last processed timestamp with the most recent article
	if len(articles) > 0 {
		// Articles are sorted newest first, so the first one is the most recent
		newestTime, err := time.Parse(time.RFC3339, articles[0].CreatedAt)
		if err == nil {
			if err := writeLastProcessedTimestamp(newestTime); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to update last processed timestamp: %v\n", err)
			} else {
				fmt.Printf("Updated last processed timestamp to: %s\n", newestTime.In(ist).Format("2006-01-02 03:04 PM MST"))
			}
		}
	}
}
