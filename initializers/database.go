package initializers

import (
	"github.com/Gonnekone/onlineStore/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectToDBDocker() {
	var err error
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("postgres_port") +
		" sslmode=" + os.Getenv("sslmode") +
		" TimeZone=" + os.Getenv("TimeZone")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + dsn)
	}
	models.AutoMigrate(DB)
}

func ConnectToDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	var err error

	dsn := "host=localhost user=postgres password=Fakon4565 dbname=onlineStore port=5432 sslmode=disable TimeZone=Europe/Moscow"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database: " + dsn)
	}

	models.AutoMigrate(DB)
}
