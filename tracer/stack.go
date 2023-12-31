package tracer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

var (
	mux sync.Mutex
)

// Stack represents a stack of messages.
type Stack struct {
	loggerType string         // "info", "debug", "error", "warning", "disabled"
	logger     *logrus.Logger // Logger for logging
	filePath   string         // File path for storing messages
}

func (s *Stack) Logger() *logrus.Logger {
	return s.logger
}

func NewStack(loggerType string, logger *logrus.Logger) (*Stack, error) {

	ms := &Stack{
		loggerType: loggerType,
		logger:     logger,
	}

	switch loggerType {
	case info:
		ms.logger.SetLevel(logrus.InfoLevel)
	case debug:
		ms.logger.SetLevel(logrus.DebugLevel)
	case errorType:
		ms.logger.SetLevel(logrus.ErrorLevel)
	case warning:
		ms.logger.SetLevel(logrus.WarnLevel)
	default:
		ms.logger.SetLevel(logrus.InfoLevel)
	}

	return ms, nil
}

// InitDatabase
// eg db, err := InitDatabase("your_db_path.db", &Message{}, &AnotherModel{})
func InitDatabase(dbPath string, models ...interface{}) (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		// Enable foreign keys by setting the foreign_key_constraint pragma to 1.
		PrepareStmt: true, // Optional, improves performance by using prepared statements.
	})
	if err != nil {
		fmt.Println("INIT SQLLite: ", err)
	}
	db.Exec("PRAGMA foreign_keys = ON;") // To make batch migration works for inheritance (or when altering (drop & transform in real) table shows error)

	// Auto Migrate the specified models
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return nil, fmt.Errorf("failed to auto-migrate model: %v", err)
		}
	}

	return db, nil
}
