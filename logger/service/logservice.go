package logservice

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	repository "github.com/RoyceAzure/sexy_gpt/logger_service/repository/mongodb"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

var Logger zerolog.Logger

type MongoLogger struct {
	mongoDao repository.IMongoDao
}

func SetUpMutiMongoLogger(mongoLogger *MongoLogger, serviceId string) error {
	if mongoLogger == nil {
		return fmt.Errorf("mongo logger is not init")
	}
	multiLogger := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, mongoLogger)
	Logger = zerolog.New(multiLogger).With().Str("service_name", serviceId).Timestamp().Logger()
	return nil
}

func NewMongoLogger(mongoDao repository.IMongoDao) *MongoLogger {
	return &MongoLogger{
		mongoDao: mongoDao,
	}
}

func (mw *MongoLogger) Write(p []byte) (n int, err error) {
	// Insert the record into the collection.

	if mw == nil {
		return 0, fmt.Errorf("mongo logger is not init")
	}

	var logEntry bson.M
	err = json.Unmarshal(p, &logEntry)
	if err != nil {
		return 0, err
	}
	err = mw.mongoDao.InsertBsonM(context.Background(), logEntry)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
