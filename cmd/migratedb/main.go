package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"notionboy/internal/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")

	// ask user to confirm
	logger.SugaredLogger.Infow("Ready to migrate DB", "db_host", db_host, "db_port", db_port, "db_name", db_name)
	fmt.Printf("Do you want to continue? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logger.SugaredLogger.Fatalw("read input failed", "err", err)
	}
	input = strings.TrimSpace(input)
	if input != "y" {
		logger.SugaredLogger.Fatalw("user cancel")
	}

	// start migrate
	logger.SugaredLogger.Infow("Start migrate DB")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", db_user, db_password, db_host, db_port, db_name))
	if err != nil {
		logger.SugaredLogger.Fatalw("open mysql failed", "err", err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		logger.SugaredLogger.Fatalw("mysql with instance failed", "err", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		logger.SugaredLogger.Fatalw("get current working directory failed", "err", err)
	}
	migratePath := fmt.Sprintf("file://%s/db/migrations", cwd)
	logger.SugaredLogger.Infow("migrate path", "migratePath", migratePath)
	m, err := migrate.NewWithDatabaseInstance(
		migratePath,
		"mysql",
		driver,
	)
	if err != nil {
		logger.SugaredLogger.Fatalw("migrate new with database instance failed", "err", err)
	}

	err = m.Up()
	if err != nil {
		logger.SugaredLogger.Fatalw("migrate up failed", "err", err)
	}
}
