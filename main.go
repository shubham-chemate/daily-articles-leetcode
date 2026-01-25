package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Article represents an article from LeetCode discuss section
type Article struct {
	UUID        string     `json:"uuid"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Summary     string     `json:"summary"`
	Author      Author     `json:"author"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
	HitCount    int        `json:"hitCount"`
	Thumbnail   string     `json:"thumbnail"`
	Tags        []Tag      `json:"tags"`
	Topic       *Topic     `json:"topic"`
	Reactions   []Reaction `json:"reactions"`
	ArticleType string     `json:"articleType"`
	Status      string     `json:"status"`
	TopicID     int        `json:"topicId"`
}

// Author represents the article author
type Author struct {
	UserName           string  `json:"userName"`
	UserSlug           string  `json:"userSlug"`
	RealName           string  `json:"realName"`
	UserAvatar         string  `json:"userAvatar"`
	NameColor          *string `json:"nameColor"`
	CertificationLevel string  `json:"certificationLevel"`
	ActiveBadge        *Badge  `json:"activeBadge"`
}

// Badge represents user badges
type Badge struct {
	DisplayName string `json:"displayName"`
	Icon        string `json:"icon"`
}

// Topic represents discussion topic info
type Topic struct {
	ID                   int `json:"id"`
	TopLevelCommentCount int `json:"topLevelCommentCount"`
}

// Reaction represents user reactions to article
type Reaction struct {
	Count        int    `json:"count"`
	ReactionType string `json:"reactionType"`
}

// Tag represents article tags
type Tag struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	TagType string `json:"tagType"`
}

// ArticlesResponse represents the GraphQL response for articles
type ArticlesResponse struct {
	Data struct {
		UgcArticleDiscussionArticles struct {
			TotalNum int `json:"totalNum"`
			Edges    []struct {
				Node Article `json:"node"`
			} `json:"edges"`
		} `json:"ugcArticleDiscussionArticles"`
	} `json:"data"`
}

const (
	leetcodeGraphQLURL = "https://leetcode.com/graphql"
	discussTopicsQuery = `
		query discussPostItems($orderBy: ArticleOrderByEnum, $keywords: [String]!, $tagSlugs: [String!], $skip: Int, $first: Int) {
			ugcArticleDiscussionArticles(
				orderBy: $orderBy
				keywords: $keywords
				tagSlugs: $tagSlugs
				skip: $skip
				first: $first
			) {
				totalNum
				edges {
					node {
						uuid
						title
						slug
						summary
						author {
							realName
							userAvatar
							userSlug
							userName
							nameColor
							certificationLevel
							activeBadge {
								icon
								displayName
							}
						}
						articleType
						thumbnail
						summary
						createdAt
						updatedAt
						status
						topicId
						hitCount
						reactions {
							count
							reactionType
						}
						tags {
							name
							slug
							tagType
						}
						topic {
							id
							topLevelCommentCount
						}
					}
				}
			}
		}
	`
)

func main() {
	fmt.Println("Fetching Latest 10 Articles from LeetCode Discuss...")

	articles, err := fetchDiscussArticles(10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching discuss articles: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d articles:\n", len(articles))

	for i, article := range articles {
		creationTime := formatStringTimestamp(article.CreatedAt)
		fmt.Printf("\n%d. %s\n", i+1, article.Title)
		fmt.Printf("   Created: %s\n", creationTime)
		fmt.Printf("   URL: https://leetcode.com/discuss/%s/%s\n", article.ArticleType, article.Slug)
	}

	filename := fmt.Sprintf("leetcode_articles_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	err = writeArticlesToFile(articles, filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing articles to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Successfully saved %d articles to %s\n", len(articles), filename)
}

// fetchDiscussArticles fetches the most recent articles from LeetCode discuss section
func fetchDiscussArticles(count int) ([]Article, error) {
	reqBody := map[string]interface{}{
		"query": discussTopicsQuery,
		"variables": map[string]interface{}{
			"orderBy":  "MOST_RECENT",
			"keywords": []string{},
			"tagSlugs": []string{},
			"skip":     0,
			"first":    count,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("POST", leetcodeGraphQLURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "LeetCode-Discuss-Fetcher/1.0")
	req.Header.Set("Referer", "https://leetcode.com/discuss/")
	req.Header.Set("Origin", "https://leetcode.com")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var articlesResp ArticlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&articlesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract articles from edges
	var articles []Article
	for _, edge := range articlesResp.Data.UgcArticleDiscussionArticles.Edges {
		articles = append(articles, edge.Node)
	}

	// Articles are already sorted by NEWEST, no need to sort again
	return articles, nil
}

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
		fmt.Fprintf(file, "Posted: %s\n", formatStringTimestamp(article.CreatedAt))
		fmt.Fprintf(file, "Updated: %s\n", formatStringTimestamp(article.UpdatedAt))
		fmt.Fprintf(file, "URL: https://leetcode.com/discuss/%s/%s\n", article.ArticleType, article.Slug)
		fmt.Fprintf(file, "Article Type: %s\n", article.ArticleType)
		fmt.Fprintf(file, "Status: %s\n", article.Status)
		fmt.Fprintf(file, "Hit Count: %d\n", article.HitCount)

		// Author info
		fmt.Fprintf(file, "\n--- Author ---\n")
		fmt.Fprintf(file, "Username: %s\n", article.Author.UserName)
		fmt.Fprintf(file, "User Slug: %s\n", article.Author.UserSlug)
		if article.Author.RealName != "" {
			fmt.Fprintf(file, "Real Name: %s\n", article.Author.RealName)
		}
		fmt.Fprintf(file, "Avatar: %s\n", article.Author.UserAvatar)
		if article.Author.NameColor != nil {
			fmt.Fprintf(file, "Name Color: %s\n", *article.Author.NameColor)
		}
		fmt.Fprintf(file, "Certification Level: %s\n", article.Author.CertificationLevel)
		if article.Author.ActiveBadge != nil {
			fmt.Fprintf(file, "Active Badge: %s (%s)\n", article.Author.ActiveBadge.DisplayName, article.Author.ActiveBadge.Icon)
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

		// Topic info
		if article.Topic != nil {
			fmt.Fprintf(file, "\n--- Discussion ---\n")
			fmt.Fprintf(file, "Topic ID: %d\n", article.Topic.ID)
			fmt.Fprintf(file, "Top Level Comments: %d\n", article.Topic.TopLevelCommentCount)
		}

		// Summary
		if article.Summary != "" {
			fmt.Fprintf(file, "\n--- Summary ---\n")
			fmt.Fprintf(file, "%s\n", article.Summary)
		}

		fmt.Fprintf(file, "\n")
	}

	return nil
}

// formatTimestamp converts Unix timestamp to readable date in IST
func formatTimestamp(ts int64) string {
	if ts == 0 {
		return "N/A"
	}
	// IST is UTC+5:30
	ist := time.FixedZone("IST", 5*3600+30*60)
	t := time.Unix(ts, 0).In(ist)
	return t.Format("2006-01-02 15:04:05 MST")
}

// formatStringTimestamp converts ISO string timestamp to readable date in IST
func formatStringTimestamp(ts string) string {
	if ts == "" {
		return "N/A"
	}
	// Parse ISO 8601 timestamp
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts // Return original if parsing fails
	}
	// IST is UTC+5:30
	ist := time.FixedZone("IST", 5*3600+30*60)
	return t.In(ist).Format("2006-01-02 15:04:05 MST")
}
