package database

import (
	"context"
	"fmt"
	"kl-assignment/pkg/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func buildDNS(user, password, host, port, dbName string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName,
	)
}

func ConnectMysql(config config.DB) (*gorm.DB, error) {
	dns := mysql.Open(buildDNS(config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))

	db, err := gorm.Open(dns, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, err
}
