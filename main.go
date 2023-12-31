package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/config"
	"auth-service/endpoint"
	"auth-service/proto/service"
	"auth-service/routes"
	"auth-service/server"
	"auth-service/transport"

	"auth-service/language"

	"github.com/jackc/pgx/v5/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	"google.golang.org/grpc"

	"auth-service/db/repo"

	kitgrpc "github.com/go-kit/kit/transport/grpc"

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

	ep := endpoint.NewEndpointModule(config, repos, server.TokenMaker)
	trans := transport.NewModule(config, repos, &multiLocalizer, server.TokenMaker, ep)

	httpRouter := routes.NewHttpRouter(server, trans, ep)
	grpcRouter := routes.NewGRPCRouter(server, trans)

	var (
		httpAddress = fmt.Sprintf(`:%d`, config.APIInfo.HTTPPort)
		grpcAddress = fmt.Sprintf(`:%d`, config.APIInfo.GRPCPort)
	)

	// Implement graceful shutdown
	errs := make(chan error, 2)

	// HTTP server
	go func() {
		log.Info().Msgf(`HTTP server (%s) started successfully 🚀`, httpAddress)
		errs <- http.ListenAndServe(httpAddress, httpRouter)
	}()

	// gRPC server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	go func() {
		grpcListener, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			log.Err(err).Msg("Start gRPC service is error")
			errs <- err
			return
		}

		service.RegisterUserServiceServer(grpcServer, grpcRouter)
		log.Info().Msgf(`GRPC server (%s) started successfully 🚀`, grpcAddress)

		errs <- grpcServer.Serve(grpcListener)
	}()

	// Implement shutdown signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errs:
		log.Err(err).Msg("exit")
	case sig := <-signalChan:
		log.Info().Msgf(`received signal: %v`, sig)

		// Graceful shutdown
		grpcServer.GracefulStop()
	}

	log.Info().Msg("exiting")

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
