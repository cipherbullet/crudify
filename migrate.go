package main

import (
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func migrate() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Tehran"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
            SingularTable: true, // Use singular table name, same as model name
        },
	})
    if err != nil {
        panic("failed to connect database")
    }

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
        {
            ID: "202409091200",
            Migrate: func(tx *gorm.DB) error {
                return tx.AutoMigrate(&Account{})
            },
            Rollback: func(tx *gorm.DB) error {
                return tx.Migrator().DropTable("account")
            },
        },
    })

    err = m.Migrate()
    if err != nil {
        panic("failed to migrate")
    } else {
        println("Migration did run successfully")
    }
}