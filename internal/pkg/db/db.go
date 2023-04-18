package db

import (
	"context"
	"fmt"
	"time"

	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"

	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var client *ent.Client

func init() {
	var err error
	client, err = openDB()
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	// migrateDB()
}

func GetClient() *ent.Client {
	return client
}

func GetTx(ctx context.Context) *ent.Tx {
	tx, err := client.Tx(ctx)
	if err != nil {
		logger.SugaredLogger.Errorw("Failed to create transaction", "error", err)
		return nil
	}
	return tx
}

func getDbConfig() (driver, dsn string) {
	database := config.GetConfig().Database
	logger.SugaredLogger.Debugw("Get database configuration", "driver", database)
	switch database.Driver {
	case config.DB_DRIVER_SQLITE:
		driver = "sqlite3"
		dsn = fmt.Sprintf("file:%s?_fk=1&_busy_timeout=5000&_synchronous=NORMAL", config.GetConfig().Database.Sqlite.File)
	case config.DB_DRIVER_MYSQL:
		driver = "mysql"
		m := database.MySQL
		fmt.Printf("Connect Mysql %s:%d", m.Host, m.Port)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			m.User, m.Pass, m.Host, m.Port, m.DB)
	default:
		panic("invalid database driver")
	}
	return driver, dsn
}

func openDB() (*ent.Client, error) {
	driver, dsn := getDbConfig()
	drv, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	// Get the underlying sql.DB object of the driver.
	db := drv.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	return ent.NewClient(ent.Driver(drv)), nil
}
