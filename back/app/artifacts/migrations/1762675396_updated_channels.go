package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3009067695")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("text1400097126")

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_961350965",
			"hidden": false,
			"id": "relation1400097126",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "country",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_3304764897",
			"hidden": false,
			"id": "relation3571151285",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "language",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3009067695")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text1400097126",
			"max": 0,
			"min": 0,
			"name": "country",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation1400097126")

		// remove field
		collection.Fields.RemoveById("relation3571151285")

		return app.Save(collection)
	})
}
