package utils

import (
	"fmt"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func SlugMigrator(app *pocketbase.PocketBase) {
	if !isMigrationApplied(app) {
		articleRecords, err := app.Dao().FindRecordsByExpr("articles", nil)
		if err != nil {
			fmt.Println("Migration failed: ", err)
			return
		}
		collection, err := app.Dao().FindCollectionByNameOrId("contents")
		if err != nil {
			fmt.Println("Migration failed: ", err)
			return
		}

		for _, articleRecord := range articleRecords {
			url := articleRecord.GetString("url")
			content := articleRecord.GetString("content")
			slug, err := NormalizeAndGenerateSlug(url)

			if err != nil {
				fmt.Println("Migration failed: ", err)
				return
			}
			contentRecord := models.NewRecord(collection)
			contentRecord.Set("id", slug)
			contentRecord.Set("content", content)
			err = app.Dao().SaveRecord(contentRecord)

			if err != nil {
				fmt.Println("Migration failed: ", err.Error())
				return
			}
			articleRecord.Set("slug", slug)

			err = app.Dao().SaveRecord(articleRecord)

			if err != nil {
				fmt.Println("Migration failed: ", err)
				return
			}
		}
	}
}

func isMigrationApplied(app *pocketbase.PocketBase) bool {
	var count int64
	res := app.Dao().CollectionQuery().Select("count(*)").From("contents").Build().Row(&count)
	if res != nil {
		return false
	}
	return count > 0
}
