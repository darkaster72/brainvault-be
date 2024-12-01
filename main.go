package main

import (
	"brain_vault/hooks"
	"brain_vault/queue"
	"brain_vault/routes"
	"brain_vault/shared"
	"brain_vault/utils"
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	appCtx := shared.AppContext{IsDev: app.IsDev()}
	queue.Initialize(appCtx)

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		utils.SlugMigrator(app)
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
	app.OnRecordBeforeCreateRequest("articles").Add(hooks.OnArticleBeforeCreate(app))
	app.OnModelAfterCreate("articles").Add(hooks.OnArticleCreate)
}
