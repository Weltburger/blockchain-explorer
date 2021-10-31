package storage

import (
	"database/sql"
	"explorer/models"
	"fmt"
	_ "github.com/mailru/go-clickhouse"
	"gorm.io/driver/clickhouse"

	//"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"log"
	"sync"
)

type Database struct {
	GormDB *gorm.DB
	sync.Once
}

func GetDB() *Database {
	gormDB := startDB()

	return &Database{GormDB: gormDB}
}

func startDB() *gorm.DB {
	dbUrl := "http://localhost:8123/blocks?debug=true"
	sqlDB, err := sql.Open("clickhouse", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("aaa")

	gormDB, err := gorm.Open(clickhouse.New(clickhouse.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return gormDB
}

func (s *Database) Migrate() {
	s.Do(func() {
		err := s.GormDB.Set("gorm:table_options", "engine = MergeTree()\n PARTITION BY toYYYYMMDD(Timestamp)\n ORDER BY (Hash)").AutoMigrate(&models.Block{})
		if err != nil {
			return
		}
	})
}

