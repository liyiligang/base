/*
 * Copyright 2021 liyiligang.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package Jorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"time"
	"unicode"
)

type OrmInitConfig struct {
	SqlDsn      string
	MaxKeepConn int
	MaxConn     int
	MaxLifetime time.Duration
	LogWrite 	io.Writer
	TableCheck  func(*gorm.DB)
}

type ormNamer struct {
	schema.NamingStrategy
}

func GormInit(config OrmInitConfig) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}
	if config.LogWrite != nil {
		newLogger := logger.New(
			log.New(config.LogWrite, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // 慢 SQL 阈值
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // 禁用彩色打印
			},
		)
		gormConfig.Logger = newLogger
	}
	gormConfig.NamingStrategy = ormNamer{}
	db, err := gorm.Open(mysql.Open(config.SqlDsn), gormConfig)
	if err != nil {
		return db, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	sqlDB.SetConnMaxLifetime(config.MaxLifetime) //每次操作数据库允许的最大时间限制
	sqlDB.SetMaxIdleConns(config.MaxKeepConn)
	sqlDB.SetMaxOpenConns(config.MaxConn)
	if config.TableCheck != nil {
		config.TableCheck(db)
	}
	return db, err
}

func (namer ormNamer) getLowerStr(str string) string{
	r := []rune(str)
	if len(r) != 0 {
		r[0] = unicode.ToLower(r[0])
	}
	return string(r)
}

func (namer ormNamer) TableName(table string) string {
	return namer.getLowerStr(table)
}

func (namer ormNamer) ColumnName(table, column string) string {
	return namer.getLowerStr(column)
}

func (namer ormNamer) IndexName(table, column string) string {
	return namer.getLowerStr(column)
}