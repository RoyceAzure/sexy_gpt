package db

import (
	"context"
	"os"
	"testing"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var testDao Dao
var testConn *pgxpool.Pool

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	config, err := config.LoadConfig("../../../")
	if err != nil {
		log.Fatal().Err(err).Msg("err load config")
	}
	ctx := context.Background()
	testConn, err = pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("err create db connect")
	}
	testDao = NewSQLDao(testConn)
}
