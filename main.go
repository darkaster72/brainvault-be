package main

import (
	"brain_vault/hooks"
	"brain_vault/queue"
	"brain_vault/routes"
	"brain_vault/shared"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	appCtx := shared.AppContext{IsDev: app.IsDev()}
	queue.Initialize(appCtx)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		registerRoutes(app, e)
		return nil
	})
	registerHooks(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(app *pocketbase.PocketBase, e *core.ServeEvent) {
	e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
	e.Router.POST("/api/brainvault/artices/:id/content", routes.SaveArticleContent(app), apis.RequireAdminAuth())
}

func registerHooks(app *pocketbase.PocketBase) {
	app.OnModelAfterCreate("articles").Add(hooks.OnArticleCreate)
}
