package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	Db   *gorm.DB
	conf = Config
)

func init() {
	connectPostgreSQL()
}

func connectPostgreSQL() {
	var err error
	database := conf.Sub("postgres")
	username := database.GetString("username")
	password := database.GetString("password")
	databasename := database.GetString("dbname")
	host := database.GetString("host")
	port := database.GetInt("port")

	dsn := fmt.Sprintf(`host=%s user=%s
		password=%s dbname=%s
		port=%d sslmode=disable 
		TimeZone=Asia/Shanghai`,
		host,
		username,
		password,
		databasename,
		port,
	)

	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
	  fmt.Println("failed to connect database: ", err)
	}
}
