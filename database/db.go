// Package database opens and init a database
package database

import (
	"log"

	"github.com/roberChen/ghip/module"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initial(db *gorm.DB) error{
	return db.AutoMigrate(
		&module.IP{},
	)
}

// Open will open a database, filename defined in config
func Open(conf *viper.Viper) *gorm.DB {
	dbfile := conf.GetString("database.file_path")
	log.Printf("using database path `%s`", dbfile)
	db, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		log.Fatalf("opening database failed: %s", err)
	}

	if err := initial(db); err != nil {
		log.Fatalf("init database failed: %s", err)
	}

	return db
}
