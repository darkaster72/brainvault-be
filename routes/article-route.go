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

		// Find the article record by ID
		record, err := app.Dao().FindRecordById("articles", id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Article not found"})
		}

		var body models.Article

		// Decode the request body into the Article struct
		if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		// Update the record fields with the new data
		record.Set("title", body.Title)
		record.Set("content_status", "loaded")
		record.Set("content", body.Content)
		record.Set("excerpt", body.Excerpt)
		record.Set("length", body.Length)
		record.Set("byline", body.Byline)
		record.Set("siteName", body.SiteName)
		record.Set("lang", body.Lang)
		record.Set("publishedTime", body.PublishedTime)

		// Save the updated record to the database
		if err := app.Dao().Save(record); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save article"})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "success"})
	}
}
