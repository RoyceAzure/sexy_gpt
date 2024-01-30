package repository

import (
	"context"
	"testing"
	"time"

	"github.com/RoyceAzure/sexy_gpt/logger_service/shared/util"
	"github.com/RoyceAzure/sexy_gpt/logger_service/shared/util/config"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	config := config.Config{
		MongodbAddress: "mongodb://localhost:27017",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongodb, err := ConnectToMongo(ctx, config.MongodbAddress)
	require.NoError(t, err)

	mongoDao := NewMongoDao(mongodb)
	err = mongoDao.Insert(context.Background(), LogEntry{
		ID:          random.RandomString(10),
		ServiceName: random.RandomString(10),
		Message:     random.RandomString(20),
		CreatedAt:   time.Now().UTC(),
	})
	require.NoError(t, err)
}

func TestGetAll(t *testing.T) {
	config := config.Config{
		MongodbAddress: "mongodb://localhost:27017",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongodb, err := ConnectToMongo(ctx, config.MongodbAddress)
	require.NoError(t, err)

	mongoDao := NewMongoDao(mongodb)

	for i := 0; i < 6; i++ {
		TestInsert(t)
	}

	res, err := mongoDao.GetAll(context.Background())
	require.Greater(t, len(res), 5)
	require.NoError(t, err)
}
