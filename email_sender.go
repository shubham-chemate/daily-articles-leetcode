package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const sendGridAPIURL = "https://api.sendgrid.com/v3/mail/send"

// SendGridEmail represents the email structure for SendGrid API
type SendGridEmail struct {
	Personalizations []Personalization `json:"personalizations"`
	From             EmailAddress      `json:"from"`
	Subject          string            `json:"subject"`
	Content          []Content         `json:"content"`
}

type Personalization struct {
	To []EmailAddress `json:"to"`
}

type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// sendEmailViaSendGrid sends an email using SendGrid API
func sendEmailViaSendGrid(apiKey, fromEmail, fromName string, toEmails []string, subject, htmlContent string) error {
	// Build recipient list
	var recipients []EmailAddress
	for _, email := range toEmails {
		recipients = append(recipients, EmailAddress{Email: email})
	}

	// Create email payload
	emailPayload := SendGridEmail{
		Personalizations: []Personalization{
			{To: recipients},
		},
		From: EmailAddress{
			Email: fromEmail,
			Name:  fromName,
		},
		Subject: subject,
		Content: []Content{
			{
				Type:  "text/html",
				Value: htmlContent,
			},
		},
	}

	jsonData, err := json.Marshal(emailPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", sendGridAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sendgrid API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// generateHTMLEmail creates an HTML email from articles
func generateHTMLEmail(articles []Article, ist *time.Location) string {
	var html strings.Builder

	html.WriteString(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 800px; margin: 0 auto; padding: 20px; background-color: #f5f5f5; }
        .container { background-color: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        h1 { color: #FFA116; border-bottom: 3px solid #FFA116; padding-bottom: 10px; margin-bottom: 20px; }
        .article { border-left: 4px solid #FFA116; padding: 15px; margin-bottom: 20px; background-color: #fafafa; border-radius: 4px; }
        .article-title { font-size: 18px; font-weight: bold; color: #262626; margin-bottom: 8px; }
        .article-title a { color: #262626; text-decoration: none; }
        .article-title a:hover { color: #FFA116; }
        .article-meta { font-size: 13px; color: #666; margin-bottom: 10px; }
        .article-summary { font-size: 14px; color: #555; line-height: 1.5; margin-bottom: 10px; }
        .article-tags { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 10px; }
        .tag { background-color: #e8f4f8; color: #0066cc; padding: 3px 10px; border-radius: 12px; font-size: 12px; }
        .reactions { font-size: 13px; color: #888; margin-top: 8px; }
        .footer { text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #ddd; color: #888; font-size: 12px; }
        .count { color: #FFA116; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <h1>📚 LeetCode Daily Articles</h1>
        <p>Found <span class="count">` + fmt.Sprintf("%d", len(articles)) + `</span> new articles:</p>
`)

	for i, article := range articles {
		html.WriteString(fmt.Sprintf(`
        <div class="article">
            <div class="article-title">%d. <a href="https://leetcode.com/discuss/%s">%s</a></div>
            <div class="article-meta">
                👤 %s | 📅 %s | 📝 %s
            </div>`,
			i+1,
			article.Slug,
			escapeHTML(article.Title),
			escapeHTML(article.Author.UserName),
			formatStringTimestamp(article.CreatedAt),
			article.ArticleType,
		))

		if article.Summary != "" {
			html.WriteString(fmt.Sprintf(`
            <div class="article-summary">%s</div>`,
				escapeHTML(truncateText(article.Summary, 200)),
			))
		}

		if len(article.Tags) > 0 {
			html.WriteString(`
            <div class="article-tags">`)
			for _, tag := range article.Tags {
				html.WriteString(fmt.Sprintf(`<span class="tag">%s</span>`, escapeHTML(tag.Name)))
			}
			html.WriteString(`</div>`)
		}

		if len(article.Reactions) > 0 {
			html.WriteString(`
            <div class="reactions">`)
			for j, reaction := range article.Reactions {
				if j > 0 {
					html.WriteString(" | ")
				}
				html.WriteString(fmt.Sprintf("%s: %d", reaction.ReactionType, reaction.Count))
			}
			html.WriteString(`</div>`)
		}

		html.WriteString(`
        </div>`)
	}

	html.WriteString(`
        <div class="footer">
            <p>Automated LeetCode Articles Digest | Generated on ` + time.Now().In(ist).Format("January 2, 2006 at 3:04 PM MST") + `</p>
        </div>
    </div>
</body>
</html>`)

	return html.String()
}

// escapeHTML escapes special HTML characters
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}

// truncateText truncates text to specified length with ellipsis
func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
