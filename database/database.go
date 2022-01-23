package database

import (
	"fmt"
	"os"

	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	DB *gorm.DB

	dbHost     string
	dbUser     string
	dbName     string
	dbPort     string
	dbSslMode  string
	dbPassword string
)

func Connect() {
	LoadEnv()
	dbUri := fmt.Sprintf("host=%s user=%s database=%s port=%s password=%s sslmode=%s",
		dbHost, dbUser, dbName, dbPort, dbPassword, dbSslMode)
	var err error
	DB, err = gorm.Open(gormPostgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	doMigrate()
}

func LoadEnv() {
	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbName = os.Getenv("DB_NAME")
	dbPort = os.Getenv("DB_PORT")
	dbSslMode = os.Getenv("DB_SSLMODE")
	dbPassword = os.Getenv("DB_PASSWORD")
}

func doMigrate() {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://_migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		handleMigrationError(m, err)
	}
}

func handleMigrationError(m *migrate.Migrate, err error) {
	version, dirty, error := m.Version()
	if error != nil {
		panic(error)
	}
	error = migrate.ErrDirty{Version: int(version)}
	if err == error && dirty {
		m.Force(int(version - 1))
		doMigrate()
	} else {
		panic(err)
	}
}
