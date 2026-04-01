package models

import (
	"fmt"
	"log"

	"admin-backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 读取配置文件
	cfg, err := utils.LoadConfig("../../conf.yaml")
	if err != nil {
		// 如果相对路径失败，尝试当前目录
		cfg, err = utils.LoadConfig("../conf.yaml")
	}
	if err != nil {
		cfg, err = utils.LoadConfig("conf.yaml")
	}
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	dbConfig := cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName,
		dbConfig.Charset, dbConfig.ParseTime, dbConfig.Loc)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	DB = db
	log.Println("Database initialized successfully")
	return nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}
