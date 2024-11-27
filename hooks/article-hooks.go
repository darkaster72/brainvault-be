package hooks

import (
	"brain_vault/queue"
	"encoding/json"
	"log"

	"github.com/pocketbase/pocketbase/core"
)

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
	queue.SendMessage(string(jsonData))
	return nil
}
