package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"auth-service/config"
	db "auth-service/db/gen"
	"auth-service/utils"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *db.Queries
var random *utils.UtilRandom

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	pgxconfig, err := pgxpool.ParseConfig(config.DataSourceForPGX())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	// pgxconfig.ConnConfig.
	pgxconfig.ConnConfig.Tracer = otelpgx.NewTracer()
	testDB, err := pgxpool.NewWithConfig(context.Background(), pgxconfig)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = db.New(testDB)
	random = utils.NewUtilRandom()

	os.Exit(m.Run())
}
