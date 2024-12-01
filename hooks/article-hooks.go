package hooks

import (
	"brain_vault/queue"
	"brain_vault/utils"
	"encoding/json"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func OnArticleBeforeCreate(app *pocketbase.PocketBase) func(e *core.RecordCreateEvent) error {
	return func(e *core.RecordCreateEvent) error {
		e.Record.Set("content_status", "loading")
		url := e.Record.GetString("url")
		normalizedURL, err := utils.NormalizeUrl(url)
		collection, _ := app.Dao().FindCollectionByNameOrId("contents")

		if err != nil {
			log.Println("Error normalizing URL:", err)
			return err
		}
		slug := utils.GenerateSlug(normalizedURL)
		content := models.NewRecord(collection)
		content.SetId(slug)

		if err := app.Dao().SaveRecord(content); err != nil {
			log.Println("Failed to create Content:", err)
			return err
		}
		e.Record.Set("url", normalizedURL)
		e.Record.Set("title", normalizedURL)
		e.Record.Set("slug", slug)
		return nil
	}
}

func OnArticleCreate(e *core.ModelEvent) error {
	record, _ := e.Dao.FindRecordById("articles", e.Model.GetId())

	data := map[string]string{
		"id":   e.Model.GetId(),
		"url":  record.GetString("url"),
		"slug": record.GetString("slug"),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return err
	}
	log.Println(string(jsonData))
	err = queue.SendMessageWithDefaults(string(jsonData))
	if err != nil {
		log.Println("Error sending message to queue:", err)
	}

	return nil
}
