package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:13306)/connectify?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to initialize database:", err)
		panic(err)
	}

	err = db.AutoMigrate(&UserModel{})
	if err != nil {
		fmt.Println("Failed to migrate database:", err)
		panic(err)
	}
	return db
}
