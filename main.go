package main

import (
	"brain_vault/queue"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		// for accepting the parsed content
		e.Router.POST("/api/brainvault/artices/:id/content", func(c echo.Context) error {
			id := c.PathParam("id")

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

			record, _ := app.Dao().FindRecordById("articles", id)
			var body Article

			if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
			}
			record.Set("title", body.Title)
			record.Set("content_status", "loaded")
			record.Set("content", body.Content)
			record.Set("excerpt", body.Excerpt)
			record.Set("length", body.Length)
			record.Set("byline", body.Byline)
			record.Set("siteName", body.SiteName)
			record.Set("lang", body.Lang)
			record.Set("publishedTime", body.PublishedTime)

			app.Dao().Save(record)

			return c.JSON(http.StatusOK, map[string]string{"status": "success"})
		}, apis.RequireAdminAuth())
		return nil
	})

	app.OnModelAfterCreate("articles").Add(func(e *core.ModelEvent) error {
		record, _ := e.Dao.FindRecordById("articles", e.Model.GetId())

		data := map[string]string{
			"id":  e.Model.GetId(),
			"url": record.GetString("url"),
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			return err
		}
		log.Println(string(jsonData))
		queue.SendMessage(string(jsonData))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
