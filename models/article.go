package models

type Article struct {
	Byline         string `json:"byline"`
	CollectionID   string `json:"collectionId"`
	CollectionName string `json:"collectionName"`
	ContentStatus  string `json:"content_status"`
	Created        string `json:"created"`
	Deleted        bool   `json:"deleted"`
	Excerpt        string `json:"excerpt"`
	ID             string `json:"id"`
	Lang           string `json:"lang"`
	Length         int64  `json:"length"`
	PublishedTime  string `json:"publishedTime"`
	SiteName       string `json:"siteName"`
	Slug           string `json:"slug"`
	Title          string `json:"title"`
	Updated        string `json:"updated"`
	URL            string `json:"url"`
	Userid         string `json:"userid"`
}
