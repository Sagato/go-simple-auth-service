package db

import (
	"authentication-service/model"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var db *pg.DB

func NewDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		User: "testuser",
		Password: "root",
		Database: "authentication",
	})
	return db
}

func NewTestDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		User: "testuser",
		Password: "root",
		Database: "authenticationTest",
	})
	return db
}

func CreateSchema(db *pg.DB) error {

	for _, model := range []interface{}{
		(*model.Activation)(nil),
		(*model.GrantTypeResponse)(nil),
		(*model.User)(nil),
	} {
 		if err := db.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:true,
		}); err != nil {
			return err
		}

		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
			IfNotExists: true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateTestSchemas(db *pg.DB) error {

	for _, model := range []interface{}{
		(*model.User)(nil),
		(*model.Activation)(nil),
		(*model.GrantTypeResponse)(nil),
	} {
		if err := db.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
		}); err != nil {
			return err
		}

		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}