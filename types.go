package main

// Article represents an article from LeetCode discuss section
type Article struct {
	UUID        string     `json:"uuid"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Summary     string     `json:"summary"`
	Author      Author     `json:"author"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
	ArticleType string     `json:"articleType"`
	Tags        []Tag      `json:"tags"`
	Reactions   []Reaction `json:"reactions"`
}

// Author represents the article author (only userName needed)
type Author struct {
	UserName string `json:"userName"`
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
