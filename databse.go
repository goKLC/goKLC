package goKLC

import (
	"fmt"
	"github.com/goKLC/goKLC/SqlProviders"
	"github.com/jinzhu/gorm"
)

const MYSQL DBType = "mysql"
const POSTGRES DBType = "postgres"
const SQLITE3 DBType = "sqlite3"
const MSSQL DBType = "mssql"
const NONE DBType = "none"

type DBType string

func connectDB(dbType DBType) {
	var err error
	var dbUrl string
	dbUser := _config.Get("DBUser", "root")
	dbPassword := _config.Get("DBPassword", "")
	dbName := _config.Get("DBName", "go_klc")
	dbHost := _config.Get("DBHost", "127.0.0.1")
	dbPort := _config.Get("DBPort", "3306")

	switch dbType {
	case MYSQL:
		SqlProviders.MysqlInit()
		dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
		break
	case POSTGRES:
		SqlProviders.PostgresInit()
		dbUrl = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
		break
	case SQLITE3:
		SqlProviders.SqliteInit()
		dbUrl = fmt.Sprintf("%s", dbHost)
		break
	case MSSQL:
		SqlProviders.MssqlInit()
		dbUrl = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", dbUser, dbPassword, dbHost, dbPort, dbName)
		break
	}

	if dbType == NONE {
		return
	}

	_DB, err = gorm.Open(string(dbType), dbUrl)
	_ = _DB

	if err != nil {
		_app.Log().Error(err.Error(), nil)
	}
}
