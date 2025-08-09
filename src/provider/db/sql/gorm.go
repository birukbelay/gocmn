package sql_db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	// "gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	conf "github.com/birukbelay/gocmn/src/config"
	cmn "github.com/birukbelay/gocmn/src/logger"
)

func NewSqlDb(config *conf.SqlDbConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var confg gorm.Config
	var dsn string
	if os.Getenv("ENABLE_GORM_LOGGER") != "" {
		confg = gorm.Config{}
	} else {
		confg = gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	switch strings.ToLower(config.Driver) {
	// case "mysql":
	// 	dsn := config.Username + ":" + config.Password + "@tcp(" + config.MongoHost + ":" + strconv.Itoa(config.MongoPort) + ")/" + config.DbName + "?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=UTC"
	// 	db, err = gorm.Open(mysql.Open(dsn), &confg)
	// 	brea

	// case "sqlite", "sqllite":
	// 	// For SQLite, the DSN is typically the file path
	// 	dsn = config.DbName // e.g., "test.db"
	// 	db, err = gorm.Open(sqlite.Open(dsn), &confg)
	// case "sqlserver", "mssql":
	// 	dsn := "sqlserver://" + config.Username + ":" + config.Password + "@" + config.MongoHost + ":" + strconv.Itoa(config.MongoPort) + "?database=" + config.DbName
	// 	db, err = gorm.Open(sqlserver.Open(dsn), &confg)
	// 	break
	case "postgresql", "postgres":
		dsn = fmt.Sprintf("user=%s dbname=%s password=%s port=%s host=%s sslmode=%s", config.Username, config.DbName, config.Password, config.SqlPort, config.SqlHost, config.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), &confg)
	}
	if err != nil || db == nil {
		cmn.LogTrace("failed to connect to database:", err.Error())
		cmn.LogTrace("conn string is :", dsn)
		log.Fatal(err)
		//panic("Failed to connect database")
	}

	cmn.LogTrace("Connection Opened to Postgress Database", config.DbName)

	return db, err
}
