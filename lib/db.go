package lib

import (
	"log"

	"github.com/xeonds/phi-plug-go/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB connects to the configured database and runs the migrator.
func NewDB(cfg *config.DatabaseConfig, migrate func(*gorm.DB) error) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	switch cfg.Type {
	case "mysql":
		dsn := cfg.User + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		db, err = gorm.Open(sqlite.Open(cfg.DB), &gorm.Config{})
	}
	if err != nil {
		log.Fatal("connect database: ", err)
	}
	if cfg.Migrate {
		if migrate == nil {
			log.Fatal("migrator is nil")
		}
		if err := migrate(db); err != nil {
			log.Fatal("migrate: ", err)
		}
	}
	return db
}
