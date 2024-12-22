package models

import (
	"log"
	"os"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var bloomFilter *bloom.BloomFilter

func DBInit() {
	// Define logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Microsecond, // Slow SQL threshold
			LogLevel:                  logger.Info,      // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound
			Colorful:                  true,             // Disable color
		},
	)

	// Initialize SQLite database
	var err error
	db, err = gorm.Open(sqlite.Open("access.db"),
		&gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate the UserAccess schema
	db.AutoMigrate(&UserAccess{})
}

func BloomFilterSetup() {
	// Initialize Bloom filter with estimates
	bloomFilter = bloom.NewWithEstimates(10000, 0.01)
}
