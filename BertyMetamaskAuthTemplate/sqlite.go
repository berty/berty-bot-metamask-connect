package main

import (
	sqlite "github.com/flyingtime/gorm-sqlcipher"
	"gorm.io/gorm"
)

type User struct {
	BertyConvPubKey string `gorm:"primary_key"`
	MetamaskAddress string
	Nonce           uint32
	Verified        uint8
}

func NewSqliteDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.db"))
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
