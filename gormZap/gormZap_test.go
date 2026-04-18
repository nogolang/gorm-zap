package gormZap

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Test_gromZap(t *testing.T) {
	// 创建zap对象
	zapLogger, _ := zap.NewProduction()
	// 创建gormZap对象
	gormZap := NewGormZap(zapLogger, gormlogger.Info, time.Second)
	// 创建gorm对象
	url := "host=192.168.80.128 user=root password=root dbname=shop_product port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ := gorm.Open(postgres.Open(url))
	db.Logger = gormZap

	type Sku struct {
		Id int `json:"id"`
	}
	tx := db.Table("sku").Where("id = ?", 1).Find(&Sku{})
	if tx.Error != nil {
		t.Error(tx.Error)
	}
}
