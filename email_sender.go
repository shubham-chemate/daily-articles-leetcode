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
        body { font-family: Georgia, 'Times New Roman', serif; line-height: 1.8; color: #333; max-width: 680px; margin: 0 auto; padding: 40px 20px; background-color: #fff; }
        h1 { font-size: 28px; font-weight: normal; color: #222; margin-bottom: 10px; letter-spacing: -0.5px; }
        .subtitle { color: #666; font-size: 14px; margin-bottom: 40px; }
        .article { margin-bottom: 40px; padding-bottom: 30px; border-bottom: 1px solid #e5e5e5; }
        .article:last-child { border-bottom: none; }
        .article-title { font-size: 20px; font-weight: 600; margin-bottom: 8px; line-height: 1.4; }
        .article-title a { color: #222; text-decoration: none; }
        .article-title a:hover { color: #0066cc; }
        .article-meta { font-size: 13px; color: #888; margin-bottom: 12px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; }
        .article-summary { font-size: 15px; color: #444; line-height: 1.7; margin-bottom: 12px; }
        .article-tags { margin-top: 12px; }
        .tag { display: inline; color: #666; font-size: 13px; margin-right: 12px; }
        .tag:before { content: "#"; color: #999; }
        .footer { text-align: center; margin-top: 50px; padding-top: 20px; border-top: 1px solid #e5e5e5; color: #999; font-size: 12px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Arial, sans-serif; }
    </style>
</head>
<body>
    <h1>LeetCode Daily Digest</h1>
    <div class="subtitle">` + fmt.Sprintf("%d new articles • %s", len(articles), time.Now().In(ist).Format("January 2, 2006")) + `</div>
`)

	for _, article := range articles {
		html.WriteString(fmt.Sprintf(`
    <div class="article">
        <div class="article-title"><a href="https://leetcode.com/discuss/post/%d/%s/">%s</a></div>
        <div class="article-meta">By %s • %s</div>`,
			article.TopicId,
			article.Slug,
			escapeHTML(article.Title),
			escapeHTML(article.Author.UserName),
			formatStringTimestamp(article.CreatedAt),
		))

		if article.Summary != "" {
			html.WriteString(fmt.Sprintf(`
        <div class="article-summary">%s</div>`,
				escapeHTML(truncateText(article.Summary, 250)),
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

		html.WriteString(`
    </div>`)
	}

	html.WriteString(`
    <div class="footer">
        <p>Automated digest • LeetCode Articles Fetcher</p>
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
