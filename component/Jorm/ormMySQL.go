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
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"time"
)

type OrmConfig struct {
	Name        string
	SqlDsn      string
	MaxKeepConn int
	MaxConn     int
	MaxLifetime time.Duration
	ShowLog     bool
	LogWrite    io.Writer
	schemaNamer schema.Namer
	TableCheck  func(*gorm.DB) error
}

func GormInit(config OrmConfig) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}
	logLevel := logger.Warn
	if config.ShowLog {
		logLevel = logger.Info
	}
	if config.LogWrite != nil {
		newLogger := logger.New(
			log.New(config.LogWrite, "", 0), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logLevel,    // Log level
				IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound(记录未找到)错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
		gormConfig.Logger = newLogger //设置日志输出
	}

	gormConfig.NamingStrategy = config.schemaNamer //设置自动生成表名, 字段名等规则

	var dialector gorm.Dialector
	switch config.Name {
	case "mysql":
		dialector = mysql.Open(config.SqlDsn)
		break
	case "postgresql":
		dialector = postgres.Open(config.SqlDsn)
		break
	case "sqlserver":
		dialector = sqlserver.Open(config.SqlDsn)
		break
	default:
		dialector = sqlite.Open(config.SqlDsn)
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(config.MaxLifetime) //每次操作数据库允许的最大时间限制
	sqlDB.SetMaxIdleConns(config.MaxKeepConn)
	sqlDB.SetMaxOpenConns(config.MaxConn)
	if config.TableCheck != nil {
		err := config.TableCheck(db)
		if err != nil {
			return nil, err
		}
	}
	return db, err
}

//type ormNamer struct {
//	schema.NamingStrategy
//}
//
//func (namer ormNamer) getLowerStr(str string) string {
//	r := []rune(str)
//	if len(r) != 0 {
//		r[0] = unicode.ToLower(r[0])
//	}
//	return string(r)
//}

//func (namer ormNamer) TableName(table string) string {
//	return namer.getLowerStr(table)
//}
//
//func (namer ormNamer) ColumnName(table, column string) string {
//	return namer.getLowerStr(column)
//}
//
//func (namer ormNamer) IndexName(table, column string) string {
//	return namer.getLowerStr(column)
//}
