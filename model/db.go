package model

import (
	"fmt"

	"../config"
	"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)


var (
	DBConn *gorm.DB
)

func GormInit() error {
	conf := &config.DBConfig{}
	err := conf.Read()

	DBConn, err = gorm.Open("postgres",
		fmt.Sprintf("host=localhost user=%s dbname=%s sslmode=disable password=%s",
			conf.DBUser, conf.DBName, conf.DBPass))
	if err != nil {
		return err
	}

  DBConn.AutoMigrate(&User{})
	return nil
}

func GormClose() error {
	if DBConn != nil {
		return DBConn.Close()
	}
	return nil
}
