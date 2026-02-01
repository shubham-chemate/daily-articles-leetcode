package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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
						topicId
						title
						slug
						summary
						author {
							userName
						}
						createdAt
						updatedAt
						articleType
						tags {
							name
							slug
							tagType
						}
						reactions {
							count
							reactionType
						}
					}
				}
			}
		}
	`
)

// fetchArticlesAfterTime fetches all articles published after the given cutoff time using pagination
func fetchArticlesAfterTime(cutoffTime time.Time) ([]Article, error) {
	var allArticles []Article
	batchSize := 100
	skip := 0

	for {
		fmt.Printf("Fetching batch starting at offset %d...\n", skip)

		// Fetch batch
		batch, err := fetchDiscussArticlesWithSkip(batchSize, skip)
		if err != nil {
			return nil, err
		}

		if len(batch) == 0 {
			break // No more articles
		}

		foundOlderArticle := false
		for _, article := range batch {
			articleTime, err := time.Parse(time.RFC3339, article.CreatedAt)
			if err != nil {
				continue // Skip if we can't parse the time
			}

			if articleTime.After(cutoffTime) {
				allArticles = append(allArticles, article)
			} else {
				foundOlderArticle = true
				break
			}
		}

		if foundOlderArticle || len(batch) < batchSize {
			break // Stop if we found older articles or reached the end
		}

		skip += batchSize
	}

	return allArticles, nil
}

// fetchDiscussArticlesWithSkip fetches articles with pagination support
func fetchDiscussArticlesWithSkip(count int, skip int) ([]Article, error) {
	reqBody := map[string]interface{}{
		"query": discussTopicsQuery,
		"variables": map[string]interface{}{
			"orderBy":  "MOST_RECENT",
			"keywords": []string{},
			"tagSlugs": []string{},
			"skip":     skip,
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

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result ArticlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var articles []Article
	for _, edge := range result.Data.UgcArticleDiscussionArticles.Edges {
		articles = append(articles, edge.Node)
	}

	// Articles are already sorted by NEWEST, no need to sort again
	return articles, nil
}
