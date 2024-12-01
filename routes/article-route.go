package routes

import (
	"brain_vault/models"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
)

// SaveArticleContent handles saving article content to the database
func SaveArticleContent(app *pocketbase.PocketBase) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")

		article, err := app.Dao().FindRecordById("articles", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Article not found"})
		}
		slug := article.GetString("slug")
		content, _ := app.Dao().FindRecordById("contents", slug)

		var body models.Article

		// Decode the request body into the Article struct
		if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		// Update the record fields with the new data
		article.Set("title", body.Title)
		article.Set("content_status", "loaded")
		article.Set("excerpt", body.Excerpt)
		article.Set("length", body.Length)
		article.Set("byline", body.Byline)
		article.Set("siteName", body.SiteName)
		article.Set("lang", body.Lang)
		article.Set("publishedTime", body.PublishedTime)
		content.Set("content", body.Content)
		content.Set("text_content", body.TextContent)

		if err := app.Dao().Save(content); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save content"})
		}

		if err := app.Dao().Save(article); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save article"})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}
