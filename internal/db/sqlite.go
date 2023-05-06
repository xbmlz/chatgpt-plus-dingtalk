package db

import (
	"os"

	"github.com/xbmlz/chatgpt-dingtalk/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		logger.Fatal(err)
	}
	db, err := gorm.Open(sqlite.Open("data/chat.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Fatal("failed to open sqlite3: %v", err)
	}
	dbObj, err := db.DB()
	if err != nil {
		logger.Fatal("failed to get sqlite3 obj: %v", err)
	}
	// See https://github.com/glebarez/sqlite/issues/52
	dbObj.SetMaxOpenConns(1)
	_ = db.AutoMigrate(
		Chat{},
	)
	DB = db
}
