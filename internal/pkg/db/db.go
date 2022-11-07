package db

import (
	"context"
	"fmt"
	"log"
	"notionboy/db/ent"
	"notionboy/internal/pkg/config"
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
	database := config.GetConfig().Databases
	if database.MySQL.Host != "" {
		driver = "mysql"
		m := database.MySQL
		fmt.Printf("Connect Mysql %s:%d", m.Host, m.Port)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			m.User, m.Pass, m.Host, m.Port, m.Database)
	} else {
		driver = "sqlite3"
		dsn = fmt.Sprintf("file:%s?_fk=1&_busy_timeout=5000&_synchronous=NORMAL", config.GetConfig().Databases.Sqlite.File)
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
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}