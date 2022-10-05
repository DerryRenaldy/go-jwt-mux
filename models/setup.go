package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectionDatabase() {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/go_jwt_mux"))
	if err != nil {
		panic(err)
	}

	errMigrate := db.AutoMigrate(&User{})
	if errMigrate != nil {
		log.Printf("error auto migrate: %v \n", errMigrate)
	}

	DB = db
}
