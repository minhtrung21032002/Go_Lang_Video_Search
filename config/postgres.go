package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	user     = "postgres"
	password = "trung123"
	port     = "5432"
)

func Protgres() (*gorm.DB, error) {
	dsn := "host=127.0.0.1 user=" + user + " password=" + password + " dbname=postgres port=" + port + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
