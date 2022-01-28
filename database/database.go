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
	M  *migrate.Migrate

	dbHost     string
	dbUser     string
	dbName     string
	dbPort     string
	dbSslMode  string
	dbPassword string

	mPath string
)

func Connect() error {
	LoadEnv()
	var err error
	fmt.Println(DbConnectionString())
	DB, err = gorm.Open(gormPostgres.Open(DbConnectionString()), &gorm.Config{})
	if err != nil {
		return err
	}
	setupMigration()
	doMigrate()
	return nil
}

func LoadEnv() {
	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbName = os.Getenv("DB_NAME")
	dbPort = os.Getenv("DB_PORT")
	dbSslMode = os.Getenv("DB_SSLMODE")
	dbPassword = os.Getenv("DB_PASSWORD")
	mPath = os.Getenv("M_PATH")
}

func DbConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSslMode)
}

func setupMigration() {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	M, err = migrate.NewWithDatabaseInstance(mPath, "postgres", driver)
	if err != nil {
		panic(err)
	}
}

func doMigrate() {
	err := M.Up()
	if err != nil && err != migrate.ErrNoChange {
		handleMigrationError(M, err)
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
