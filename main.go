package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"auth-service/config"
	"auth-service/endpoint"
	"auth-service/routes"
	"auth-service/server"

	"auth-service/language"

	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/rubenv/sql-migrate"

	"auth-service/db/repo"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)

	migrationDB(config)

	conn, err := pgxpool.New(context.Background(), config.DataSourceForPGX())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	defer conn.Close()

	repos := repo.NewRepo(config, conn)
	multiLocalizer := language.LoadAllFileLanguage("./language/localizations_src/")

	server, err := server.NewServer(config, repos, &multiLocalizer)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	ep := endpoint.NewEndpointModule(config, repos, &multiLocalizer, server.TokenMaker)
	httpRouter := routes.NewHttpRouter(server, ep)

	httpAddress := fmt.Sprintf(`:%d`, config.APIInfo.HTTPPort)
	log.Info().Msgf("HTTP address: %s", httpAddress)
	err = http.ListenAndServe(httpAddress, httpRouter)
	log.Fatal().Err(err)
}

func migrationDB(config *config.Config) {
	dataSource := config.DataSourceForSQL()
	db, err := sql.Open(config.Database.Driver, dataSource)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer func() {
		log.Info().Msg("Migration: Close db")
		db.Close()
	}()

	log.Info().Msg(`Migration database is starting...`)
	migrations := migrate.FileMigrationSource{Dir: "./db/migration"}
	migrate.SetTable("migrations")
	n, err := migrate.Exec(db, config.Database.Driver, migrations, migrate.Up)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("Applied %d migrations!", n)
}
