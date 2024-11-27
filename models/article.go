package models

type Article struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	TextContent   string `json:"textContent"`
	Length        int    `json:"length"`
	Excerpt       string `json:"excerpt"`
	Byline        string `json:"byline"`
	Dir           string `json:"dir"`
	SiteName      string `json:"siteName"`
	Lang          string `json:"lang"`
	PublishedTime string `json:"publishedTime"`
}
