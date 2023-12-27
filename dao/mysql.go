package dao

import (
	"douyin/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 内涵连接池管理
var db *gorm.DB

// DbInit  数据库初始函数
func DbInit() error {
	var err error
	dsn := config.C.Mysql.Username + ":" + config.C.Mysql.Password + "@tcp(" + config.C.Mysql.Ipaddress + ":" + config.C.Mysql.Port + ")/" + config.C.Mysql.Dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println(dsn)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Comment{}, &Follow{}, &Like{}, &Message{}, &User{}, &Video{}, &User2{}, &UserCount{})
	return err
}

func ExecuteTransaction(operation func(*gorm.DB) error) error {
	tx := db.Begin()
	err := operation(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("事务执行失败：%w", err)
	}
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("事务提交失败：%w", err)
	}
	return nil
}
func DBInit() {
	dsn := "root:@tcp(localhost:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
