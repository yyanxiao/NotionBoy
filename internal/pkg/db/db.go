package db

import (
	"context"
	"fmt"
	"notionboy/db/ent"
	"notionboy/db/ent/migrate"
	"notionboy/internal/pkg/config"
	"notionboy/internal/pkg/logger"
	"time"

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
	migrateDB()
}

func GetClient() *ent.Client {
	return client
}

func getDbConfig() (driver string, dsn string) {
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
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
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

func migrateDB() {
	ctx := context.Background()
	if err := client.Debug().Schema.Create(
		ctx,
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	); err != nil {
		logger.SugaredLogger.Fatalw("Failed creating schema resources", "err", err)
	}
}
