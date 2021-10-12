package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"brucosijada.kset.org/app/models"
)

type databaseProvider struct {
	conn *gorm.DB
}

var db_ = databaseProvider{}

func DatabaseProvider() databaseProvider {
	return db_
}

func (p databaseProvider) Client() *gorm.DB {
	return db_.conn
}

func (p databaseProvider) log(format string, v ...interface{}) {
	log.Printf("|DB> "+format, v...)
}

func (p databaseProvider) Register() (err error) {
	db_.log("Connecting to PostgreSQL...")

	logger_ := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags|log.LUTC),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db_.conn, err = gorm.Open(
		postgres.Open(os.Getenv("DATABASE_URL")),
		&gorm.Config{
			Logger: logger_,
		},
	)

	db_.log("Connected to PostgreSQL")
	if err != nil {
		return err
	}

	db_.log("Running the migrations...")
	err = db_.conn.AutoMigrate(
		&models.Image{},
		&models.ImageVariation{},
		&models.Sponsor{},
		&models.User{},
		&models.Artist{},
		&models.UserInvitation{},
		&models.Page{},
	)
	if err != nil {
		return err
	}
	db_.log("Done with migrations")

	return nil
}
