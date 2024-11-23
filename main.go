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
			return c.JSON(http.StatusOK, map[string]string{"message": "Hello " + id})
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
