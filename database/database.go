package database

import (
	"fmt"
	"maintainman/config"
	"maintainman/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	dbType := config.AppConfig.GetString("database.driver")
	switch dbType {
	case "mysql":
		DB = initMysql()
	case "sqlite":
		DB = initSqlite()
	default:
		panic(fmt.Errorf("support mysql and sqlite only"))
	}
	initTable()
}

func initSqlite() *gorm.DB {
	dbPath := config.AppConfig.GetString("database.sqlite.path")
	db, err := gorm.Open(sqlite.Open(dbPath))
	if err != nil {
		panic(fmt.Errorf("No error should happen when connecting to database, but got: %+v", err))
	}
	return db
}

func initMysql() *gorm.DB {
	dbHost := config.AppConfig.GetString("database.mysql.host")
	dbPort := config.AppConfig.GetInt("database.mysql.port")
	dbName := config.AppConfig.GetString("database.mysql.name")
	dbParams := config.AppConfig.GetString("database.mysql.params")
	dbUser := config.AppConfig.GetString("database.mysql.user")
	dbPasswd := config.AppConfig.GetString("database.mysql.password")
	dbURL := fmt.Sprintf("%s:%s@(%s:%d)/%s?%s", dbUser, dbPasswd, dbHost, dbPort, dbName, dbParams)

	db, err := gorm.Open(mysql.Open(dbURL))
	if err != nil {
		panic(fmt.Errorf("No error should happen when connecting to database, but got: %+v", err))
	}
	return db
}

func initTable() {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Division{},
		&model.Announce{},
		&model.Order{},
		&model.Status{},
		&model.Tag{},
		&model.Comment{},
		&model.Item{},
		&model.ItemLog{},
	)
	if err != nil {
		panic(fmt.Errorf("Unable to sync the struct to database: %+v", err))
	}
}
