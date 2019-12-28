package goKLC

import "fmt"

const MYSQL DBType = "mysql"
const POSTGRES DBType = "postgres"
const SQLITE3 DBType = "sqlite3"
const MSSQL DBType = "mssql"
const NONE DBType = "none"

type DBType string

func connectDB(dbType DBType) string {
	var dbUrl string
	dbUser := _config.Get("DBUser", "root")
	dbPassword := _config.Get("DBPassword", "")
	dbName := _config.Get("DBName", "go_klc")
	dbHost := _config.Get("DBHost", "127.0.0.1")
	dbPort := _config.Get("DBPort", "3306")

	switch dbType {
	case MYSQL:
		dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
		break
	case POSTGRES:
		dbUrl = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
		break
	case SQLITE3:
		dbUrl = fmt.Sprintf("%s", dbHost)
		break
	case MSSQL:
		dbUrl = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", dbUser, dbPassword, dbHost, dbPort, dbName)
		break
	}

	return dbUrl
}
