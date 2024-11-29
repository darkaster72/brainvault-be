package hooks

import (
	"brain_vault/queue"
	"brain_vault/utils"
	"encoding/json"
	"log"

	"github.com/pocketbase/pocketbase/core"
)

func OnArticleBeforeCreate(e *core.RecordCreateEvent) error {
	e.Record.Set("content_status", "loading")
	url := e.Record.GetString("url")
	normalizedURL, err := utils.NormalizeUrl(url)
	if err != nil {
		log.Println("Error normalizing URL:", err)
		return err
	}
	e.Record.Set("url", normalizedURL)
	e.Record.Set("title", normalizedURL)
	return nil
}

func OnArticleCreate(e *core.ModelEvent) error {
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
	err = queue.SendMessageWithDefaults(string(jsonData))
	if err != nil {
		log.Println("Error sending message to queue:", err)
	}

	return nil
}
