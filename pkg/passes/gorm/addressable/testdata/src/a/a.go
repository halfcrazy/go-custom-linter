package a

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID uint `gorm:"primaryKey,autoIncrement"`
}

func do() {
	db, err := gorm.Open(sqlite.Open("./gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(User{})
	var rv []User
	if err := db.Model(&User{}).Find(&rv).Error; err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Println(rv)
	var t = User{}
	if err := db.Model(User{}).Create(&t).Error; err != nil {
		fmt.Printf("%+v", err)
	}
	if err := db.Model(User{}).Find(&rv).Error; err != nil {
		fmt.Printf("%+v", err)
	}
	if err := db.Model(User{}).Find(rv).Error; err != nil { // want "not addressable param passed to gorm Find"
		fmt.Printf("%+v", err)
	}
	if err := db.Model(User{}).Delete(t).Error; err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Println(rv)
}
